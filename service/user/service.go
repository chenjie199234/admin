package user

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	//"github.com/chenjie199234/Corelib/web"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Service subservice for user business
type Service struct {
	stop *graceful.Graceful

	userDao       *userdao.Dao
	permissionDao *permissiondao.Dao
}

// Start -
func Start() (*Service, error) {
	return &Service{
		stop: graceful.New(),

		userDao:       userdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}, nil
}

func (s *Service) GetOauth2(ctx context.Context, req *api.GetOauth2Req) (*api.GetOauth2Resp, error) {
	switch req.GetSrcType() {
	case "DingDing":
		if config.AC.Service.DingDingOauth2 == "" {
			return nil, ecode.ErrBan
		}
		resp := &api.GetOauth2Resp{}
		resp.SetUrl(config.AC.Service.DingDingOauth2)
		return resp, nil
	case "FeiShu":
		if config.AC.Service.FeiShuOauth2 == "" {
			return nil, ecode.ErrBan
		}
		resp := &api.GetOauth2Resp{}
		resp.SetUrl(config.AC.Service.FeiShuOauth2)
		return resp, nil
	case "WXWork":
		if config.AC.Service.WXWorkOauth2 == "" {
			return nil, ecode.ErrBan
		}
		resp := &api.GetOauth2Resp{}
		resp.SetUrl(config.AC.Service.WXWorkOauth2)
		return resp, nil
	}
	slog.ErrorContext(ctx, "[GetOauth2] unsupported oauth2 type", slog.String("type", req.GetSrcType()))
	return nil, ecode.ErrReq
}
func (s *Service) UserLogin(ctx context.Context, req *api.UserLoginReq) (*api.UserLoginResp, error) {
	var userid bson.ObjectID
	var e error
	var oauth2username, oauth2mobile string
	switch req.GetSrcType() {
	case "DingDing":
		oauth2username, oauth2mobile, e = util.GetDingDingOAuth2(ctx, req.GetCode())
	case "FeiShu":
		oauth2username, oauth2mobile, e = util.GetFeiShuOAuth2(ctx, req.GetCode())
	case "WXWork":
		oauth2username, oauth2mobile, e = util.GetWXWorkOAuth2(ctx, req.GetCode())
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	userid, e = s.userDao.MongoUserLogin(ctx, oauth2mobile, oauth2username, req.GetSrcType())
	if e != nil {
		slog.ErrorContext(ctx, "[UserLogin] db op failed",
			slog.String("oauth2_service", req.GetSrcType()), slog.String("code", req.GetCode()), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	tokenstr := publicmids.MakeToken(ctx, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, userid.Hex(), "", config.AC.Service.TokenExpire.StdDuration())
	resp := &api.UserLoginResp{}
	resp.SetToken(tokenstr)
	return resp, nil
}
func (s *Service) LoginInfo(ctx context.Context, req *api.LoginInfoReq) (*api.LoginInfoResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[LoginInfo] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	if operator.IsZero() {
		slog.ErrorContext(ctx, "[LoginInfo] root shouldn't send this request", slog.String("operator", md["Token-User"]))
		resp := &api.LoginInfoResp{}
		resp.SetUser(nil)
		return resp, nil
	}
	users, e := s.userDao.MongoGetUsers(ctx, []bson.ObjectID{operator})
	if e != nil {
		slog.ErrorContext(ctx, "[LoginInfo] db op failed", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	user, ok := users[operator]
	if !ok {
		slog.ErrorContext(ctx, "[LoginInfo] operator not exist", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrUserNotExist
	}
	respuser := &api.UserInfo{}
	respuser.SetUserId(user.ID.Hex())
	respuser.SetFeishuUserName(user.FeiShuUserName)
	respuser.SetDingdingUserName(user.DingDingUserName)
	respuser.SetWxworkUserName(user.WXWorkUserName)
	respuser.SetCtime(uint32(user.ID.Timestamp().Unix()))
	respuser.SetProjectRoles(make([]*api.ProjectRoles, 0, len(user.Projects)))
	for projecridstr, roles := range user.Projects {
		projectid, e := util.ParseID(projecridstr)
		if e != nil {
			slog.ErrorContext(ctx, "[LoginInfo] operator's joined project's projectid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projecridstr))
			return nil, ecode.ErrSystem
		}
		tmp := &api.ProjectRoles{}
		tmp.SetProjectId(projectid)
		tmp.SetRoles(roles)
		respuser.SetProjectRoles(append(respuser.GetProjectRoles(), tmp))
	}
	sort.Slice(respuser.GetProjectRoles(), func(i, j int) bool {
		return respuser.GetProjectRoles()[i].GetProjectId()[1] < respuser.GetProjectRoles()[j].GetProjectId()[1]
	})
	for _, v := range respuser.GetProjectRoles() {
		sort.Strings(v.GetRoles())
	}
	resp := &api.LoginInfoResp{}
	resp.SetUser(respuser)
	return resp, nil
}
func (s *Service) InviteProject(ctx context.Context, req *api.InviteProjectReq) (*api.InviteProjectResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[InviteProject] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	target, e := bson.ObjectIDFromHex(req.GetUserId())
	if e != nil {
		slog.ErrorContext(ctx, "[InviteProject] target's userid format wrong", slog.String("user_id", req.GetUserId()))
		return nil, ecode.ErrReq
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[InviteProject] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoInvite(ctx, operator, projectid, target); e != nil {
		slog.ErrorContext(ctx, "[InviteProject] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("user_id", req.GetUserId()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[InviteProject] success",
		slog.String("operator", md["Token-User"]),
		slog.String("user_id", req.GetUserId()),
		slog.String("project_id", projectid))
	return &api.InviteProjectResp{}, nil
}
func (s *Service) KickProject(ctx context.Context, req *api.KickProjectReq) (*api.KickProjectResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[KickProject] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	target, e := bson.ObjectIDFromHex(req.GetUserId())
	if e != nil {
		slog.ErrorContext(ctx, "[KickProject] target's userid format wrong", slog.String("user_id", req.GetUserId()))
		return nil, ecode.ErrReq
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, target, projectid, true)
		if e != nil {
			slog.ErrorContext(ctx, "[KickProject] get target's permission info failed",
				slog.String("user_id", req.GetUserId()),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if admin {
			//target is admin in this project,only root can kick this target from this project
			return nil, ecode.ErrPermission
		}
		//target is not admin in this project
		_, _, admin, e = s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[KickProject] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoKick(ctx, operator, projectid, target); e != nil {
		slog.ErrorContext(ctx, "[KickProject] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("user_id", req.GetUserId()),
			slog.String("project_id", projectid),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[KickProject] success",
		slog.String("operator", md["Token-User"]),
		slog.String("user_id", req.GetUserId()),
		slog.String("project_id", projectid))
	return &api.KickProjectResp{}, nil
}
func (s *Service) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[SearchUsers] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[SearchUsers] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if req.GetOnlyProject() {
			if !canread && !admin {
				return nil, ecode.ErrPermission
			}
		} else if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	searchProjectid := projectid
	if !req.GetOnlyProject() {
		searchProjectid = ""
	}
	users, page, totalsize, e := s.userDao.MongoSearchUsers(ctx, searchProjectid, req.GetUserName(), 20, int64(req.GetPage()))
	if e != nil {
		slog.ErrorContext(ctx, "[SearchUsers] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("search_user_name", req.GetUserName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.SearchUsersResp{}
	resp.SetUsers(make([]*api.UserInfo, 0, len(users)))
	resp.SetPage(uint32(page))
	resp.SetPagesize(20)
	resp.SetTotalsize(uint32(totalsize))
	if resp.GetPage() == 0 {
		resp.SetPagesize(resp.GetTotalsize())
	}
	//only return the role in the project
	for _, user := range users {
		if user.ID.IsZero() {
			//jump the superadmin
			continue
		}
		respuser := &api.UserInfo{}
		respuser.SetUserId(user.ID.Hex())
		respuser.SetFeishuUserName(user.FeiShuUserName)
		respuser.SetDingdingUserName(user.DingDingUserName)
		respuser.SetWxworkUserName(user.WXWorkUserName)
		respuser.SetCtime(uint32(user.ID.Timestamp().Unix()))
		respuser.SetProjectRoles(make([]*api.ProjectRoles, 0, 1))
		if roles, ok := user.Projects[projectid]; ok {
			tmp := &api.ProjectRoles{}
			tmp.SetProjectId(req.GetProjectId())
			tmp.SetRoles(roles)
			respuser.SetProjectRoles(append(respuser.GetProjectRoles(), tmp))
		}
		for _, v := range respuser.GetProjectRoles() {
			sort.Strings(v.GetRoles())
		}
		resp.SetUsers(append(resp.GetUsers(), respuser))
	}
	sort.Slice(resp.GetUsers(), func(i, j int) bool {
		if resp.GetUsers()[i].GetCtime() == resp.GetUsers()[j].GetCtime() {
			return resp.GetUsers()[i].GetUserId() > resp.GetUsers()[j].GetUserId()
		}
		return resp.GetUsers()[i].GetCtime() > resp.GetUsers()[j].GetCtime()
	})
	return resp, nil
}

func (s *Service) CreateRole(ctx context.Context, req *api.CreateRoleReq) (*api.CreateRoleResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}
	req.SetRoleName(strings.TrimSpace(req.GetRoleName()))
	if req.GetRoleName() == "" {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[CreateRole] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[CreateRole] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoCreateRole(ctx, projectid, req.GetRoleName(), req.GetComment()); e != nil {
		slog.ErrorContext(ctx, "[CreateRole] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("role_name", req.GetRoleName()),
			slog.String("role_comment", req.GetComment()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[CreateRole] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("role_name", req.GetRoleName()),
		slog.String("role_comment", req.GetComment()))
	return &api.CreateRoleResp{}, nil
}
func (s *Service) SearchRoles(ctx context.Context, req *api.SearchRolesReq) (*api.SearchRolesResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[SearchRoles] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[SearchRoles] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	roles, page, totalsize, e := s.userDao.MongoSearchRoles(ctx, projectid, req.GetRoleName(), 20, int64(req.GetPage()))
	if e != nil {
		slog.ErrorContext(ctx, "[SearchRoles] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("search_role_name", req.GetRoleName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.SearchRolesResp{}
	resp.SetRoles(make([]*api.RoleInfo, 0, len(roles)))
	resp.SetPage(uint32(page))
	resp.SetPagesize(20)
	resp.SetTotalsize(uint32(totalsize))
	if resp.GetPage() == 0 {
		resp.SetPagesize(resp.GetTotalsize())
	}
	for _, role := range roles {
		tmp := &api.RoleInfo{}
		tmp.SetProjectId(req.GetProjectId())
		tmp.SetRoleName(role.RoleName)
		tmp.SetComment(role.Comment)
		tmp.SetCtime(uint32(role.ID.Timestamp().Unix()))
		resp.SetRoles(append(resp.GetRoles(), tmp))
	}
	sort.Slice(resp.GetRoles(), func(i, j int) bool {
		if resp.GetRoles()[i].GetCtime() == resp.GetRoles()[j].GetCtime() {
			return resp.GetRoles()[i].GetRoleName() > resp.GetRoles()[j].GetRoleName()
		}
		return resp.GetRoles()[i].GetCtime() > resp.GetRoles()[j].GetCtime()
	})
	return resp, nil
}
func (s *Service) UpdateRole(ctx context.Context, req *api.UpdateRoleReq) (*api.UpdateRoleResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateRole] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateRole] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	oldrole, e := s.userDao.MongoUpdateRole(ctx, projectid, req.GetRoleName(), req.GetNewComment())
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateRole] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("role_name", req.GetRoleName()),
			slog.String("new_role_comment", req.GetNewComment()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if oldrole.Comment != req.GetNewComment() {
		slog.InfoContext(ctx, "[UpdateRole] success",
			slog.String("project_id", oldrole.ProjectID),
			slog.String("role_name", oldrole.RoleName),
			slog.String("old_role_comment", oldrole.Comment),
			slog.String("new_role_comment", req.GetNewComment()))
	} else {
		slog.InfoContext(ctx, "[UpdateRole] success,nothing changed",
			slog.String("project_id", oldrole.ProjectID),
			slog.String("role_name", oldrole.RoleName),
			slog.String("role_comment", oldrole.Comment))
	}
	return &api.UpdateRoleResp{}, nil
}
func (s *Service) DelRoles(ctx context.Context, req *api.DelRolesReq) (*api.DelRolesResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[DelRoles] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	//permission check
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[DelRoles] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoDelRoles(ctx, projectid, req.GetRoleNames()); e != nil {
		slog.ErrorContext(ctx, "[DelRoles] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.Any("role_names", req.GetRoleNames()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[DelRoles] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.Any("role_names", req.GetRoleNames()))
	return &api.DelRolesResp{}, nil
}
func (s *Service) AddUserRole(ctx context.Context, req *api.AddUserRoleReq) (*api.AddUserRoleResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[AddUserRole] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	target, e := bson.ObjectIDFromHex(req.GetUserId())
	if e != nil {
		slog.ErrorContext(ctx, "[AddUserRole] target's userid format wrong", slog.String("user_id", req.GetUserId()))
		return nil, ecode.ErrReq
	}

	projectid := util.FormID(req.GetProjectId())

	//permission check
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[AddUserRole] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e = s.userDao.MongoAddUserRole(ctx, target, projectid, req.GetRoleName()); e != nil {
		slog.ErrorContext(ctx, "[AddUserRole] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("user_id", req.GetUserId()),
			slog.String("role_name", req.GetRoleName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[AddUserRole] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("user_id", req.GetUserId()),
		slog.String("role_name", req.GetRoleName()))
	return &api.AddUserRoleResp{}, nil
}
func (s *Service) DelUserRole(ctx context.Context, req *api.DelUserRoleReq) (*api.DelUserRoleResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[DelUserRole] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	target, e := bson.ObjectIDFromHex(req.GetUserId())
	if e != nil {
		slog.ErrorContext(ctx, "[DelUserRole] target's userid format wrong", slog.String("user_id", req.GetUserId()))
		return nil, ecode.ErrReq
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[DelUserRole] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e = s.userDao.MongoDelUserRole(ctx, target, projectid, req.GetRoleName()); e != nil {
		slog.ErrorContext(ctx, "[DelUserRole] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("user_id", req.GetUserId()),
			slog.String("role_name", req.GetRoleName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[DelUserRole] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("user_id", req.GetUserId()),
		slog.String("role_name", req.GetRoleName()))
	return &api.DelUserRoleResp{}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
