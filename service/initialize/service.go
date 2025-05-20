package initialize

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	initializedao "github.com/chenjie199234/admin/dao/initialize"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/pool/bpool"
	"github.com/chenjie199234/Corelib/secure"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/graceful"
	"github.com/chenjie199234/Corelib/util/name"
	"go.mongodb.org/mongo-driver/v2/bson"
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
func Start() (*Service, error) {
	return &Service{
		stop: graceful.New(),

		initializeDao: initializedao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		userDao:       userdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}, nil
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
	slog.ErrorContext(ctx, "[InitStatus] db op failed", slog.String("error", e.Error()))
	return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
}

// Init 初始化项目
func (s *Service) Init(ctx context.Context, req *api.InitReq) (*api.InitResp, error) {
	if e := s.initializeDao.MongoInit(ctx, req.Password); e != nil {
		slog.ErrorContext(ctx, "[Init] db op failed", slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.InitResp{}, nil
}

// RootLogin 登录
func (s *Service) RootLogin(ctx context.Context, req *api.RootLoginReq) (*api.RootLoginResp, error) {
	user, e := s.initializeDao.MongoRootLogin(ctx)
	if e != nil {
		slog.ErrorContext(ctx, "[RootLogin] db op failed", slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if e := secure.SignCheck(req.Password, user.Password); e != nil {
		if e == ecode.ErrDataBroken {
			e = ecode.ErrDBDataBroken
		}
		slog.ErrorContext(ctx, "[RootLogin] sign check failed", slog.String("error", e.Error()))
		return nil, e
	}
	tokenstr := publicmids.MakeToken(ctx, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, user.ID.Hex(), "", config.AC.Service.TokenExpire.StdDuration())
	return &api.RootLoginResp{Token: tokenstr}, nil
}

// RootPassword 更新密码
func (s *Service) UpdateRootPassword(ctx context.Context, req *api.UpdateRootPasswordReq) (*api.UpdateRootPasswordResp, error) {
	md := metadata.GetMetadata(ctx)
	//only super admin can change password
	if md["Token-User"] != bson.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
	}
	if e := s.initializeDao.MongoUpdateRootPassword(ctx, req.OldPassword, req.NewPassword); e != nil {
		slog.ErrorContext(ctx, "[UpdateRootPassword] db op failed", slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[UpdateRootPassword] success")
	return &api.UpdateRootPasswordResp{}, nil
}

// CreateProject 创建项目
func (s *Service) CreateProject(ctx context.Context, req *api.CreateProjectReq) (*api.CreateProjectResp, error) {
	if e := name.SingleCheck(req.ProjectName, false); e != nil {
		slog.ErrorContext(ctx, "[CreateProject] project name format wrong", slog.String("project_name", req.ProjectName))
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	//only super admin can create project
	if md["Token-User"] != bson.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
	}
	projectidstr, e := s.initializeDao.MongoCreateProject(ctx, req.ProjectName, req.ProjectData)
	if e != nil {
		slog.ErrorContext(ctx, "[CreateProject] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("name", req.ProjectName),
			slog.String("data", req.ProjectData),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	projectid, _ := util.ParseNodeIDstr(projectidstr)
	slog.InfoContext(ctx, "[CreateProject] success",
		slog.String("project_id", projectidstr),
		slog.String("project_name", req.ProjectName),
		slog.String("project_data", req.ProjectData))
	return &api.CreateProjectResp{ProjectId: projectid}, nil
}

// UpdateProject 更新项目
func (s *Service) UpdateProject(ctx context.Context, req *api.UpdateProjectReq) (*api.UpdateProjectResp, error) {
	//0,1 -> project:admin can't be updated
	if req.ProjectId[0] != 0 || req.ProjectId[1] == 1 {
		return nil, ecode.ErrReq
	}
	if e := name.SingleCheck(req.NewProjectName, false); e != nil {
		slog.ErrorContext(ctx, "[UpdateProject] project name format wrong", slog.String("new_project_name", req.NewProjectName))
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	//only super admin can update project
	if md["Token-User"] != bson.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
	}
	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	oldnode, e := s.initializeDao.MongoUpdateProject(ctx, projectid, req.NewProjectName, req.NewProjectData)
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateProject] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("new_project_name", req.NewProjectName),
			slog.String("new_project_data", req.NewProjectData),
			slog.String("error", e.Error()))
		return nil, e
	}
	if oldnode.NodeName != req.NewProjectName && oldnode.NodeData != req.NewProjectData {
		slog.InfoContext(ctx, "[UpdateProject] success",
			slog.String("project_id", projectid),
			slog.String("old_project_name", oldnode.NodeName),
			slog.String("new_project_name", req.NewProjectName),
			slog.String("old_project_data", oldnode.NodeData),
			slog.String("new_project_data", req.NewProjectData))
	} else if oldnode.NodeName != req.NewProjectName {
		slog.InfoContext(ctx, "[UpdateProject] success",
			slog.String("project_id", projectid),
			slog.String("old_project_name", oldnode.NodeName),
			slog.String("new_project_name", req.NewProjectName))
	} else if oldnode.NodeData != req.NewProjectData {
		slog.InfoContext(ctx, "[UpdateProject] success",
			slog.String("project_id", projectid),
			slog.String("old_project_data", oldnode.NodeData),
			slog.String("new_project_data", req.NewProjectData))
	} else {
		slog.InfoContext(ctx, "[UpdateProject] success,nothing changed",
			slog.String("project_id", projectid),
			slog.String("project_name", oldnode.NodeName),
			slog.String("project_data", oldnode.NodeData))
	}
	return &api.UpdateProjectResp{}, nil
}

// GetProjectIdByName 获取项目id
func (s *Service) GetProjectIdByName(ctx context.Context, req *api.GetProjectIdByNameReq) (*api.GetProjectIdByNameResp, error) {
	projectid, e := s.initializeDao.MongoGetProjectIDByName(ctx, req.ProjectName)
	if e != nil {
		slog.ErrorContext(ctx, "[GetProjectIdByName] db op failed", slog.String("project_name", req.ProjectName), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	projectids, e := util.ParseNodeIDstr(projectid)
	if e != nil {
		slog.ErrorContext(ctx, "[GetProjectIdByName] project's projectid format wrong", slog.String("project_name", req.ProjectName), slog.String("project_id", projectid))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.GetProjectIdByNameResp{ProjectId: projectids}, nil
}

// ListProject 获取项目列表
func (s *Service) ListProject(ctx context.Context, req *api.ListProjectReq) (*api.ListProjectResp, error) {
	md := metadata.GetMetadata(ctx)
	nodes, e := s.initializeDao.MongoListProject(ctx)
	if e != nil {
		slog.ErrorContext(ctx, "[ListProject] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var user *model.User
	if md["Token-User"] != bson.NilObjectID.Hex() {
		operator, e := bson.ObjectIDFromHex(md["Token-User"])
		if e != nil {
			slog.ErrorContext(ctx, "[ListProject] operator's token format wrong", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrToken
		}
		users, e := s.userDao.MongoGetUsers(ctx, []bson.ObjectID{operator})
		if e != nil {
			slog.ErrorContext(ctx, "[ListProject] get operator's user info failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		var ok bool
		user, ok = users[operator]
		if !ok {
			slog.ErrorContext(ctx, "[ListProject] operator not exist", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrSystem
		}
	}
	resp := &api.ListProjectResp{
		Projects: make([]*api.ProjectInfo, 0, len(nodes)),
	}
	for _, node := range nodes {
		if user != nil {
			find := false
			for userprojectid := range user.Projects {
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
			slog.ErrorContext(ctx, "[ListProject] project's projectid format wrong", slog.String("project_id", node.NodeId))
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
	if md["Token-User"] != bson.NilObjectID.Hex() {
		return nil, ecode.ErrPermission
	}
	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	node, e := s.initializeDao.MongoDelProject(ctx, projectid)
	if e != nil {
		slog.ErrorContext(ctx, "[DeleteProject] db op failed", slog.String("operator", md["Token-User"]), slog.String("project_id", projectid), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[DeleteProject] success", slog.String("project_id", projectid), slog.String("project_name", node.NodeName), slog.String("project_data", node.NodeData))
	return &api.DeleteProjectResp{}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
