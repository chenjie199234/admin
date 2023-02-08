package proxy

import (
	"context"
	"sync"
	"time"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/dao"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	proxydao "github.com/chenjie199234/admin/dao/proxy"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/log"
	//"github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service subservice for proxy business
type Service struct {
	stop *graceful.Graceful

	proxyDao      *proxydao.Dao
	permissionDao *permissiondao.Dao

	lker    sync.Mutex
	clients map[string]*crpc.CrpcClient //key:appgroup.appname,value:crpc client
	refresh map[string]int64            //key:appgroup.appname,value:crpc client's last use timestamp(nanosecond)
}

// Start -
func Start() *Service {
	s := &Service{
		stop: graceful.New(),

		proxyDao:      proxydao.NewDao(nil, nil, nil),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),

		clients: make(map[string]*crpc.CrpcClient),
		refresh: make(map[string]int64),
	}
	go s.job()
	return s
}
func (s *Service) job() {
	tker := time.NewTicker(time.Minute)
	for {
		now := <-tker.C
		s.lker.Lock()
		for name, last := range s.refresh {
			if now.UnixNano()-last < time.Minute.Nanoseconds() {
				continue
			}
			delete(s.refresh, name)
			client, ok := s.clients[name]
			if !ok {
				continue
			}
			go client.Close(false)
			delete(s.clients, name)
		}
		s.lker.Unlock()
	}
}
func (s *Service) Tob(ctx context.Context, req *api.TobReq) (*api.TobResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[Tob] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.AdminProjectID+model.Proxy, true)
		if e != nil {
			log.Error(ctx, "[Tob] operator:", md["Token-Data"], "project:", model.AdminProjectID, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	s.lker.Lock()
	client, ok := s.clients[req.Groupname+"."+req.Appname]
	if !ok {
		var e error
		client, e = crpc.NewCrpcClient(dao.GetCrpcClientConfig(), model.Group, model.Name, req.Groupname, req.Appname)
		if e != nil {
			log.Error(ctx, "[Tob] create crpc client for group:", req.Groupname, "app:", req.Appname, "error:", e)
			return nil, ecode.ErrReq
		}
		s.clients[req.Groupname+"."+req.Appname] = client
	}
	s.refresh[req.Groupname+"."+req.Appname] = time.Now().UnixNano()
	s.lker.Unlock()
	respdata, e := client.Call(ctx, req.Path, common.Str2byte(req.Data), metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	return &api.TobResp{Data: common.Byte2str(respdata)}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
