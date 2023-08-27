package app

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/dao"
	appdao "github.com/chenjie199234/admin/dao/app"
	initializedao "github.com/chenjie199234/admin/dao/initialize"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/pool"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "github.com/chenjie199234/Corelib/web"
	//"github.com/chenjie199234/Corelib/cgrpc"
)

// Service subservice for config business
type Service struct {
	stop *graceful.Graceful

	appDao        *appdao.Dao
	permissionDao *permissiondao.Dao
	initializeDao *initializedao.Dao

	sync.Mutex

	apps          map[string]*model.AppSummary            //key:projectid-group.app,value:appinfo
	notices       map[string]map[chan *struct{}]*struct{} //key:projectid-group.app,value:waiting chans
	clients       map[string]*clientinfo                  //key:projectid-group.app,value:clientinfo
	clientsActive map[string]int64                        //key:projectid-group.app,value:last use timestamp(unixnano)
}
type clientinfo struct {
	di     discover.DI
	client *crpc.CrpcClient
}

// Start -
func Start() *Service {
	s := &Service{
		stop: graceful.New(),

		appDao:        appdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		initializeDao: initializedao.NewDao(nil, nil, config.GetMongo("admin_mongo")),

		apps:          make(map[string]*model.AppSummary),
		notices:       make(map[string]map[chan *struct{}]*struct{}),
		clients:       make(map[string]*clientinfo),
		clientsActive: make(map[string]int64),
	}
	if e := s.appDao.MongoWatchConfig(s.drop, s.update, s.del, s.apps); e != nil {
		panic("[Config.Start] watch error: " + e.Error())
	}
	go s.job()
	return s
}

func (s *Service) job() {
	tker := time.NewTicker(time.Minute)
	for {
		now := <-tker.C
		s.Lock()
		for fullname, last := range s.clientsActive {
			if now.UnixNano()-last < time.Minute.Nanoseconds() {
				continue
			}
			delete(s.clientsActive, fullname)
			c, ok := s.clients[fullname]
			if !ok {
				continue
			}
			go func() {
				c.client.Close(false)
				c.di.Stop()
			}()
			delete(s.clients, fullname)
		}
		s.Unlock()
	}
}
func (s *Service) drop() {
	s.Lock()
	defer s.Unlock()
	s.apps = make(map[string]*model.AppSummary)
	for _, notices := range s.notices {
		for notice := range notices {
			select {
			case notice <- nil:
			default:
			}
		}
	}
	for _, v := range s.clients {
		c := v
		go func() {
			c.client.Close(false)
			c.di.Stop()
		}()
	}
	s.clients = make(map[string]*clientinfo)
	s.clientsActive = make(map[string]int64)
}
func (s *Service) update(cur *model.AppSummary) {
	log.Debug(nil, "[update]", map[string]interface{}{"project_id": cur.ProjectID, "group": cur.Group, "app": cur.App, "keys": cur.Keys})
	s.Lock()
	defer s.Unlock()
	if s.stop.Closed() {
		return
	}
	s.apps[cur.GetFullName()] = cur
	for notice := range s.notices[cur.GetFullName()] {
		select {
		case notice <- nil:
		default:
		}
	}
}
func (s *Service) del(id string) {
	s.Lock()
	defer s.Unlock()
	if s.stop.Closed() {
		return
	}
	var app *model.AppSummary
	for _, v := range s.apps {
		if v.ID.Hex() == id {
			app = v
			break
		}
	}
	if app == nil {
		//the deleted doc is not the summary doc
		return
	}
	//delete the summary,this is same as delete the app
	log.Debug(nil, "[del]", map[string]interface{}{"project_id": app.ProjectID, "group": app.Group, "app": app.App})
	delete(s.apps, app.GetFullName())
	if c, ok := s.clients[app.GetFullName()]; ok {
		delete(s.clients, app.GetFullName())
		delete(s.clientsActive, app.GetFullName())
		go func() {
			c.client.Close(false)
			c.di.Stop()
		}()
	}
	for notice := range s.notices[app.GetFullName()] {
		select {
		case notice <- nil:
		default:
		}
	}
}

