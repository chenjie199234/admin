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
	"github.com/chenjie199234/admin/util"

	cerror "github.com/chenjie199234/Corelib/error"
	"github.com/chenjie199234/Corelib/log"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
)

//Service subservice for config business
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

//Start -
func Start() *Service {
	s := &Service{
		configDao:  configdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		noticepool: &sync.Pool{},
		apps:       make(map[string]map[string]*app),
		status:     true,
	}
	if e := s.configDao.MongoWatchConfig(s.refresh, s.update, s.delgroup, s.delapp, s.delconfig, util.Decrypt); e != nil {
		panic("[Config.Start] watch error: " + e.Error())
	}
	return s
}

//first key groupname,second key appname,value curconfig
func (s *Service) refresh(curs map[string]map[string]*model.AppSummary) {
	s.Lock()
	defer s.Unlock()
	//delete not exist
	for gname, g := range s.apps {
		curg, ok := curs[gname]
		if !ok {
			log.Debug(nil, "[refresh.delgroup] group:", gname)
			for aname, a := range g {
				if len(a.notices) == 0 {
					//if there are no watchers,clean right now
					delete(g, aname)
				} else {
					//if there are watchers,clean will happened when watcher return
					a.appsummary.Cipher = ""
					for _, keysummary := range a.appsummary.Keys {
						keysummary.CurVersion = 0
						keysummary.CurIndex = 0
						keysummary.MaxIndex = 0
						keysummary.CurValue = ""
					}
					for notice := range a.notices {
						notice <- nil
					}
				}
			}
			if len(g) == 0 {
				//if there are no watchers,clean right now
				delete(s.apps, gname)
			}
			continue
		}
		for aname, a := range g {
			if _, ok := curg[aname]; !ok {
				log.Debug(nil, "[refresh.delapp] group:", gname, "app:", aname)
				if len(a.notices) == 0 {
					//if there are no watchers,clean right now
					delete(g, aname)
				} else {
					//if there are watchers,clean will happened when watcher return
					a.appsummary.Cipher = ""
					for _, keysummary := range a.appsummary.Keys {
						keysummary.CurVersion = 0
						keysummary.CurIndex = 0
						keysummary.MaxIndex = 0
						keysummary.CurValue = ""
					}
					for notice := range a.notices {
						notice <- nil
					}
				}
			}
		}
		if len(g) == 0 {
			//if there are no watchers,clean right now
			delete(s.apps, gname)
		}
	}
	//add new or refresh exist
	for gname, curg := range curs {
		g, gok := s.apps[gname]
		if !gok {
			g = make(map[string]*app)
		}
		for aname, cura := range curg {
			log.Debug(nil, "[refresh.update] group:", gname, "app:", aname, "cipher:", cura.Cipher, "keys:", cura.Keys)
			a, ok := g[aname]
			if !ok {
				//this is a new
				if len(cura.Keys) == 0 {
					//this is same as not exist
					continue
				}
				has := false
				for _, keysummary := range cura.Keys {
					if keysummary.CurVersion > 0 {
						has = true
						break
					}

				}
				if !has {
					//this is same as not exist
					continue
				}
				g[aname] = &app{
					appsummary: cura,
					notices:    make(map[chan *struct{}]*struct{}),
				}
			} else {
				//already exist
				if len(cura.Keys) > 0 {
					a.appsummary = cura
					for notice := range a.notices {
						notice <- nil
					}
					continue
				}
				if len(a.notices) == 0 {
					//this is same as not exist and there are no watchers,clean right now
					delete(g, aname)
				}
			}
		}
		if !gok && len(g) > 0 {
			s.apps[gname] = g
		} else if gok && len(g) == 0 {
			delete(s.apps, gname)
		}
	}
}
func (s *Service) update(gname, aname string, cur *model.AppSummary) {
	log.Debug(nil, "[update] group:", gname, "app:", aname, "cipher:", cur.Cipher, "keys:", cur.Keys)
	s.Lock()
	defer s.Unlock()
	g, gok := s.apps[gname]
	if !gok {
		g = make(map[string]*app)
	}
	defer func() {
		if !gok && len(g) > 0 {
			s.apps[gname] = g
		} else if gok && len(g) == 0 {
			delete(s.apps, gname)
		}
	}()
	a, ok := g[aname]
	if !ok {
		//this is a new
		if len(cur.Keys) == 0 {
			//this is same as not exist
			return
		}
		has := false
		for _, keysummary := range cur.Keys {
			if keysummary.CurVersion > 0 {
				break
			}
		}
		if !has {
			//this is same as not exist
			return
		}
		g[aname] = &app{
			appsummary: cur,
			notices:    make(map[chan *struct{}]*struct{}),
		}
		return
	}
	//already exist
	if len(cur.Keys) > 0 {
		a.appsummary = cur
		for notice := range a.notices {
			notice <- nil
		}
		return
	}
	if len(a.notices) == 0 {
		//this is same as not exist and there are no watchers,clean right now
		delete(g, aname)
	}
}
func (s *Service) delgroup(groupname string) {
	log.Debug(nil, "[delgroup] group:", groupname)
	s.Lock()
	defer s.Unlock()
	g, ok := s.apps[groupname]
	if !ok {
		return
	}
	for aname, a := range g {
		if len(a.notices) == 0 {
			//if there are no watchers,clean right now
			delete(g, aname)
		} else {
			//if there are watchers,clean will happened when watcher return
			a.appsummary.Cipher = ""
			for _, keysummary := range a.appsummary.Keys {
				keysummary.CurVersion = 0
				keysummary.CurIndex = 0
				keysummary.MaxIndex = 0
				keysummary.CurValue = ""
			}
			for notice := range a.notices {
				notice <- nil
			}
		}
	}
	if len(g) == 0 {
		//if there are no watchers,clean right now
		delete(s.apps, groupname)
	}
}
func (s *Service) delapp(groupname, appname string) {
	log.Debug(nil, "[delapp] group:", groupname, "app:", appname)
	s.Lock()
	defer s.Unlock()
	g, ok := s.apps[groupname]
	if !ok {
		return
	}
	a, ok := g[appname]
	if !ok {
		return
	}
	if len(a.notices) == 0 {
		//if there are no watchers,clean right now
		delete(g, appname)
		if len(g) == 0 {
			delete(s.apps, groupname)
		}
	} else {
		//if there are watchers,clean will happened when watcher return
		a.appsummary.Cipher = ""
		for _, keysummary := range a.appsummary.Keys {
			keysummary.CurVersion = 0
			keysummary.CurIndex = 0
			keysummary.MaxIndex = 0
			keysummary.CurValue = ""
		}
		for notice := range a.notices {
			notice <- nil
		}
	}
}
func (s *Service) delconfig(groupname, appname, summaryid string) {
	s.Lock()
	defer s.Unlock()
	g, ok := s.apps[groupname]
	if !ok {
		return
	}
	a, ok := g[appname]
	if !ok {
		return
	}
	if a.appsummary.ID.Hex() != summaryid {
		log.Debug(nil, "[delconfig] group:", groupname, "app:", appname, "config log")
		return
	}
	//delete the summary,this is same as delete the app
	log.Debug(nil, "[delconfig] group:", groupname, "app:", appname)
	if len(a.notices) == 0 {
		//if there are no watchers,clean right now
		delete(g, appname)
		if len(g) == 0 {
			delete(s.apps, groupname)
		}
	} else {
		//if there are watchers,clean will happened when watcher return
		a.appsummary.Cipher = ""
		for _, keysummary := range a.appsummary.Keys {
			keysummary.CurVersion = 0
			keysummary.CurIndex = 0
			keysummary.MaxIndex = 0
			keysummary.CurValue = ""
		}
		for notice := range a.notices {
			notice <- nil
		}
	}
}

