package app

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/dao"
	appdao "github.com/chenjie199234/admin/dao/app"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/pool"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/web"
)

// Service subservice for config business
type Service struct {
	stop *graceful.Graceful

	configDao     *appdao.Dao
	permissionDao *permissiondao.Dao

	sync.Mutex

	apps          map[string]*model.AppSummary            //key:appgroup.appname,value:appinfo
	notices       map[string]map[chan *struct{}]*struct{} //key:appgroup.appname,value:waiting chans
	clients       map[string]*crpc.CrpcClient             //key:appgroup.appname,value:client
	clientsActive map[string]int64                        //key:appgroup.appname,value:last use timestamp(unixnano)
}

// Start -
func Start() *Service {
	s := &Service{
		stop: graceful.New(),

		configDao:     appdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),

		apps:          make(map[string]*model.AppSummary),
		notices:       make(map[string]map[chan *struct{}]*struct{}),
		clients:       make(map[string]*crpc.CrpcClient),
		clientsActive: make(map[string]int64),
	}
	if e := s.configDao.MongoWatchConfig(s.drop, s.update, s.del, s.apps); e != nil {
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
		for name, last := range s.clientsActive {
			if now.UnixNano()-last < time.Minute.Nanoseconds() {
				continue
			}
			delete(s.clientsActive, name)
			client, ok := s.clients[name]
			if !ok {
				continue
			}
			go client.Close(false)
			delete(s.clients, name)
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
	for _, client := range s.clients {
		go client.Close(false)
	}
	s.clients = make(map[string]*crpc.CrpcClient)
	s.clientsActive = make(map[string]int64)
}
func (s *Service) update(gname, aname string, cur *model.AppSummary) {
	log.Debug(nil, "[update] group:", gname, "app:", aname, "keys:", cur.Keys)
	s.Lock()
	defer s.Unlock()
	if s.stop.Closed() {
		return
	}
	s.apps[gname+"."+aname] = cur
	for notice := range s.notices[gname+"."+aname] {
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
	log.Debug(nil, "[del] group:", app.Group, "app:", app.App, "summary")
	delete(s.apps, app.Group+"."+app.App)
	if client, ok := s.clients[app.Group+"."+app.App]; ok {
		delete(s.clients, app.Group+"."+app.App)
		delete(s.clientsActive, app.Group+"."+app.App)
		go client.Close(false)
	}
	for notice := range s.notices[app.Group+"."+app.App] {
		select {
		case notice <- nil:
		default:
		}
	}
}

// get all groups
func (s *Service) ListGroup(ctx context.Context, req *api.ListGroupReq) (*api.ListGroupResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[Groups] operator:", md["Token-Data"], "format wrong:", e)
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
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.AppControl, true)
		if e != nil {
			log.Error(ctx, "[Groups] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	groups, e := s.configDao.MongoGetAllGroups(ctx, projectid)
	if e != nil {
		log.Error(ctx, "[Groups]", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.ListGroupResp{Groups: groups}, nil
}

// get all apps in one specific group
func (s *Service) ListApp(ctx context.Context, req *api.ListAppReq) (*api.ListAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[Apps] operator:", md["Token-Data"], "format wrong:", e)
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
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.AppControl, true)
		if e != nil {
			log.Error(ctx, "[Apps] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	apps, e := s.configDao.MongoGetAllAppsInGroup(ctx, projectid, req.GName)
	if e != nil {
		log.Error(ctx, "[Apps] group:", req.GName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.ListAppResp{Apps: apps}, nil
}

// create one specific app
func (s *Service) CreateApp(ctx context.Context, req *api.CreateAppReq) (*api.CreateAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[CreateApp] operator:", md["Token-Data"], "format wrong:", e)
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
			log.Error(ctx, "[CreateApp] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.configDao.MongoCreateApp(ctx, projectid, req.GName, req.AName, req.Secret); e != nil {
		log.Error(ctx, "[CreateApp] group:", req.GName, "app:", req.AName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[CreateApp] group:", req.GName, "app:", req.AName, "success")
	return &api.CreateAppResp{}, nil
}

// del one specific app in one specific group
func (s *Service) DelApp(ctx context.Context, req *api.DelAppReq) (*api.DelAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelApp] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.configDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[DelApp] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "3" {
			log.Error(ctx, "[DelApp] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		confignodeid := nodeid[:strings.LastIndex(nodeid, ",")]
		projectid := confignodeid[:strings.LastIndex(confignodeid, ",")]
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, confignodeid, true)
		if e != nil {
			log.Error(ctx, "[DelApp] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.configDao.MongoDelApp(ctx, req.GName, req.AName, req.Secret); e != nil {
		log.Error(ctx, "[DelApp] group:", req.GName, "app:", req.AName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelApp] group:", req.GName, "app:", req.AName, "success")
	return &api.DelAppResp{}, nil
}

// update one specific app's cipher
func (s *Service) UpdateAppSecret(ctx context.Context, req *api.UpdateAppSecretReq) (*api.UpdateAppSecretResp, error) {
	if req.OldSecret == req.NewSecret {
		return &api.UpdateAppSecretResp{}, nil
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateAppSecret] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.configDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[UpdateAppSecret] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "3" {
			log.Error(ctx, "[UpdateAppSecret] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[UpdateAppSecret] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.configDao.MongoUpdateAppSecret(ctx, req.GName, req.AName, req.OldSecret, req.NewSecret); e != nil {
		log.Error(ctx, "[UpdateAppSecret] group:", req.GName, "app:", req.AName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[UpdateAppSecret] group:", req.GName, "app:", req.AName, "success")
	return &api.UpdateAppSecretResp{}, nil
}

// get all config's keys in this app
func (s *Service) ListKey(ctx context.Context, req *api.ListKeyReq) (*api.ListKeyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[Keys] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.configDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[Keys] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "3" {
			log.Error(ctx, "[Keys] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[Keys] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	keys, e := s.configDao.MongoGetAllKeys(ctx, req.GName, req.AName, req.Secret)
	if e != nil {
		log.Error(ctx, "[Keys] group:", req.GName, "app:", req.AName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.ListKeyResp{Keys: keys}, nil
}

// del one specific key in one specific app
func (s *Service) DelKey(ctx context.Context, req *api.DelKeyReq) (*api.DelKeyResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelKey] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.configDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[DelKey] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "3" {
			log.Error(ctx, "[DelKey] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[DelKey] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	if e := s.configDao.MongoDelKey(ctx, req.GName, req.AName, req.Key, req.Secret); e != nil {
		log.Error(ctx, "[DelKey] group:", req.GName, "app:", req.AName, "key:", req.Key, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelKey] group:", req.GName, "app:", req.AName, "key:", req.Key)
	return &api.DelKeyResp{}, nil
}

// get config
func (s *Service) GetKeyConfig(ctx context.Context, req *api.GetKeyConfigReq) (*api.GetKeyConfigResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[GetKeyConfig] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.configDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[GetKeyConfig] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "3" {
			log.Error(ctx, "[GetKeyConfig] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[GetKeyConfig] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	keysummary, configlog, e := s.configDao.MongoGetKeyConfig(ctx, req.GName, req.AName, req.Key, req.Index, req.Secret)
	if e != nil {
		log.Error(ctx, "[GetKeyConfig] group:", req.GName, "app:", req.AName, "key:", req.Key, "index:", req.Index, e)
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

// set config
func (s *Service) SetKeyConfig(ctx context.Context, req *api.SetKeyConfigReq) (*api.SetKeyConfigResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SetKeyConfig] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.configDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[SetKeyConfig] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "3" {
			log.Error(ctx, "[SetKeyConfig] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[SetKeyConfig] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	req.Key = strings.TrimSpace(req.Key)
	if req.Key == "" {
		log.Error(ctx, "[SetKeyConfig] group:", req.GName, "app:", req.AName, "key empty")
		return nil, ecode.ErrReq
	}
	req.Value = strings.TrimSpace(req.Value)
	if req.Value == "" {
		log.Error(ctx, "[SetKeyConfig] group:", req.GName, "app:", req.AName, "value empty")
		return nil, ecode.ErrReq
	}
	index, version, e := s.configDao.MongoSetKeyConfig(ctx, req.GName, req.AName, req.Key, req.Secret, req.Value, req.ValueType)
	if e != nil {
		log.Error(ctx, "[SetKeyConfig] group:", req.GName, "app:", req.AName, "key:", req.Key, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[SetKeyConfig] group:", req.GName, "app:", req.AName, "key:", req.Key, "index:", index, "version:", version, "success")
	return &api.SetKeyConfigResp{}, nil
}

// rollback config
func (s *Service) Rollback(ctx context.Context, req *api.RollbackReq) (*api.RollbackResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[Rollback] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.configDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[Rollback] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "3" {
			log.Error(ctx, "[Rollback] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[Rollback] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.configDao.MongoRollbackKeyConfig(ctx, req.GName, req.AName, req.Key, req.Secret, req.Index); e != nil {
		log.Error(ctx, "[Rollback] group:", req.GName, "app:", req.AName, "key:", req.Key, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[Rollback] group:", req.GName, "app:", req.AName, "key:", req.Key, "index:", req.Index, "success")
	return &api.RollbackResp{}, nil
}

// watch config
func (s *Service) Watch(ctx context.Context, req *api.WatchReq) (*api.WatchResp, error) {
	for k := range req.Keys {
		if strings.Contains(k, ".") {
			return nil, ecode.ErrReq
		}
	}
	if !s.stop.AddOne() {
		return nil, cerror.ErrClosing
	}
	defer s.stop.DoneOne()

	resp := &api.WatchResp{
		Datas: make(map[string]*api.WatchData, len(req.Keys)+3),
	}

	s.Lock()
	a, ok := s.apps[req.GName+"."+req.AName]
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
		if _, ok := s.notices[req.GName+"."+req.AName]; !ok {
			s.notices[req.GName+"."+req.AName] = map[chan *struct{}]*struct{}{ch: nil}
		} else {
			s.notices[req.GName+"."+req.AName][ch] = nil
		}
		s.Unlock()
		select {
		case <-ctx.Done():
			s.Lock()
			delete(s.notices[req.GName+"."+req.AName], ch)
			s.Unlock()
			return nil, ctx.Err()
		case _, ok := <-ch:
			if !ok {
				return nil, cerror.ErrClosing
			}
		}
		s.Lock()
		delete(s.notices[req.GName+"."+req.AName], ch)
		a, ok = s.apps[req.GName+"."+req.AName]
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

func (s *Service) ListProxy(ctx context.Context, req *api.ListProxyReq) (*api.ListProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[ListProxy] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.configDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[ListProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "3" {
			log.Error(ctx, "[ListProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[ListProxy] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread || !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	paths, e := s.configDao.MongoListProxyPath(ctx, req.GName, req.AName, req.Secret)
	if e != nil {
		log.Error(ctx, "[ListProxy] group:", req.GName, "app:", req.AName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.ListProxyResp{
		Paths: make(map[string]*api.ProxyPathInfo),
	}
	for path, info := range paths {
		nodeid, e := util.ParseNodeIDstr(info.PermissionNodeID)
		if e != nil {
			log.Error(ctx, "[ListProxy] group:", req.GName, "app:", req.AName, "path:", path, "nodeid:", info.PermissionNodeID, "format wrong:", e)
			return nil, ecode.ErrConfigDataBroken
		}
		resp.Paths[path] = &api.ProxyPathInfo{
			NodeId: nodeid,
			Read:   info.PermissionRead,
			Write:  info.PermissionWrite,
		}
	}
	return resp, nil
}
func (s *Service) SetProxy(ctx context.Context, req *api.SetProxyReq) (*api.SetProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SetProxy] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.configDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[SetProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "3" {
			log.Error(ctx, "[SetProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[SetProxy] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.configDao.MongoSetProxyPath(ctx, req.GName, req.AName, req.Secret, req.Path, req.Read, req.Write); e != nil {
		log.Error(ctx, "[SetProxy] group:", req.GName, "app:", req.AName, "path:", req.Path, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[SetProxy] group:", req.GName, "app:", req.AName, "path:", req.Path, "read:", req.Read, "write:", req.Write, "success")
	return &api.SetProxyResp{}, nil
}
func (s *Service) DelProxy(ctx context.Context, req *api.DelProxyReq) (*api.DelProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelProxy] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.configDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[DelProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "3" {
			log.Error(ctx, "[DelProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[DelProxy] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.configDao.MongoDelProxyPath(ctx, req.GName, req.AName, req.Secret, req.Path); e != nil {
		log.Error(ctx, "[DelProxy] group:", req.GName, "app:", req.AName, "path:", req.Path, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelProxy] group:", req.GName, "app:", req.AName, "path:", req.Path, "success")
	return &api.DelProxyResp{}, nil
}
func (s *Service) Proxy(ctx context.Context, req *api.ProxyReq) (*api.ProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[Proxy] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	s.Lock()
	app, ok := s.apps[req.GName+"."+req.AName]
	if !ok {
		s.Unlock()
		return nil, ecode.ErrAppNotExist
	}
	pathinfo, ok := app.Paths[req.Path]
	if !ok {
		s.Unlock()
		return nil, ecode.ErrProxyPathNotExist
	}
	client, ok := s.clients[req.GName+"."+req.AName]
	if !ok {
		var e error
		client, e = crpc.NewCrpcClient(dao.GetCrpcClientConfig(), model.Group, model.Name, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[Proxy] new crpc client to group:", req.GName, "app:", req.AName, e)
			s.Unlock()
			return nil, ecode.ErrSystem
		}
		s.clients[req.GName+"."+req.AName] = client
	}
	s.clientsActive[req.GName+"."+req.AName] = time.Now().UnixNano()
	s.Unlock()
	if !operator.IsZero() && (pathinfo.PermissionRead || pathinfo.PermissionWrite) {
		//permission check
		canread, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, pathinfo.PermissionNodeID, true)
		if e != nil {
			log.Error(ctx, "[Proxy] operator:", md["Token-Data"], "nodeid:", pathinfo.PermissionNodeID, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if pathinfo.PermissionRead && !canread && !admin {
			return nil, ecode.ErrPermission
		}
		if pathinfo.PermissionWrite && !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	out, e := client.Call(ctx, req.Path, common.Str2byte(req.Data), metadata.GetMetadata(ctx))
	if e != nil {
		log.Error(ctx, "[Proxy] operator:", md["Token-Data"], "call group:", req.GName, "app:", req.AName, "path:", req.Path, "reqdata:", req.Data, e)
		return nil, e
	}
	log.Info(ctx, "[Proxy] operator:", md["Token-Data"], "call group:", req.GName, "app:", req.AName, "path:", req.Path, "reqdata:", req.Data, "respdata:", out)
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