func (s *Service) GetApp(ctx context.Context, req *api.GetAppReq) (*api.GetAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[GetApp] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[GetApp] get app's permission nodeid failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[GetApp] app's permission nodeid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "nodeid": nodeid})
			return nil, ecode.ErrDataBroken
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[GetApp] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": nodeid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	app, e := s.appDao.MongoGetApp(ctx, projectid, req.GName, req.AName, req.Secret)
	if e != nil {
		log.Error(ctx, "[GetApp] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.GetAppResp{
		Keys:  make(map[string]*api.KeyConfigInfo),
		Paths: make(map[string]*api.ProxyPathInfo),
	}
	for k, v := range app.Keys {
		resp.Keys[k] = &api.KeyConfigInfo{
			CurIndex:     v.CurIndex,
			MaxIndex:     v.MaxIndex,
			CurVersion:   v.CurVersion,
			CurValue:     v.CurValue,
			CurValueType: v.CurValueType,
		}
	}
	for k, v := range app.Paths {
		nodeid, e := util.ParseNodeIDstr(v.PermissionNodeID)
		if e != nil {
			log.Error(ctx, "[GetApp] app's path's permission nodeid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "path": k, "nodeid": v.PermissionNodeID, "error": e})
			return nil, ecode.ErrDataBroken
		}
		resp.Paths[k] = &api.ProxyPathInfo{
			NodeId: nodeid,
			Read:   v.PermissionRead,
			Write:  v.PermissionWrite,
			Admin:  v.PermissionAdmin,
		}
	}
	return resp, nil
}

func (s *Service) CreateApp(ctx context.Context, req *api.CreateAppReq) (*api.CreateAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[CreateApp] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	if !operator.IsZero() {
		//config control permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.AppControl, true)
		if e != nil {
			log.Error(ctx, "[CreateApp] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": projectid + model.AppControl, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	nodeidstr, e := s.appDao.MongoCreateApp(ctx, projectid, req.GName, req.AName, req.Secret)
	if e != nil {
		log.Error(ctx, "[CreateApp] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeid, _ := util.ParseNodeIDstr(nodeidstr)
	log.Info(ctx, "[CreateApp] success", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName})
	return &api.CreateAppResp{NodeId: nodeid}, nil
}

func (s *Service) DelApp(ctx context.Context, req *api.DelAppReq) (*api.DelAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelApp] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
	if e != nil {
		log.Error(ctx, "[DelApp] get app's permission nodeid failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		log.Error(ctx, "[DelApp] app's permission nodeid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "nodeid": nodeid})
		return nil, ecode.ErrDataBroken
	}
	//self can't be deleted
	if nodeids[1] == "1" && nodeids[3] == "1" {
		log.Error(ctx, "[DelApp] can't delete self", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName})
		return nil, ecode.ErrPermission
	}

	if !operator.IsZero() {
		//config control permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.AppControl, true)
		if e != nil {
			log.Error(ctx, "[DelApp] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": projectid + model.AppControl, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoDelApp(ctx, projectid, req.GName, req.AName, req.Secret); e != nil {
		log.Error(ctx, "[DelApp] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelApp] success", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName})
	return &api.DelAppResp{}, nil
}

func (s *Service) UpdateAppSecret(ctx context.Context, req *api.UpdateAppSecretReq) (*api.UpdateAppSecretResp, error) {
	if req.OldSecret == req.NewSecret {
		return &api.UpdateAppSecretResp{}, nil
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateAppSecret] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[UpdateAppSecret] get app's permission nodeid failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[UpdateAppSecret] app's permission nodeid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "nodeid": nodeid})
			return nil, ecode.ErrDataBroken
		}
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[UpdateAppSecret] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": nodeid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoUpdateAppSecret(ctx, projectid, req.GName, req.AName, req.OldSecret, req.NewSecret); e != nil {
		log.Error(ctx, "[UpdateAppSecret] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[UpdateAppSecret] success", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName})
	return &api.UpdateAppSecretResp{}, nil
}

func (s *Service) DelKey(ctx context.Context, req *api.DelKeyReq) (*api.DelKeyResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelKey] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
	if e != nil {
		log.Error(ctx, "[DelKey] get app's permission nodeid failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		log.Error(ctx, "[DelKey] app's permission nodeid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "nodeid": nodeid})
		return nil, ecode.ErrDataBroken
	}
	if nodeids[1] == "1" && nodeids[3] == "1" && (req.Key == "AppConfig" || req.Key == "SourceConfig") {
		log.Error(ctx, "[DelKey] can't delete self's 'AppConfig' or 'SourceConfig' key", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName})
		return nil, ecode.ErrPermission
	}

	if !operator.IsZero() {
		//config control permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[DelKey] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": nodeid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	if e := s.appDao.MongoDelKey(ctx, projectid, req.GName, req.AName, req.Key, req.Secret); e != nil {
		log.Error(ctx, "[DelKey] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "key": req.Key, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelKey] success", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "key": req.Key})
	return &api.DelKeyResp{}, nil
}

func (s *Service) GetKeyConfig(ctx context.Context, req *api.GetKeyConfigReq) (*api.GetKeyConfigResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[GetKeyConfig] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[GetKeyConfig] get app's permission nodeid failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[GetKeyConfig] app's permission nodeid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "nodeid": nodeid})
			return nil, ecode.ErrDataBroken
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[GetKeyConfig] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": nodeid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	keysummary, configlog, e := s.appDao.MongoGetKeyConfig(ctx, projectid, req.GName, req.AName, req.Key, req.Index, req.Secret)
	if e != nil {
		log.Error(ctx, "[GetKeyConfig] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "key": req.Key, "index": req.Index, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.GetKeyConfigResp{
		CurIndex:   keysummary.CurIndex,
		MaxIndex:   keysummary.MaxIndex,
		CurVersion: keysummary.CurVersion,
		ThisIndex:  configlog.Index,
		Value:      configlog.Value,
		ValueType:  configlog.ValueType,
	}, nil
}

func (s *Service) SetKeyConfig(ctx context.Context, req *api.SetKeyConfigReq) (*api.SetKeyConfigResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SetKeyConfig] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	req.Key = strings.TrimSpace(req.Key)
	if req.Key == "" {
		log.Error(ctx, "[SetKeyConfig] key empty", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName})
		return nil, ecode.ErrReq
	}
	req.Value = strings.TrimSpace(req.Value)
	if req.Value == "" {
		log.Error(ctx, "[SetKeyConfig] value empty", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "key": req.Key})
		return nil, ecode.ErrReq
	}
	switch req.ValueType {
	case "json":
		buf := bytes.NewBuffer(nil)
		if e := json.Compact(buf, common.Str2byte(req.Value)); e != nil {
			log.Error(ctx, "[SetKeyConfig] json value format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "key": req.Key, "error": e})
			return nil, ecode.ErrReq
		}
		req.Value = common.Byte2str(buf.Bytes())
	case "toml":
		//TODO
		fallthrough
	case "yaml":
		//TODO
		fallthrough
	case "raw":
		//TODO
		fallthrough
	default:
		log.Error(ctx, "[SetKeyConfig] unsupported value type", map[string]interface{}{"operator": md["Token-Data"], "projectid": projectid, "group": req.GName, "app": req.AName, "key": req.Key, "valuetype": req.ValueType})
		return nil, ecode.ErrReq
	}

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[SetKeyConfig] get app's permission nodeid failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[SetKeyConfig] app's permission nodeid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "nodeid": nodeid})
			return nil, ecode.ErrDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[SetKeyConfig] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": nodeid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	index, version, e := s.appDao.MongoSetKeyConfig(ctx, projectid, req.GName, req.AName, req.Key, req.Secret, req.Value, req.ValueType, req.NewKey)
	if e != nil {
		log.Error(ctx, "[SetKeyConfig] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "key": req.Key, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[SetKeyConfig] success", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "key": req.Key, "new_index": index, "new_version": version})
	return &api.SetKeyConfigResp{}, nil
}

