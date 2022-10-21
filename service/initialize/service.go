package initialize

import (
	"context"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	initializedao "github.com/chenjie199234/admin/dao/initialize"
	"github.com/chenjie199234/admin/ecode"

	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/log"
	//"github.com/chenjie199234/Corelib/web"
)

// Service subservice for init business
type Service struct {
	initializeDao *initializedao.Dao
}

// Start -
func Start() *Service {
	return &Service{
		initializeDao: initializedao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}
}

func (s *Service) Initialize(ctx context.Context, req *api.InitializeReq) (*api.InitializeResp, error) {
	if e := s.initializeDao.MongoInit(ctx, req.SuperAdminPassword); e != nil {
		log.Error(ctx, "[Initialize]", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.InitializeResp{}, nil
}

// Stop -
func (s *Service) Stop() {

}
