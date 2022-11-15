package config

import (
	"context"
	"strings"
	"sync"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	configdao "github.com/chenjie199234/admin/dao/config"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/log"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
)

// Service subservice for config business
type Service struct {
	configDao  *configdao.Dao
	noticepool *sync.Pool
	sync.Mutex
	apps   map[string]map[string]*app //first key groupname,second key appname,value appinfo
	status bool
}
type app struct {
	appsummary *model.AppSummary
	notices    map[chan *struct{}]*struct{}
}

// Start -
func Start() *Service {
	s := &Service{
		configDao:  configdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		noticepool: &sync.Pool{},
		apps:       make(map[string]map[string]*app),
		status:     true,
	}
	if e := s.configDao.MongoWatchConfig(s.update, s.delapp, s.delconfig); e != nil {
		panic("[Config.Start] watch error: " + e.Error())
	}
	return s
}

func (s *Service) update(gname, aname string, cur *model.AppSummary) {
	log.Debug(nil, "[update] group:", gname, "app:", aname, "keys:", cur.Keys)
	s.Lock()
	defer s.Unlock()
	if !s.status {
		return
	}
	g, gok := s.apps[gname]
	if !gok {
		return
	}
	a, ok := g[aname]
	if !ok {
		return
	}
	a.appsummary = cur
	for notice := range a.notices {
		select {
		case notice <- nil:
		default:
		}
	}
}
func (s *Service) delapp(groupname, appname string) {
	log.Debug(nil, "[delapp] group:", groupname, "app:", appname)
	s.Lock()
	defer s.Unlock()
	if !s.status {
		return
	}
	g, ok := s.apps[groupname]
	if !ok {
		return
	}
	a, ok := g[appname]
	if !ok {
		return
	}
	a.appsummary = nil
	for notice := range a.notices {
		select {
		case notice <- nil:
		default:
		}
	}
}
func (s *Service) delconfig(groupname, appname, id string) {
	s.Lock()
	defer s.Unlock()
	if !s.status {
		return
	}
	g, ok := s.apps[groupname]
	if !ok {
		return
	}
	a, ok := g[appname]
	if !ok {
		return
	}
	if a.appsummary.ID.Hex() != id {
		log.Debug(nil, "[delconfig] group:", groupname, "app:", appname, "config log")
		return
	}
	//delete the summary,this is same as delete the app
	log.Debug(nil, "[delconfig] group:", groupname, "app:", appname, "summary")
	a.appsummary = nil
	for notice := range a.notices {
		select {
		case notice <- nil:
		default:
		}
	}
}

