package initialize

import (
	"context"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	initializedao "github.com/chenjie199234/admin/dao/initialize"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/pool"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service subservice for init business
type Service struct {
	initializeDao *initializedao.Dao
	userDao       *userdao.Dao
}

// Start -
func Start() *Service {
	return &Service{
		initializeDao: initializedao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		userDao:       userdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}
}

// Init 初始化项目
func (s *Service) Init(ctx context.Context, req *api.InitReq) (*api.InitResp, error) {
	if e := s.initializeDao.MongoInit(ctx, req.Password); e != nil {
		log.Error(ctx, "[Init]", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.InitResp{}, nil
}

// RootLogin 登录
func (s *Service) RootLogin(ctx context.Context, req *api.RootLoginReq) (*api.RootLoginResp, error) {
	user, e := s.initializeDao.MongoRootLogin(ctx)
	if e != nil {
		log.Error(ctx, "[RootLogin]", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if e := util.SignCheck(req.Password, user.Password); e != nil {
		return nil, ecode.ErrPasswordWrong
	}
	tokenstr := publicmids.MakeToken(ctx, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, user.ID.Hex())
	return &api.RootLoginResp{Token: tokenstr}, nil
}

// RootPassword 更新密码
func (s *Service) RootPassword(ctx context.Context, req *api.RootPasswordReq) (*api.RootPasswordResp, error) {
	if e := s.initializeDao.MongoRootPassword(ctx, req.OldPassword, req.NewPassword); e != nil {
		log.Error(ctx, "[RootPassword]", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.RootPasswordResp{}, nil
}

// CreateProject 创建项目
func (s *Service) CreateProject(ctx context.Context, req *api.CreateProjectReq) (*api.CreateProjectResp, error) {
	md := metadata.GetMetadata(ctx)
	//only super admin can create project
	if md["Token-Data"] != primitive.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
	}
	str, e := s.initializeDao.MongoCreateProject(ctx, req.ProjectName, req.ProjectData)
	if e != nil {
		log.Error(ctx, "[CreateProject] operator:", md["Token-Data"], "Name:", req.ProjectName, "Data:", req.ProjectData, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	projectid, _ := util.ParseNodeIDstr(str)
	return &api.CreateProjectResp{ProjectId: projectid}, nil
}

// UpdateProject 更新项目
func (s *Service) UpdateProject(ctx context.Context, req *api.UpdateProjectReq) (*api.UpdateProjectResp, error) {
	md := metadata.GetMetadata(ctx)
	//only super admin can update project
	if md["Token-Data"] != primitive.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
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
	if e := s.initializeDao.MongoUpdateProject(ctx, projectid, req.NewProjectName, req.NewProjectData); e != nil {
		return nil, e
	}
	return &api.UpdateProjectResp{}, nil
}

// ListProject 获取项目列表
func (s *Service) ListProject(ctx context.Context, req *api.ListProjectReq) (*api.ListProjectResp, error) {
	md := metadata.GetMetadata(ctx)
	nodes, e := s.initializeDao.MongoListProject(ctx)
	if e != nil {
		log.Error(ctx, "[ListProject] operator:", md["Token-Data"], e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var user *model.User
	if md["Token-Data"] != primitive.NilObjectID.Hex() {
		operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
		if e != nil {
			log.Error(ctx, "[ListProject] operator:", md["Token-Data"], "format wrong:", e)
			return nil, ecode.ErrToken
		}
		users, e := s.userDao.MongoGetUsers(ctx, []primitive.ObjectID{operator})
		if e != nil {
			log.Error(ctx, "[ListProject] operator:", md["Token-Data"], "get user info:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		var ok bool
		user, ok = users[operator]
		if !ok {
			log.Error(ctx, "[ListProject] operator:", md["Token-Data"], "doesn't exist")
			return nil, ecode.ErrSystem
		}
	}
	resp := &api.ListProjectResp{
		Projects: make([]*api.ProjectInfo, 0, len(nodes)),
	}
	for _, node := range nodes {
		if user != nil {
			find := false
			for _, projectid := range user.Projects {
				if projectid == node.NodeId {
					find = true
					break
				}
			}
			if !find {
				continue
			}
		}
		nodeid, e := util.ParseNodeIDstr(node.NodeId)
		if e != nil {
			log.Error(ctx, "[ListProject] project:", node.NodeId, "format wrong:", e)
			return nil, ecode.ErrSystem
		}
		resp.Projects = append(resp.Projects, &api.ProjectInfo{
			ProjectId:   nodeid,
			ProjectName: node.NodeName,
			ProjectData: node.NodeData,
		})
	}
	return resp, nil
}

// DeleteProject 删除项目
func (s *Service) DeleteProject(ctx context.Context, req *api.DeleteProjectReq) (*api.DeleteProjectResp, error) {
	//0,1 -> project:admin can't be deleted
	if req.ProjectId[0] != 0 || req.ProjectId[1] == 1 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	//only super admin can create project
	if md["Token-Data"] != primitive.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
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
	if e := s.initializeDao.MongoDelProject(ctx, projectid); e != nil {
		log.Error(ctx, "[DeleteProject] operator:", md["Token-Data"], "Name:", projectid, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.DeleteProjectResp{}, nil
}

// Stop -
func (s *Service) Stop() {

}