func (s *Service) Rollback(ctx context.Context, req *api.RollbackReq) (*api.RollbackResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[Rollback] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[Rollback] get app's permission nodeid failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[Rollback] app's permission nodeid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "nodeid": nodeid})
			return nil, ecode.ErrDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[Rollback] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": nodeid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoRollbackKeyConfig(ctx, projectid, req.GName, req.AName, req.Key, req.Secret, req.Index); e != nil {
		log.Error(ctx, "[Rollback] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "key": req.Key, "index": req.Index, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[Rollback] success", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "key": req.Key, "index": req.Index})
	return &api.RollbackResp{}, nil
}

func (s *Service) Watch(ctx context.Context, req *api.WatchReq) (*api.WatchResp, error) {
	for k := range req.Keys {
		if strings.Contains(k, ".") {
			return nil, ecode.ErrReq
		}
	}
	if !s.stop.AddOne() {
		return nil, cerror.ErrServerClosing
	}
	defer s.stop.DoneOne()
	projectid, e := s.initializeDao.MongoGetProjectIDByName(ctx, req.ProjectName)
	if e != nil {
		log.Error(ctx, "[Watch] get projectid failed", map[string]interface{}{"project_name": req.ProjectName, "group": req.GName, "app": req.AName, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}

	resp := &api.WatchResp{
		Datas: make(map[string]*api.WatchData, len(req.Keys)+3),
	}

	s.Lock()
	a, ok := s.apps[projectid+"-"+req.GName+"."+req.AName]
	if !ok {
		s.Unlock()
		return nil, ecode.ErrAppNotExist
	}
	needreturn := false
	for key, clientversion := range req.Keys {
		k, ok := a.Keys[key]
		if !ok || k == nil || k.CurVersion == 0 {
			s.Unlock()
			return nil, ecode.ErrKeyNotExist
		}
		if clientversion == 0 || clientversion != k.CurVersion {
			needreturn = true
			resp.Datas[key] = &api.WatchData{
				Key:       key,
				Value:     k.CurValue,
				ValueType: k.CurValueType,
				Version:   k.CurVersion,
			}
		} else {
			resp.Datas[key] = &api.WatchData{
				Key:       key,
				Value:     "",
				ValueType: "",
				Version:   k.CurVersion,
			}
		}
	}
	if needreturn {
		s.Unlock()
		return resp, nil
	}
	for {
		ch := make(chan *struct{})
		if _, ok := s.notices[projectid+"-"+req.GName+"."+req.AName]; !ok {
			s.notices[projectid+"-"+req.GName+"."+req.AName] = map[chan *struct{}]*struct{}{ch: nil}
		} else {
			s.notices[projectid+"-"+req.GName+"."+req.AName][ch] = nil
		}
		s.Unlock()
		select {
		case <-ctx.Done():
			s.Lock()
			delete(s.notices[projectid+"-"+req.GName+"."+req.AName], ch)
			s.Unlock()
			return nil, cerror.ConvertStdError(ctx.Err())
		case _, ok := <-ch:
			if !ok {
				return nil, cerror.ErrServerClosing
			}
		}
		s.Lock()
		delete(s.notices[projectid+"-"+req.GName+"."+req.AName], ch)
		a, ok = s.apps[projectid+"-"+req.GName+"."+req.AName]
		if !ok {
			s.Unlock()
			return nil, ecode.ErrAppNotExist
		}
		for key, clientversion := range req.Keys {
			k, ok := a.Keys[key]
			if !ok || k == nil || k.CurVersion == 0 {
				s.Unlock()
				return nil, ecode.ErrKeyNotExist
			}
			if clientversion == 0 || clientversion != k.CurVersion {
				needreturn = true
				resp.Datas[key] = &api.WatchData{
					Key:       key,
					Value:     k.CurValue,
					ValueType: k.CurValueType,
					Version:   k.CurVersion,
				}
			} else {
				resp.Datas[key] = &api.WatchData{
					Key:       key,
					Value:     "",
					ValueType: "",
					Version:   k.CurVersion,
				}
			}
		}
		if needreturn {
			s.Unlock()
			return resp, nil
		}
	}
}

func (s *Service) SetProxy(ctx context.Context, req *api.SetProxyReq) (*api.SetProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SetProxy] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[SetProxy] get app's permission nodeid failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[SetProxy] app's permission nodeid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "nodeid": nodeid})
			return nil, ecode.ErrDataBroken
		}
		if nodeids[1] == "1" && nodeids[3] == "1" {
			log.Error(ctx, "[SetProxy] can't set proxy path for self", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName})
			return nil, ecode.ErrPermission
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[SetProxy] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": nodeid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if req.Path[0] != '/' {
		req.Path = "/" + req.Path
	}
	newnodeidstr, e := s.appDao.MongoSetProxyPath(ctx, projectid, req.GName, req.AName, req.Secret, req.Path, req.Read, req.Write, req.Admin, req.NewPath)
	if e != nil {
		log.Error(ctx, "[SetProxy] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "path": req.Path, "read_write_admin": []bool{req.Read, req.Write, req.Admin}, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	newnodeid, _ := util.ParseNodeIDstr(newnodeidstr)
	log.Info(ctx, "[SetProxy] success", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "path": req.Path, "read_write_admin": []bool{req.Read, req.Write, req.Admin}})
	return &api.SetProxyResp{NodeId: newnodeid}, nil
}
func (s *Service) DelProxy(ctx context.Context, req *api.DelProxyReq) (*api.DelProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelProxy] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[DelProxy] get app's permission nodeid failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[DelProxy] app's permission nodeid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "nodeid": nodeid})
			return nil, ecode.ErrDataBroken
		}
		if nodeids[1] == "1" && nodeids[3] == "1" {
			log.Error(ctx, "[DelProxy] can't delete proxy path for self", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName})
			return nil, ecode.ErrPermission
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[DelProxy] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": nodeid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoDelProxyPath(ctx, projectid, req.GName, req.AName, req.Secret, req.Path); e != nil {
		log.Error(ctx, "[DelProxy] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "path": req.Path, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelProxy] success", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "path": req.Path})
	return &api.DelProxyResp{}, nil
}
func (s *Service) Proxy(ctx context.Context, req *api.ProxyReq) (*api.ProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[Proxy] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()

	s.Lock()
	app, ok := s.apps[projectid+"-"+req.GName+"."+req.AName]
	if !ok {
		s.Unlock()
		return nil, ecode.ErrAppNotExist
	}
	if req.Path[0] != '/' {
		req.Path = "/" + req.Path
	}
	pathinfo, ok := app.Paths[req.Path]
	if !ok {
		s.Unlock()
		return nil, ecode.ErrProxyPathNotExist
	}
	c, ok := s.clients[projectid+"-"+req.GName+"."+req.AName]
	if !ok {
		projectname, e := s.initializeDao.MongoGetProjectNameByID(ctx, projectid)
		if e != nil {
			log.Error(ctx, "[Proxy] get project name failed", map[string]interface{}{"project_id": projectid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		di, e := discover.NewDNSDiscover(projectname, req.GName, req.AName, req.AName+"-headless."+projectname+"-"+req.GName, time.Second*10, 9000, 10000, 8000)
		if e != nil {
			log.Error(ctx, "[Proxy] new dns discover failed", map[string]interface{}{"project": projectname, "group": req.GName, "app": req.AName, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		client, e := crpc.NewCrpcClient(dao.GetCrpcClientConfig(), di, model.Project, model.Group, model.Name, projectname, req.GName, req.AName, nil)
		if e != nil {
			log.Error(ctx, "[Proxy] new crpc client failed", map[string]interface{}{"project_id": projectid, "group": req.GName, "app": req.AName, "error": e})
			s.Unlock()
			return nil, ecode.ErrSystem
		}
		c = &clientinfo{
			di:     di,
			client: client,
		}
		s.clients[projectid+"-"+req.GName+"."+req.AName] = c
	}
	s.clientsActive[projectid+"-"+req.GName+"."+req.AName] = time.Now().UnixNano()
	s.Unlock()
	if !operator.IsZero() && (pathinfo.PermissionRead || pathinfo.PermissionWrite || pathinfo.PermissionAdmin) {
		//permission check
		canread, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, pathinfo.PermissionNodeID, true)
		if e != nil {
			log.Error(ctx, "[Proxy] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "nodeid": pathinfo.PermissionNodeID, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if pathinfo.PermissionAdmin && !admin {
			return nil, ecode.ErrPermission
		} else if pathinfo.PermissionWrite && !canwrite && !admin {
			return nil, ecode.ErrPermission
		} else if pathinfo.PermissionRead && !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	out, e := c.client.Call(ctx, req.Path, common.Str2byte(req.Data), metadata.GetMetadata(ctx))
	if e != nil {
		log.Error(ctx, "[Proxy] call server failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "path": req.Path, "reqdata": req.Data, "error": e})
		return nil, e
	}
	log.Info(ctx, "[Proxy] success", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "group": req.GName, "app": req.AName, "path": req.Path, "reqdata": req.Data, "respdata": out})
	return &api.ProxyResp{Data: common.Byte2str(out)}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(func() {
		s.Lock()
		for _, notices := range s.notices {
			for n := range notices {
				close(n)
			}
		}
		s.Unlock()
	}, nil)
}