// get all groups
func (s *Service) Groups(ctx context.Context, req *api.GroupsReq) (*api.GroupsResp, error) {
	groups, e := s.configDao.MongoGetAllGroups(ctx, req.SearchFilter)
	if e != nil {
		log.Error(ctx, "[Groups]", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.GroupsResp{Groups: groups}, nil
}

// get all apps in one specific group
func (s *Service) Apps(ctx context.Context, req *api.AppsReq) (*api.AppsResp, error) {
	apps, e := s.configDao.MongoGetAllApps(ctx, req.Groupname, req.SearchFilter)
	if e != nil {
		log.Error(ctx, "[Apps] group:", req.Groupname, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.AppsResp{Apps: apps}, nil
}

// create one specific app
func (s *Service) CreateApp(ctx context.Context, req *api.CreateAppReq) (*api.CreateAppResp, error) {
	if e := s.configDao.MongoCreateApp(ctx, req.Groupname, req.Appname, req.Secret); e != nil {
		log.Error(ctx, "[CreateApp] group:", req.Groupname, "app:", req.Appname, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[CreateApp] group:", req.Groupname, "app:", req.Appname, "success")
	return &api.CreateAppResp{}, nil
}

// del one specific app in one specific group
func (s *Service) DelApp(ctx context.Context, req *api.DelAppReq) (*api.DelAppResp, error) {
	e := s.configDao.MongoDelApp(ctx, req.Groupname, req.Appname, req.Secret)
	if e != nil {
		log.Error(ctx, "[DelApp] group:", req.Groupname, "app:", req.Appname, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelApp] group:", req.Groupname, "app:", req.Appname, "success")
	return &api.DelAppResp{}, nil
}

// update one specific app's cipher
func (s *Service) UpdateAppSecret(ctx context.Context, req *api.UpdateAppSecretReq) (*api.UpdateAppSecretResp, error) {
	if req.OldSecret == req.NewSecret {
		return &api.UpdateAppSecretResp{}, nil
	}
	if e := s.configDao.MongoUpdateAppSecret(ctx, req.Groupname, req.Appname, req.OldSecret, req.NewSecret); e != nil {
		log.Error(ctx, "[UpdateAppSecret] group:", req.Groupname, "app:", req.Appname, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[UpdateAppSecret] group:", req.Groupname, "app:", req.Appname, "success")
	return &api.UpdateAppSecretResp{}, nil
}

// get all config's keys in this app
func (s *Service) Keys(ctx context.Context, req *api.KeysReq) (*api.KeysResp, error) {
	keys, e := s.configDao.MongoGetAllKeys(ctx, req.Groupname, req.Appname, req.Secret)
	if e != nil {
		log.Error(ctx, "[Keys] group:", req.Groupname, "app:", req.Appname, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.KeysResp{Keys: keys}, nil
}

// del one specific key in one specific app
func (s *Service) DelKey(ctx context.Context, req *api.DelKeyReq) (*api.DelKeyResp, error) {
	e := s.configDao.MongoDelKey(ctx, req.Groupname, req.Appname, req.Key, req.Secret)
	if e != nil {
		log.Error(ctx, "[DelKey] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelKey] group:", req.Groupname, "app:", req.Appname, "key:", req.Key)
	return &api.DelKeyResp{}, nil
}

// get config
func (s *Service) GetKeyConfig(ctx context.Context, req *api.GetKeyConfigReq) (*api.GetKeyConfigResp, error) {
	keysummary, configlog, e := s.configDao.MongoGetKeyConfig(ctx, req.Groupname, req.Appname, req.Key, req.Index, req.Secret)
	if e != nil {
		log.Error(ctx, "[GetKeyConfig] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, "index:", req.Index, e)
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
	req.Key = strings.TrimSpace(req.Key)
	if req.Key == "" {
		log.Error(ctx, "[SetKeyConfig] group:", req.Groupname, "app:", req.Appname, "key empty")
		return nil, ecode.ErrReq
	}
	req.Value = strings.TrimSpace(req.Value)
	if req.Value == "" {
		log.Error(ctx, "[SetKeyConfig] group:", req.Groupname, "app:", req.Appname, "value empty")
		return nil, ecode.ErrReq
	}
	index, version, e := s.configDao.MongoSetKeyConfig(ctx, req.Groupname, req.Appname, req.Key, req.Secret, req.Value, req.ValueType)
	if e != nil {
		log.Error(ctx, "[SetKeyConfig] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[SetKeyConfig] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, "index:", index, "version:", version, "success")
	return &api.SetKeyConfigResp{}, nil
}

// rollback config
func (s *Service) Rollback(ctx context.Context, req *api.RollbackReq) (*api.RollbackResp, error) {
	if e := s.configDao.MongoRollbackKeyConfig(ctx, req.Groupname, req.Appname, req.Key, req.Secret, req.Index); e != nil {
		log.Error(ctx, "[Rollback] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[Rollback] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, "index:", req.Index, "success")
	return &api.RollbackResp{}, nil
}

func (s *Service) getnotice() chan *struct{} {
	ch, ok := s.noticepool.Get().(chan *struct{})
	if !ok {
		return make(chan *struct{}, 1)
	}
	return ch
}
func (s *Service) putnotice(ch chan *struct{}) {
	s.noticepool.Put(ch)
}

// watch config
func (s *Service) Watch(ctx context.Context, req *api.WatchReq) (*api.WatchResp, error) {
	resp := &api.WatchResp{
		Datas: make(map[string]*api.WatchData, len(req.Keys)+3),
	}
	s.Lock()
	if !s.status {
		s.Unlock()
		return nil, cerror.ErrClosing
	}
	g, gok := s.apps[req.Groupname]
	if !gok {
		g = make(map[string]*app)
		s.apps[req.Groupname] = g
	}
	a, aok := g[req.Appname]
	if !aok {
		a = &app{notices: make(map[chan *struct{}]*struct{})}
		g[req.Appname] = a
	}
	if !gok || !aok {
		//lazy init
		appsummary, e := s.configDao.MongoGetAppConfig(ctx, req.Groupname, req.Appname)
		if e != nil {
			s.Unlock()
			log.Error(ctx, "[Watch] group:", req.Groupname, "app:", req.Appname, e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		a.appsummary = appsummary
	}
	needreturn := false
	for key, clientversion := range req.Keys {
		if a.appsummary == nil {
			if clientversion != 0 {
				needreturn = true
			}
			resp.Datas[key] = &api.WatchData{
				Key:       key,
				Value:     "",
				ValueType: "raw",
				Version:   0,
			}
		} else if keysummary, ok := a.appsummary.Keys[key]; !ok {
			if clientversion != 0 {
				needreturn = true
			}
			resp.Datas[key] = &api.WatchData{
				Key:       key,
				Value:     "",
				ValueType: "raw",
				Version:   0,
			}
		} else if clientversion != int32(keysummary.CurVersion) {
			needreturn = true
			resp.Datas[key] = &api.WatchData{
				Key:       key,
				Value:     keysummary.CurValue,
				ValueType: keysummary.CurValueType,
				Version:   int32(keysummary.CurVersion),
			}
		} else {
			resp.Datas[key] = &api.WatchData{
				Key:       key,
				Value:     "",
				ValueType: "",
				Version:   clientversion,
			}
		}
	}
	if needreturn {
		s.Unlock()
		return resp, nil
	}
	for {
		ch := s.getnotice()
		a.notices[ch] = nil
		s.Unlock()
		select {
		case <-ctx.Done():
			s.Lock()
			delete(a.notices, ch)
			s.putnotice(ch)
			s.Unlock()
			return nil, ctx.Err()
		case _, ok := <-ch:
			if !ok {
				return nil, cerror.ErrClosing
			}
		}
		s.Lock()
		delete(a.notices, ch)
		s.putnotice(ch)
		for key, respdata := range resp.Datas {
			if a.appsummary == nil {
				if respdata.Version != 0 {
					needreturn = true
				}
				respdata.Value = ""
				respdata.ValueType = "raw"
				respdata.Version = 0
			} else if keysummary, ok := a.appsummary.Keys[key]; !ok {
				if respdata.Version != 0 {
					needreturn = true
				}
				respdata.Value = ""
				respdata.ValueType = "raw"
				respdata.Version = 0
				continue
			} else if int32(keysummary.CurVersion) != respdata.Version {
				needreturn = true
				respdata.Value = keysummary.CurValue
				respdata.ValueType = keysummary.CurValueType
				respdata.Version = int32(keysummary.CurVersion)
			} else {
				respdata.Value = ""
				respdata.ValueType = ""
			}
		}
		if needreturn {
			s.Unlock()
			return resp, nil
		}
	}
}

// Stop -
func (s *Service) Stop() {
	s.Lock()
	defer s.Unlock()
	s.status = false
	for _, g := range s.apps {
		for _, a := range g {
			for n := range a.notices {
				close(n)
			}
		}
	}
}
