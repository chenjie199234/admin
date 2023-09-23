package initialize

import (
	"context"
	"strconv"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	initializedao "github.com/chenjie199234/admin/dao/initialize"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/pool"
	"github.com/chenjie199234/Corelib/secure"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/graceful"
	"github.com/chenjie199234/Corelib/util/name"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
)

// Service subservice for init business
type Service struct {
	stop *graceful.Graceful

	initializeDao *initializedao.Dao
	userDao       *userdao.Dao
}

// Start -
func Start() *Service {
	return &Service{
		stop: graceful.New(),

		initializeDao: initializedao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		userDao:       userdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}
}

// 初始化状态
func (s *Service) InitStatus(ctx context.Context, req *api.InitStatusReq) (*api.InitStatusResp, error) {
	_, e := s.initializeDao.MongoRootLogin(ctx)
	if e == nil {
		return &api.InitStatusResp{Status: true}, nil
	}
	if e == ecode.ErrNotInited {
		return &api.InitStatusResp{Status: false}, nil
	}
	log.Error(ctx, "[InitStatus] db op failed", log.CError(e))
	return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
}

// Init 初始化项目
func (s *Service) Init(ctx context.Context, req *api.InitReq) (*api.InitResp, error) {
	if e := s.initializeDao.MongoInit(ctx, req.Password); e != nil {
		log.Error(ctx, "[Init] db op failed", log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.InitResp{}, nil
}

// RootLogin 登录
func (s *Service) RootLogin(ctx context.Context, req *api.RootLoginReq) (*api.RootLoginResp, error) {
	user, e := s.initializeDao.MongoRootLogin(ctx)
	if e != nil {
		log.Error(ctx, "[RootLogin] db op failed", log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if e := secure.SignCheck(req.Password, user.Password); e != nil {
		if e == ecode.ErrDataBroken {
			e = ecode.ErrDBDataBroken
		}
		log.Error(ctx, "[RootLogin] sign check failed", log.CError(e))
		return nil, e
	}
	tokenstr := publicmids.MakeToken(ctx, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, user.ID.Hex(), "")
	return &api.RootLoginResp{Token: tokenstr}, nil
}

// RootPassword 更新密码
func (s *Service) UpdateRootPassword(ctx context.Context, req *api.UpdateRootPasswordReq) (*api.UpdateRootPasswordResp, error) {
	md := metadata.GetMetadata(ctx)
	//only super admin can change password
	if md["Token-User"] != primitive.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
	}
	if e := s.initializeDao.MongoUpdateRootPassword(ctx, req.OldPassword, req.NewPassword); e != nil {
		log.Error(ctx, "[RootPassword] db op failed", log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.UpdateRootPasswordResp{}, nil
}

// CreateProject 创建项目
func (s *Service) CreateProject(ctx context.Context, req *api.CreateProjectReq) (*api.CreateProjectResp, error) {
	if e := name.SingleCheck(req.ProjectName, false); e != nil {
		log.Error(ctx, "[CreateProject] project name format wrong", log.String("project_name", req.ProjectName))
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	//only super admin can create project
	if md["Token-User"] != primitive.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
	}
	str, e := s.initializeDao.MongoCreateProject(ctx, req.ProjectName, req.ProjectData)
	if e != nil {
		log.Error(ctx, "[CreateProject] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("name", req.ProjectName),
			log.String("data", req.ProjectData),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	projectid, _ := util.ParseNodeIDstr(str)
	return &api.CreateProjectResp{ProjectId: projectid}, nil
}

// UpdateProject 更新项目
func (s *Service) UpdateProject(ctx context.Context, req *api.UpdateProjectReq) (*api.UpdateProjectResp, error) {
	//0,1 -> project:admin can't be updated
	if req.ProjectId[0] != 0 || req.ProjectId[1] == 1 {
		return nil, ecode.ErrReq
	}
	if e := name.SingleCheck(req.NewProjectName, false); e != nil {
		log.Error(ctx, "[UpdateProject] project name format wrong", log.String("project_name", req.NewProjectName))
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	//only super admin can update project
	if md["Token-User"] != primitive.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
	}
	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	if e := s.initializeDao.MongoUpdateProject(ctx, projectid, req.NewProjectName, req.NewProjectData); e != nil {
		log.Error(ctx, "[UpdateProject] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("new_name", req.NewProjectName),
			log.String("new_data", req.NewProjectData),
			log.CError(e))
		return nil, e
	}
	return &api.UpdateProjectResp{}, nil
}

// GetProjectIdByName 获取项目id
func (s *Service) GetProjectIdByName(ctx context.Context, req *api.GetProjectIdByNameReq) (*api.GetProjectIdByNameResp, error) {
	projectid, e := s.initializeDao.MongoGetProjectIDByName(ctx, req.ProjectName)
	if e != nil {
		log.Error(ctx, "[GetProjectIdByName] db op failed", log.String("project_name", req.ProjectName), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	projectids, e := util.ParseNodeIDstr(projectid)
	if e != nil {
		log.Error(ctx, "[GetProjectIdByName] project's projectid format wrong",
			log.String("project_name", req.ProjectName),
			log.String("project_id", projectid),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.GetProjectIdByNameResp{ProjectId: projectids}, nil
}

// ListProject 获取项目列表
func (s *Service) ListProject(ctx context.Context, req *api.ListProjectReq) (*api.ListProjectResp, error) {
	md := metadata.GetMetadata(ctx)
	nodes, e := s.initializeDao.MongoListProject(ctx)
	if e != nil {
		log.Error(ctx, "[ListProject] db op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var user *model.User
	if md["Token-User"] != primitive.NilObjectID.Hex() {
		operator, e := primitive.ObjectIDFromHex(md["Token-User"])
		if e != nil {
			log.Error(ctx, "[ListProject] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
			return nil, ecode.ErrToken
		}
		users, e := s.userDao.MongoGetUsers(ctx, []primitive.ObjectID{operator})
		if e != nil {
			log.Error(ctx, "[ListProject] get operator's user info failed", log.String("operator", md["Token-User"]), log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		var ok bool
		user, ok = users[operator]
		if !ok {
			log.Error(ctx, "[ListProject] operator not exist", log.String("operator", md["Token-User"]))
			return nil, ecode.ErrSystem
		}
	}
	resp := &api.ListProjectResp{
		Projects: make([]*api.ProjectInfo, 0, len(nodes)),
	}
	for _, node := range nodes {
		if user != nil {
			find := false
			for _, userprojectid := range user.ProjectIDs {
				if userprojectid == node.NodeId {
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
			log.Error(ctx, "[ListProject] project's projectid format wrong", log.String("project_id", node.NodeId), log.CError(e))
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
	if md["Token-User"] != primitive.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
	}
	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)
	if e := s.initializeDao.MongoDelProject(ctx, projectid); e != nil {
		log.Error(ctx, "[DeleteProject] db op failed", log.String("operator", md["Token-User"]), log.String("project_id", projectid), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.DeleteProjectResp{}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
