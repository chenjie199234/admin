package initialize

import (
	"context"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	initializedao "github.com/chenjie199234/admin/dao/initialize"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/util"

	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (s *Service) Init(ctx context.Context, req *api.InitReq) (*api.InitResp, error) {
	if e := s.initializeDao.MongoInit(ctx, req.Password); e != nil {
		log.Error(ctx, "[Init]", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.InitResp{}, nil
}
func (s *Service) RootLogin(ctx context.Context, req *api.RootLoginReq) (*api.RootLoginResp, error) {
	user, e := s.initializeDao.MongoRootLogin(ctx)
	if e != nil {
		log.Error(ctx, "[RootLogin]", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if user.Password != req.Password {
		return nil, ecode.ErrPasswordWrong
	}
	tokenstr := publicmids.MakeToken(ctx, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, user.ID.Hex())
	return &api.RootLoginResp{Token: tokenstr}, nil
}
func (s *Service) RootPassword(ctx context.Context, req *api.RootPasswordReq) (*api.RootPasswordResp, error) {
	if e := s.initializeDao.MongoRootPassword(ctx, req.OldPassword, req.NewPassword); e != nil {
		log.Error(ctx, "[RootPassword]", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.RootPasswordResp{}, nil
}

// 创建项目
func (s *Service) CreateProject(ctx context.Context, req *api.CreateProjectReq) (*api.CreateProjectResp, error) {
	md := metadata.GetMetadata(ctx)
	//only super admin can create project
	if md["Token-Data"] != primitive.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
	}
	if e := s.initializeDao.MongoCreateProject(ctx, req.ProjectName, req.ProjectData); e != nil {
		log.Error(ctx, "[CreateProject] operator:", md["Token-Data"], "Name:", req.ProjectName, "Data:", req.ProjectData, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.CreateProjectResp{}, nil
}

// 获取项目列表
func (s *Service) ListProject(ctx context.Context, req *api.ListProjectReq) (*api.ListProjectResp, error) {
	md := metadata.GetMetadata(ctx)
	nodes, e := s.initializeDao.MongoListProject(ctx)
	if e != nil {
		log.Error(ctx, "[ListProject] operator:", md["Token-Data"], e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.ListProjectResp{
		Projects: make([]*api.ProjectInfo, 0, len(nodes)),
	}
	for _, node := range nodes {
		nodeid, e := util.ParseNodeIDstr(node.NodeId)
		if e != nil {
			log.Error(ctx, "[ListProject] project:", node.NodeId, "format wrong:", e)
			return nil, ecode.ErrSystem
		}
		resp.Projects = append(resp.Projects, &api.ProjectInfo{
			NodeId:      nodeid,
			ProjectName: node.NodeName,
			ProjectData: node.NodeData,
		})
	}
	return resp, nil
}

// Stop -
func (s *Service) Stop() {

}