//get all groups
func (s *Service) Groups(ctx context.Context, req *api.GroupsReq) (*api.GroupsResp, error) {
	groups, e := s.configDao.MongoGetAllGroups(ctx, req.SearchFilter)
	if e != nil {
		log.Error(ctx, "[Groups]", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.GroupsResp{Groups: groups}, nil
}

//del one specific group
func (s *Service) DelGroup(ctx context.Context, req *api.DelGroupReq) (*api.DelGroupResp, error) {
	e := s.configDao.MongoDelGroup(ctx, req.Groupname)
	if e != nil {
		log.Error(ctx, "[DelGroup] group:", req.Groupname, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.DelGroupResp{}, nil
}

//get all apps in one specific group
func (s *Service) Apps(ctx context.Context, req *api.AppsReq) (*api.AppsResp, error) {
	apps, e := s.configDao.MongoGetAllApps(ctx, req.Groupname, req.SearchFilter)
	if e != nil {
		log.Error(ctx, "[Apps] group:", req.Groupname, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.AppsResp{Apps: apps}, nil
}

//del one specific app in one specific group
func (s *Service) DelApp(ctx context.Context, req *api.DelAppReq) (*api.DelAppResp, error) {
	e := s.configDao.MongoDelApp(ctx, req.Groupname, req.Appname)
	if e != nil {
		log.Error(ctx, "[DelApp] group:", req.Groupname, "app:", req.Appname, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.DelAppResp{}, nil
}

//get all config's keys in this app
func (s *Service) Keys(ctx context.Context, req *api.KeysReq) (*api.KeysResp, error) {
	keys, e := s.configDao.MongoGetAllKeys(ctx, req.Groupname, req.Appname)
	if e != nil {
		log.Error(ctx, "[Keys] group:", req.Groupname, "app:", req.Appname, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.KeysResp{Keys: keys}, nil
}

//del one specific key in one specific app
func (s *Service) DelKey(ctx context.Context, req *api.DelKeyReq) (*api.DelKeyResp, error) {
	e := s.configDao.MongoDelKey(ctx, req.Groupname, req.Appname, req.Key)
	if e != nil {
		log.Error(ctx, "[DelKey] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.DelKeyResp{}, nil
}

//create one specific app
func (s *Service) Create(ctx context.Context, req *api.CreateReq) (*api.CreateResp, error) {
	if req.Cipher != "" && len(req.Cipher) != 32 {
		log.Error(ctx, "[Create] group:", req.Groupname, "app:", req.Appname, ecode.ErrCipherLength)
		return nil, ecode.ErrCipherLength
	}
	if e := s.configDao.MongoCreate(ctx, req.Groupname, req.Appname, req.Cipher, util.Encrypt); e != nil {
		log.Error(ctx, "[Create] group:", req.Groupname, "app:", req.Appname, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	log.Info(ctx, "[Create] group:", req.Groupname, "app:", req.Appname, "success")
	return &api.CreateResp{}, nil
}

//update one specific app's cipher
func (s *Service) Updatecipher(ctx context.Context, req *api.UpdatecipherReq) (*api.UpdatecipherResp, error) {
	if req.New != "" && len(req.New) != 32 {
		log.Error(ctx, "[Updatechiper] group:", req.Groupname, "app:", req.Appname, ecode.ErrCipherLength)
		return nil, ecode.ErrCipherLength
	}
	if req.Old == req.New {
		return &api.UpdatecipherResp{}, nil
	}
	if e := s.configDao.MongoUpdateCipher(ctx, req.Groupname, req.Appname, req.Old, req.New, util.Decrypt, util.Encrypt); e != nil {
		log.Error(ctx, "[Updatechiper] group:", req.Groupname, "app:", req.Appname, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	log.Info(ctx, "[Updatecipher] group:", req.GetGroupname, "app:", req.Appname, "success")
	return &api.UpdatecipherResp{}, nil
}

//get config
func (s *Service) Get(ctx context.Context, req *api.GetReq) (*api.GetResp, error) {
	keysummary, configlog, e := s.configDao.MongoGetConfig(ctx, req.Groupname, req.Appname, req.Key, req.Index, util.Decrypt)
	if e != nil {
		log.Error(ctx, "[Get] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, "index:", req.Index, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.GetResp{
		CurIndex:   keysummary.CurIndex,
		MaxIndex:   keysummary.MaxIndex,
		CurVersion: keysummary.CurVersion,
		ThisIndex:  configlog.Index,
		Value:      configlog.Value,
	}, nil
}

//set config
func (s *Service) Set(ctx context.Context, req *api.SetReq) (*api.SetResp, error) {
	req.Key = strings.TrimSpace(req.Key)
	if req.Key == "" {
		log.Error(ctx, "[Set] group:", req.Groupname, "app:", req.Appname, "key empty")
		return nil, ecode.ErrReq
	}
	req.Value = strings.TrimSpace(req.Value)
	if req.Value == "" {
		log.Error(ctx, "[Set] group:", req.Groupname, "app:", req.Appname, "value empty")
		return nil, ecode.ErrReq
	}
	index, e := s.configDao.MongoSetConfig(ctx, req.Groupname, req.Appname, req.Key, req.Value, util.Encrypt)
	if e != nil {
		log.Error(ctx, "[Set] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	log.Info(ctx, "[Set] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, "index:", index, "success")
	return &api.SetResp{}, nil
}

//rollback config
func (s *Service) Rollback(ctx context.Context, req *api.RollbackReq) (*api.RollbackResp, error) {
	if e := s.configDao.MongoRollbackConfig(ctx, req.Groupname, req.Appname, req.Key, req.Index); e != nil {
		log.Error(ctx, "[Rollback] group:", req.Groupname, "app:", req.Appname, "key:", req.Key, e)
		if e != ecode.ErrAppNotExist && e != ecode.ErrIndexNotExist {
			e = ecode.ErrSystem
		}
		return nil, e
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

//watch config
func (s *Service) Watch(ctx context.Context, req *api.WatchReq) (*api.WatchResp, error) {
	resp := &api.WatchResp{
		Datas: make(map[string]*api.WatchData, len(req.Keys)+3),
	}
	s.Lock()
	if !s.status {
		s.Unlock()
		return nil, cerror.ErrClosing
	}
	g, ok := s.apps[req.Groupname]
	if !ok {
		g = make(map[string]*app)
		s.apps[req.Groupname] = g
	}
	a, ok := g[req.Appname]
	if !ok {
		a = &app{
			appsummary: &model.AppSummary{},
			notices:    make(map[chan *struct{}]*struct{}),
		}
	}
	needreturn := false
	for key, clientversion := range req.Keys {
		keysummary, ok := a.appsummary.Keys[key]
		if !ok {
			resp.Datas[key] = &api.WatchData{
				Key:     key,
				Value:   "",
				Version: 0,
			}
			if clientversion != 0 {
				needreturn = true
			}
			continue
		}
		if clientversion != int32(keysummary.CurVersion) {
			resp.Datas[key] = &api.WatchData{
				Key:     key,
				Value:   keysummary.CurValue,
				Version: int32(keysummary.CurVersion),
			}
			needreturn = true
		} else {
			resp.Datas[key] = &api.WatchData{
				Key:     key,
				Value:   "",
				Version: clientversion,
			}
		}
	}
	if needreturn {
		s.Unlock()
		return resp, nil
	}
	if !ok {
		g[req.Appname] = a
	}
	//if int32(a.summary.CurVersion) != req.CurVersion {
	//        resp := &api.WatchResp{
	//                Version:      int32(a.summary.CurVersion),
	//                AppConfig:    a.summary.CurAppConfig,
	//                SourceConfig: a.summary.CurSourceConfig,
	//        }
	//        s.Unlock()
	//        return resp, nil
	//}
	for {
		ch := s.getnotice()
		a.notices[ch] = nil
		s.Unlock()
		select {
		case <-ctx.Done():
			s.Lock()
			delete(a.notices, ch)
			s.putnotice(ch)
			if len(a.notices) == 0 {
				if len(a.appsummary.Keys) == 0 {
					delete(g, req.Appname)
				} else {
					has := false
					for _, keysmmary := range a.appsummary.Keys {
						if keysmmary.CurVersion > 0 {
							has = true
							break
						}
					}
					if !has {
						//this is same as not exist
						delete(g, req.Appname)
					}
				}
			}
			if len(g) == 0 {
				delete(s.apps, req.Groupname)
			}
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
		if int32(a.summary.CurVersion) != req.CurVersion {
			if len(a.notices) == 0 && a.summary.CurVersion == 0 {
				delete(g, req.Appname)
			}
			if len(g) == 0 {
				delete(s.apps, req.Groupname)
			}
			s.Unlock()
			return &api.WatchResp{
				Version:      int32(a.summary.CurVersion),
				AppConfig:    a.summary.CurAppConfig,
				SourceConfig: a.summary.CurSourceConfig,
			}, nil
		}
	}
}

//Stop -
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
