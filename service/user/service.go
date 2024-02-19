package user

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/pool"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/chenjie199234/Corelib/web"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
)

// Service subservice for user business
type Service struct {
	stop *graceful.Graceful

	userDao       *userdao.Dao
	permissionDao *permissiondao.Dao
}

// Start -
func Start() *Service {
	return &Service{
		stop: graceful.New(),

		userDao:       userdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}
}

func (s *Service) GetOauth2(ctx context.Context, req *api.GetOauth2Req) (*api.GetOauth2Resp, error) {
	switch req.SrcType {
	case "DingDing":
		if config.AC.Service.DingDingOauth2 == "" {
			return nil, ecode.ErrBan
		}
		return &api.GetOauth2Resp{Url: config.AC.Service.DingDingOauth2}, nil
	case "FeiShu":
		if config.AC.Service.FeiShuOauth2 == "" {
			return nil, ecode.ErrBan
		}
		return &api.GetOauth2Resp{Url: config.AC.Service.FeiShuOauth2}, nil
	case "WXWork":
		if config.AC.Service.WXWorkOauth2 == "" {
			return nil, ecode.ErrBan
		}
		return &api.GetOauth2Resp{Url: config.AC.Service.WXWorkOauth2}, nil
	}
	log.Error(ctx, "[GetOauth2] unsupported oauth2 type", log.String("type", req.SrcType))
	return nil, ecode.ErrReq
}
func (s *Service) UserLogin(ctx context.Context, req *api.UserLoginReq) (*api.UserLoginResp, error) {
	var userid primitive.ObjectID
	var e error
	var oauth2username, oauth2mobile string
	switch req.SrcType {
	case "DingDing":
		oauth2username, oauth2mobile, e = util.GetDingDingOAuth2(ctx, req.Code)
	case "FeiShu":
		oauth2username, oauth2mobile, e = util.GetFeiShuOAuth2(ctx, req.Code)
	case "WXWork":
		oauth2username, oauth2mobile, e = util.GetWXWorkOAuth2(ctx, req.Code)
	}
	if e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	userid, e = s.userDao.MongoUserLogin(ctx, oauth2mobile, oauth2username, req.SrcType)
	if e != nil {
		log.Error(ctx, "[UserLogin] db op failed", log.String("oauth2_service", req.SrcType), log.String("code", req.Code), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	tokenstr := publicmids.MakeToken(ctx, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, userid.Hex(), "")
	return &api.UserLoginResp{Token: tokenstr}, nil
}
func (s *Service) LoginInfo(ctx context.Context, req *api.LoginInfoReq) (*api.LoginInfoResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[LoginInfo] operator's token format wrong", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	if operator.IsZero() {
		log.Error(ctx, "[LoginInfo] root shouldn't send this request", log.String("operator", md["Token-User"]))
		return &api.LoginInfoResp{User: nil}, nil
	}
	users, e := s.userDao.MongoGetUsers(ctx, []primitive.ObjectID{operator})
	if e != nil {
		log.Error(ctx, "[LoginInfo] db op failed", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	user, ok := users[operator]
	if !ok {
		log.Error(ctx, "[LoginInfo] operator not exist", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrUserNotExist
	}
	respuser := &api.UserInfo{
		UserId:           user.ID.Hex(),
		FeishuUserName:   user.FeiShuUserName,
		DingdingUserName: user.DingDingUserName,
		Ctime:            uint32(user.ID.Timestamp().Unix()),
		ProjectRoles:     make([]*api.ProjectRoles, 0, len(user.Projects)),
	}
	for projecridstr, roles := range user.Projects {
		projectid, e := util.ParseNodeIDstr(projecridstr)
		if e != nil {
			log.Error(ctx, "[LoginInfo] operator's joined project's projectid format wrong",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projecridstr))
			return nil, ecode.ErrSystem
		}
		respuser.ProjectRoles = append(respuser.ProjectRoles, &api.ProjectRoles{ProjectId: projectid, Roles: roles})
	}
	sort.Slice(respuser.ProjectRoles, func(i, j int) bool {
		return respuser.ProjectRoles[i].ProjectId[1] < respuser.ProjectRoles[j].ProjectId[1]
	})
	for _, v := range respuser.ProjectRoles {
		sort.Strings(v.Roles)
	}
	return &api.LoginInfoResp{User: respuser}, nil
}
func (s *Service) InviteProject(ctx context.Context, req *api.InviteProjectReq) (*api.InviteProjectResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[InviteProject] operator's token format wrong", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[InviteProject] target's userid format wrong", log.String("user_id", req.UserId))
		return nil, ecode.ErrReq
	}
	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[InviteProject] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoInvite(ctx, operator, projectid, target); e != nil {
		log.Error(ctx, "[InviteProject] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("user_id", req.UserId),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[InviteProject] success",
		log.String("operator", md["Token-User"]),
		log.String("user_id", req.UserId),
		log.String("project_id", projectid))
	return &api.InviteProjectResp{}, nil
}
func (s *Service) KickProject(ctx context.Context, req *api.KickProjectReq) (*api.KickProjectResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[KickProject] operator's token format wrong", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[KickProject] target's userid format wrong", log.String("user_id", req.UserId))
		return nil, ecode.ErrReq
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, target, projectid, true)
		if e != nil {
			log.Error(ctx, "[KickProject] get target's permission info failed",
				log.String("user_id", req.UserId),
				log.String("project_id", projectid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if admin {
			//target is admin in this project,only root can kick this target from this project
			return nil, ecode.ErrPermission
		}
		//target is not admin in this project
		_, _, admin, e = s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[KickProject] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoKick(ctx, operator, projectid, target); e != nil {
		log.Error(ctx, "[KickProject] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("user_id", req.UserId),
			log.String("project_id", projectid),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[KickProject] success",
		log.String("operator", md["Token-User"]),
		log.String("user_id", req.UserId),
		log.String("project_id", projectid))
	return &api.KickProjectResp{}, nil
}
func (s *Service) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[SearchUsers] operator's token format wrong", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[SearchUsers] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if req.OnlyProject {
			if !canread && !admin {
				return nil, ecode.ErrPermission
			}
		} else if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	searchProjectid := projectid
	if !req.OnlyProject {
		searchProjectid = ""
	}
	users, page, totalsize, e := s.userDao.MongoSearchUsers(ctx, searchProjectid, req.UserName, 20, int64(req.Page))
	if e != nil {
		log.Error(ctx, "[SearchUsers] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("search_user_name", req.UserName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.SearchUsersResp{
		Users:     make([]*api.UserInfo, 0, len(users)),
		Page:      uint32(page),
		Pagesize:  20,
		Totalsize: uint32(totalsize),
	}
	if resp.Page == 0 {
		resp.Pagesize = resp.Totalsize
	}
	//only return the role in the project
	for _, user := range users {
		if user.ID.IsZero() {
			//jump the superadmin
			continue
		}
		respuser := &api.UserInfo{
			UserId:           user.ID.Hex(),
			FeishuUserName:   user.FeiShuUserName,
			DingdingUserName: user.DingDingUserName,
			Ctime:            uint32(user.ID.Timestamp().Unix()),
			ProjectRoles:     make([]*api.ProjectRoles, 0, 1),
		}
		if roles, ok := user.Projects[projectid]; ok {
			respuser.ProjectRoles = append(respuser.ProjectRoles, &api.ProjectRoles{ProjectId: req.ProjectId, Roles: roles})
		}
		for _, v := range respuser.ProjectRoles {
			sort.Strings(v.Roles)
		}
		resp.Users = append(resp.Users, respuser)
	}
	sort.Slice(resp.Users, func(i, j int) bool {
		if resp.Users[i].Ctime == resp.Users[j].Ctime {
			return resp.Users[i].UserId > resp.Users[j].UserId
		}
		return resp.Users[i].Ctime > resp.Users[j].Ctime
	})
	return resp, nil
}

func (s *Service) CreateRole(ctx context.Context, req *api.CreateRoleReq) (*api.CreateRoleResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}
	req.RoleName = strings.TrimSpace(req.RoleName)
	if req.RoleName == "" {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[CreateRole] operator's token format wrong", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[CreateRole] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoCreateRole(ctx, projectid, req.RoleName, req.Comment); e != nil {
		log.Error(ctx, "[CreateRole] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("role_name", req.RoleName),
			log.String("role_comment", req.Comment),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[CreateRole] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("role_name", req.RoleName),
		log.String("role_comment", req.Comment))
	return &api.CreateRoleResp{}, nil
}
func (s *Service) SearchRoles(ctx context.Context, req *api.SearchRolesReq) (*api.SearchRolesResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[SearchRoles] operator's token format wrong", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[SearchRoles] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	roles, page, totalsize, e := s.userDao.MongoSearchRoles(ctx, projectid, req.RoleName, 20, int64(req.Page))
	if e != nil {
		log.Error(ctx, "[SearchRoles] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("search_role_name", req.RoleName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.SearchRolesResp{
		Roles:     make([]*api.RoleInfo, 0, len(roles)),
		Page:      uint32(page),
		Pagesize:  20,
		Totalsize: uint32(totalsize),
	}
	if resp.Page == 0 {
		resp.Pagesize = resp.Totalsize
	}
	for _, role := range roles {
		resp.Roles = append(resp.Roles, &api.RoleInfo{
			ProjectId: req.ProjectId,
			RoleName:  role.RoleName,
			Comment:   role.Comment,
			Ctime:     uint32(role.ID.Timestamp().Unix()),
		})
	}
	sort.Slice(resp.Roles, func(i, j int) bool {
		if resp.Roles[i].Ctime == resp.Roles[j].Ctime {
			return resp.Roles[i].RoleName > resp.Roles[j].RoleName
		}
		return resp.Roles[i].Ctime > resp.Roles[j].Ctime
	})
	return resp, nil
}
func (s *Service) UpdateRole(ctx context.Context, req *api.UpdateRoleReq) (*api.UpdateRoleResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateRole] operator's token format wrong", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	if !operator.IsZero() {
		//permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[UpdateRole] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	oldrole, e := s.userDao.MongoUpdateRole(ctx, projectid, req.RoleName, req.NewComment)
	if e != nil {
		log.Error(ctx, "[UpdateRole] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("role_name", req.RoleName),
			log.String("new_role_comment", req.NewComment),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if oldrole.Comment != req.NewComment {
		log.Info(ctx, "[UpdateRole] success",
			log.String("project_id", oldrole.ProjectID),
			log.String("role_name", oldrole.RoleName),
			log.String("old_role_comment", oldrole.Comment),
			log.String("new_role_comment", req.NewComment))
	} else {
		log.Info(ctx, "[UpdateRole] success,nothing changed",
			log.String("project_id", oldrole.ProjectID),
			log.String("role_name", oldrole.RoleName),
			log.String("role_comment", oldrole.Comment))
	}
	return &api.UpdateRoleResp{}, nil
}
func (s *Service) DelRoles(ctx context.Context, req *api.DelRolesReq) (*api.DelRolesResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[DelRoles] operator's token format wrong", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	//permission check
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[DelRoles] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoDelRoles(ctx, projectid, req.RoleNames); e != nil {
		log.Error(ctx, "[DelRoles] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.Any("role_names", req.RoleNames),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelRoles] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.Any("role_names", req.RoleNames))
	return &api.DelRolesResp{}, nil
}
func (s *Service) AddUserRole(ctx context.Context, req *api.AddUserRoleReq) (*api.AddUserRoleResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[AddUserRole] operator's token format wrong", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[AddUserRole] target's userid format wrong", log.String("user_id", req.UserId))
		return nil, ecode.ErrReq
	}
	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	//permission check
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[AddUserRole] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e = s.userDao.MongoAddUserRole(ctx, target, projectid, req.RoleName); e != nil {
		log.Error(ctx, "[AddUserRole] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("user_id", req.UserId),
			log.String("role_name", req.RoleName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[AddUserRole] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("user_id", req.UserId),
		log.String("role_name", req.RoleName))
	return &api.AddUserRoleResp{}, nil
}
func (s *Service) DelUserRole(ctx context.Context, req *api.DelUserRoleReq) (*api.DelUserRoleResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[DelUserRole] operator's token format wrong", log.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[DelUserRole] target's userid format wrong", log.String("user_id", req.UserId))
		return nil, ecode.ErrReq
	}
	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[DelUserRole] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e = s.userDao.MongoDelUserRole(ctx, target, projectid, req.RoleName); e != nil {
		log.Error(ctx, "[DelUserRole] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("user_id", req.UserId),
			log.String("role_name", req.RoleName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelUserRole] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("user_id", req.UserId),
		log.String("role_name", req.RoleName))
	return &api.DelUserRoleResp{}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
