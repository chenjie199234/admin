package app

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
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
	"github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/pool"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/egroup"
	"github.com/chenjie199234/Corelib/util/graceful"
	"github.com/chenjie199234/Corelib/web"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/chenjie199234/Corelib/cgrpc"
)

// Service subservice for config business
type Service struct {
	stop *graceful.Graceful

	appDao        *appdao.Dao
	permissionDao *permissiondao.Dao

	sync.Mutex

	apps          map[string]*model.AppSummary            //key:appgroup.appname,value:appinfo
	notices       map[string]map[chan *struct{}]*struct{} //key:appgroup.appname,value:waiting chans
	clients       map[string]*clientinfo                  //key:appgroup.appname,value:clientinfo
	clientsActive map[string]int64                        //key:appgroup.appname,value:last use timestamp(unixnano)
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
		for name, last := range s.clientsActive {
			if now.UnixNano()-last < time.Minute.Nanoseconds() {
				continue
			}
			delete(s.clientsActive, name)
			c, ok := s.clients[name]
			if !ok {
				continue
			}
			go func() {
				c.client.Close(false)
				c.di.Stop()
			}()
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
	if c, ok := s.clients[app.Group+"."+app.App]; ok {
		delete(s.clients, app.Group+"."+app.App)
		delete(s.clientsActive, app.Group+"."+app.App)
		go func() {
			c.client.Close(false)
			c.di.Stop()
		}()
	}
	for notice := range s.notices[app.Group+"."+app.App] {
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
		log.Error(ctx, "[GetApp] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[GetApp] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[GetApp] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
			return nil, ecode.ErrConfigDataBroken
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[GetApp] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	app, e := s.appDao.MongoGetApp(ctx, req.GName, req.AName, req.Secret)
	if e != nil {
		log.Error(ctx, "[GetApp] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, e)
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
			log.Error(ctx, "[GetApp] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "path:", k, "nodeid:", v.PermissionNodeID, "format wrong:", e)
			return nil, ecode.ErrConfigDataBroken
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

func (s *Service) AppInstances(ctx context.Context, req *api.AppInstancesReq) (*api.AppInstancesResp, error) {
	md := metadata.GetMetadata(ctx)
	s.Lock()
	app, ok := s.apps[req.GName+"."+req.AName]
	if !ok {
		s.Unlock()
		return nil, ecode.ErrAppNotExist
	}
	s.Unlock()
	if e := util.SignCheck(req.Secret, app.Value); e != nil {
		log.Error(ctx, "[AppInstances] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var ips []string
	tmpdi := discover.NewDNSDiscover(req.GName, req.AName, req.AName+"-headless."+req.GName, time.Second*10, 9000, 10000, 8000)
	defer tmpdi.Stop()
	notice, cancel := tmpdi.GetNotice()
	defer cancel()
	select {
	case <-notice:
		addrs, e := tmpdi.GetAddrs(discover.NotNeed)
		if e != nil {
			log.Error(ctx, "[AppInstances] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		ips = make([]string, 0, len(addrs))
		for addr := range addrs {
			ips = append(ips, addr)
		}
	case <-ctx.Done():
		return nil, cerror.ConvertStdError(ctx.Err())
	}
	resp := &api.AppInstancesResp{
		Instances: make([]*api.InstanceInfo, 0, len(ips)),
	}
	eg := egroup.GetGroup(ctx)
	for _, v := range ips {
		ip := v
		info := &api.InstanceInfo{}
		resp.Instances = append(resp.Instances, info)
		eg.Go(func(gctx context.Context) error {
			di := discover.NewDirectDiscover(req.GName, req.AName, "http://"+ip+":6060", 9000, 10000, 8000)
			defer di.Stop()
			webclient, e := web.NewWebClient(dao.GetWebClientConfig(), di, model.Group, model.Name, req.GName, req.AName, nil)
			if e != nil {
				log.Error(ctx, "[AppInstances] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "ip:", ip, e)
				return e
			}
			defer webclient.Close(false)
			r, e := webclient.Get(gctx, "/info", "", nil, nil)
			if e != nil {
				log.Error(ctx, "[AppInstances] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "ip:", ip, e)
				return e
			}
			defer r.Body.Close()
			body, e := io.ReadAll(r.Body)
			if e != nil {
				log.Error(ctx, "[AppInstances] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "ip:", ip, e)
				return e
			}
			tmp := &struct {
				HostIP   string  `json:"host_ip"`
				HostName string  `json:"host_name"`
				CpuNum   float64 `json:"cpu_num"`
				CpuUsage float64 `json:"cur_usage"`
				MemTotal uint64  `json:"mem_total"`
				MemUsage float64 `json:"mem_usage"`
			}{}
			if e = json.Unmarshal(body, tmp); e != nil {
				log.Error(ctx, "[AppInstances] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "ip:", ip, e)
				return e
			}
			info.HostIp = ip
			info.HostName = tmp.HostName
			info.CpuNum = tmp.CpuNum
			info.CpuUsage = tmp.CpuUsage
			info.MemTotal = float64(tmp.MemTotal) / 1024.0 / 1024.0
			info.MemUsage = tmp.MemUsage
			return nil
		})
	}
	if e := egroup.PutGroup(eg); e != nil {
		return nil, ecode.ErrSystem
	}
	return resp, nil
}
func (s *Service) AppInstanceCmd(ctx context.Context, req *api.AppInstanceCmdReq) (*api.AppInstanceCmdResp, error) {
	md := metadata.GetMetadata(ctx)
	s.Lock()
	app, ok := s.apps[req.GName+"."+req.AName]
	if !ok {
		s.Unlock()
		return nil, ecode.ErrAppNotExist
	}
	s.Unlock()
	if e := util.SignCheck(req.Secret, app.Value); e != nil {
		log.Error(ctx, "[AppInstanceCmd] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "host:", req.HostIp, "cmd:", req.Cmd, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	tmpdi := discover.NewDNSDiscover(req.GName, req.AName, req.AName+"-headless."+req.GName, time.Second*10, 9000, 10000, 8000)
	defer tmpdi.Stop()
	notice, cancel := tmpdi.GetNotice()
	defer cancel()
	var ips []string
	select {
	case <-notice:
		addrs, e := tmpdi.GetAddrs(discover.NotNeed)
		if e != nil {
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		ips = make([]string, 0, len(addrs))
		for addr := range addrs {
			ips = append(ips, addr)
		}
	case <-ctx.Done():
		return nil, cerror.ConvertStdError(ctx.Err())
	}
	find := false
	for _, ip := range ips {
		if ip == req.HostIp {
			find = true
			break
		}
	}
	if !find {
		return nil, ecode.ErrAppInstanceNotExist
	}
	switch req.Cmd {
	case "pprof":
		di := discover.NewDirectDiscover(req.GName, req.AName, "http://"+req.HostIp+":6060", 9000, 10000, 8000)
		defer di.Stop()
		webclient, e := web.NewWebClient(dao.GetWebClientConfig(), di, model.Group, model.Name, req.GName, req.AName, nil)
		if e != nil {
			log.Error(ctx, "[AppInstanceCmd] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "ip:", req.HostIp, e)
			return nil, ecode.ErrSystem
		}
		defer webclient.Close(false)
		r, e := webclient.Get(ctx, "/debug/pprof/profile", "", nil, nil)
		if e != nil {
			log.Error(ctx, "[AppInstanceCmd] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "host:", req.HostIp, e)
			return nil, ecode.ErrSystem
		}
		defer r.Body.Close()
		body, e := io.ReadAll(r.Body)
		if e != nil {
			log.Error(ctx, "[AppInstanceCmd] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "host:", req.HostIp, e)
			return nil, ecode.ErrSystem
		}
		return &api.AppInstanceCmdResp{Data: hex.EncodeToString(body)}, nil
	}
	return &api.AppInstanceCmdResp{}, nil
}
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
	nodeidstr, e := s.appDao.MongoCreateApp(ctx, projectid, req.GName, req.AName, req.Secret)
	if e != nil {
		log.Error(ctx, "[CreateApp] group:", req.GName, "app:", req.AName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeid, _ := util.ParseNodeIDstr(nodeidstr)
	log.Info(ctx, "[CreateApp] group:", req.GName, "app:", req.AName, "success")
	return &api.CreateAppResp{NodeId: nodeid}, nil
}

func (s *Service) DelApp(ctx context.Context, req *api.DelAppReq) (*api.DelAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelApp] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
	if e != nil {
		log.Error(ctx, "[DelApp] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		log.Error(ctx, "[DelApp] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
		return nil, ecode.ErrConfigDataBroken
	}
	//self can't be deleted
	if nodeids[1] == "1" && nodeids[3] == "1" {
		return nil, ecode.ErrPermission
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

	//logic
	if e := s.appDao.MongoDelApp(ctx, req.GName, req.AName, req.Secret); e != nil {
		log.Error(ctx, "[DelApp] group:", req.GName, "app:", req.AName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelApp] group:", req.GName, "app:", req.AName, "success")
	return &api.DelAppResp{}, nil
}

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
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[UpdateAppSecret] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
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
	if e := s.appDao.MongoUpdateAppSecret(ctx, req.GName, req.AName, req.OldSecret, req.NewSecret); e != nil {
		log.Error(ctx, "[UpdateAppSecret] group:", req.GName, "app:", req.AName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[UpdateAppSecret] group:", req.GName, "app:", req.AName, "success")
	return &api.UpdateAppSecretResp{}, nil
}

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
	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
	if e != nil {
		log.Error(ctx, "[DelKey] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		log.Error(ctx, "[DelKey] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
		return nil, ecode.ErrConfigDataBroken
	}
	if nodeids[1] == "1" && nodeids[3] == "1" && (req.Key == "AppConfig" || req.Key == "SourceConfig") {
		log.Error(ctx, "[DelKey] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "can't delete key:", req.Key)
		return nil, ecode.ErrPermission
	}
	if !operator.IsZero() {
		//config control permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[DelKey] operator:", md["Token-Data"], "nodeid:", nodeid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	if e := s.appDao.MongoDelKey(ctx, req.GName, req.AName, req.Key, req.Secret); e != nil {
		log.Error(ctx, "[DelKey] group:", req.GName, "app:", req.AName, "key:", req.Key, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelKey] group:", req.GName, "app:", req.AName, "key:", req.Key)
	return &api.DelKeyResp{}, nil
}

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
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[GetKeyConfig] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
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
	keysummary, configlog, e := s.appDao.MongoGetKeyConfig(ctx, req.GName, req.AName, req.Key, req.Index, req.Secret)
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

func (s *Service) SetKeyConfig(ctx context.Context, req *api.SetKeyConfigReq) (*api.SetKeyConfigResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}
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
	switch req.ValueType {
	case "json":
		buf := bytes.NewBuffer(nil)
		if e := json.Compact(buf, common.Str2byte(req.Value)); e != nil {
			log.Error(ctx, "[SetKeyConfig] group:", req.GName, "app:", req.AName, "json value format check failed:", e)
			return nil, ecode.ErrReq
		}
		req.Value = common.Byte2str(buf.Bytes())
	case "toml":
		//TODO
	case "yaml":
		//TODO
	case "raw":
		//TODO
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SetKeyConfig] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[SetKeyConfig] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
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
	index, version, e := s.appDao.MongoSetKeyConfig(ctx, req.GName, req.AName, req.Key, req.Secret, req.Value, req.ValueType, req.NewKey)
	if e != nil {
		log.Error(ctx, "[SetKeyConfig] group:", req.GName, "app:", req.AName, "key:", req.Key, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[SetKeyConfig] group:", req.GName, "app:", req.AName, "key:", req.Key, "index:", index, "version:", version, "success")
	return &api.SetKeyConfigResp{}, nil
}

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
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[Rollback] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
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
	if e := s.appDao.MongoRollbackKeyConfig(ctx, req.GName, req.AName, req.Key, req.Secret, req.Index); e != nil {
		log.Error(ctx, "[Rollback] group:", req.GName, "app:", req.AName, "key:", req.Key, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[Rollback] group:", req.GName, "app:", req.AName, "key:", req.Key, "index:", req.Index, "success")
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
				return nil, cerror.ErrServerClosing
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

func (s *Service) SetProxy(ctx context.Context, req *api.SetProxyReq) (*api.SetProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SetProxy] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
	if e != nil {
		log.Error(ctx, "[SetProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		log.Error(ctx, "[SetProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
		return nil, ecode.ErrConfigDataBroken
	}
	if nodeids[1] == "1" && nodeids[3] == "1" {
		log.Error(ctx, "[SetProxy] operator:", md["Token-Data"], "group:", req.GName, "app:", req.AName, "can't have proxy")
		return nil, ecode.ErrPermission
	}
	if !operator.IsZero() {
		//config control permission check
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
	if req.Path[0] != '/' {
		req.Path = "/" + req.Path
	}
	newnodeidstr, e := s.appDao.MongoSetProxyPath(ctx, req.GName, req.AName, req.Secret, req.Path, req.Read, req.Write, req.Admin, req.NewPath)
	if e != nil {
		log.Error(ctx, "[SetProxy] group:", req.GName, "app:", req.AName, "path:", req.Path, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	newnodeid, _ := util.ParseNodeIDstr(newnodeidstr)
	log.Info(ctx, "[SetProxy] group:", req.GName, "app:", req.AName, "path:", req.Path, "read:", req.Read, "write:", req.Write, "success")
	return &api.SetProxyResp{NodeId: newnodeid}, nil
}
func (s *Service) DelProxy(ctx context.Context, req *api.DelProxyReq) (*api.DelProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelProxy] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, req.GName, req.AName)
	if e != nil {
		log.Error(ctx, "[DelProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid failed:", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		log.Error(ctx, "[DelProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "permission nodeid:", nodeid, "format wrong")
		return nil, ecode.ErrConfigDataBroken
	}
	if nodeids[1] == "1" && nodeids[3] == "1" {
		log.Error(ctx, "[DelProxy] operator:", md["Token-Data"], "get group:", req.GName, "app:", req.AName, "can't have proxy")
		return nil, ecode.ErrPermission
	}
	if !operator.IsZero() {
		//config control permission check
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
	if e := s.appDao.MongoDelProxyPath(ctx, req.GName, req.AName, req.Secret, req.Path); e != nil {
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
	if req.Path[0] != '/' {
		req.Path = "/" + req.Path
	}
	pathinfo, ok := app.Paths[req.Path]
	if !ok {
		s.Unlock()
		return nil, ecode.ErrProxyPathNotExist
	}
	c, ok := s.clients[req.GName+"."+req.AName]
	if !ok {
		di := discover.NewDNSDiscover(req.GName, req.AName, req.AName+"-headless."+req.GName, time.Second*10, 9000, 10000, 8000)
		client, e := crpc.NewCrpcClient(dao.GetCrpcClientConfig(), di, model.Group, model.Name, req.GName, req.AName, nil)
		if e != nil {
			log.Error(ctx, "[Proxy] new crpc client to group:", req.GName, "app:", req.AName, e)
			s.Unlock()
			return nil, ecode.ErrSystem
		}
		c = &clientinfo{
			di:     di,
			client: client,
		}
		s.clients[req.GName+"."+req.AName] = c
	}
	s.clientsActive[req.GName+"."+req.AName] = time.Now().UnixNano()
	s.Unlock()
	if !operator.IsZero() && (pathinfo.PermissionRead || pathinfo.PermissionWrite || pathinfo.PermissionAdmin) {
		//permission check
		canread, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, pathinfo.PermissionNodeID, true)
		if e != nil {
			log.Error(ctx, "[Proxy] operator:", md["Token-Data"], "nodeid:", pathinfo.PermissionNodeID, "get permission failed:", e)
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
