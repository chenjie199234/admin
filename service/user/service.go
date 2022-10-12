package user

import (
	"context"
	"time"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	cerror "github.com/chenjie199234/Corelib/error"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/chenjie199234/Corelib/web"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
)

// Service subservice for user business
type Service struct {
	userDao       *userdao.Dao
	permissionDao *permissiondao.Dao
}

// Start -
func Start() *Service {
	return &Service{
		userDao:       userdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}
}
func (s *Service) SuperAdminLogin(ctx context.Context, req *api.SuperAdminLoginReq) (*api.SuperAdminLoginResp, error) {
	users, e := s.userDao.MongoGetUsers(ctx, []primitive.ObjectID{primitive.NilObjectID})
	if e != nil {
		log.Error(ctx, "[SuperAdminLogin]", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	user, ok := users[primitive.NilObjectID]
	if !ok {
		return nil, ecode.ErrNotInited
	}
	if user.Password != req.Password {
		return nil, ecode.ErrAuth
	}
	start := time.Now()
	end := start.Add(config.AC.TokenExpire.StdDuration())
	tokenstr := publicmids.MakeToken(config.AC.TokenSecret, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, user.ID.Hex(), uint64(start.Unix()), uint64(end.Unix()))
	return &api.SuperAdminLoginResp{Token: tokenstr}, nil
}
func (s *Service) Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error) {
	var userid string
	//TODO login and get userid
	start := time.Now()
	end := start.Add(config.AC.TokenExpire.StdDuration())
	tokenstr := publicmids.MakeToken(config.AC.TokenSecret, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, userid, uint64(start.Unix()), uint64(end.Unix()))
	return &api.LoginResp{Token: tokenstr}, nil
}
func (s *Service) GetUsers(ctx context.Context, req *api.GetUsersReq) (*api.GetUsersResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[GetUsers] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.UserControlNodeId)
	if e != nil {
		log.Error(ctx, "[GetUsers] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !canread && !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	undup := make(map[primitive.ObjectID]*struct{})
	for _, userid := range req.UserIds {
		obj, e := primitive.ObjectIDFromHex(userid)
		if e != nil {
			log.Error(ctx, "[GetUsers] userid:", userid, "format wrong:", e)
			return nil, ecode.ErrReq
		}
		undup[obj] = nil
	}
	userids := make([]primitive.ObjectID, 0, len(undup))
	for userid := range undup {
		userids = append(userids, userid)
	}
	users, e := s.userDao.MongoGetUsers(ctx, userids)
	if e != nil {
		log.Error(ctx, "[GetUsers]", req.UserIds, e)
		return nil, ecode.ErrSystem
	}
	resp := &api.GetUsersResp{
		Users: make([]*api.UserInfo, 0, len(users)),
	}
	for _, user := range users {
		resp.Users = append(resp.Users, &api.UserInfo{
			UserId:     user.ID.Hex(),
			UserName:   user.UserName,
			Department: user.Department,
			Ctime:      user.Ctime,
		})
	}
	return resp, nil
}
func (s *Service) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SearchUsers] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.UserControlNodeId)
	if e != nil {
		log.Error(ctx, "[SearchUsers] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !canread && !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	var users map[primitive.ObjectID]*model.User
	var totalsize int64
	if req.Page == 0 {
		users, totalsize, e = s.userDao.MongoSearchUsers(ctx, req.UserName, 0, 0)
	} else {
		skip := int64(req.Page-1) * 20
		users, totalsize, e = s.userDao.MongoSearchUsers(ctx, req.UserName, 20, skip)
	}
	if e != nil {
		log.Error(ctx, "[SearchUsers]", req.UserName, e)
		return nil, ecode.ErrSystem
	}
	resp := &api.SearchUsersResp{
		Users:     make([]*api.UserInfo, 0, len(users)),
		Page:      req.Page,
		Pagesize:  20,
		Totalsize: uint32(totalsize),
	}
	for _, user := range users {
		resp.Users = append(resp.Users, &api.UserInfo{
			UserId:     user.ID.Hex(),
			UserName:   user.UserName,
			Department: user.Department,
			Ctime:      user.Ctime,
		})
	}
	return resp, nil
}
func (s *Service) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (*api.UpdateUserResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateUser] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.UserControlNodeId)
	if e != nil {
		log.Error(ctx, "[UpdateUser] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !canwrite && !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	userid, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[UpdateUser] userid:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	if e := s.userDao.MongoUpdateUser(ctx, userid, req.UserName, req.Department); e != nil {
		log.Error(ctx, "[UpdateUser]", req.UserId, e)
		return nil, ecode.ErrSystem
	}
	return &api.UpdateUserResp{}, nil
}
func (s *Service) DelUsers(ctx context.Context, req *api.DelUsersReq) (*api.DelUsersResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelUsers] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.UserControlNodeId)
	if e != nil {
		log.Error(ctx, "[DelUsers] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	undup := make(map[primitive.ObjectID]*struct{})
	for _, userid := range req.UserIds {
		if userid == md["Token-Data"] {
			//can't delete self
			continue
		}
		obj, e := primitive.ObjectIDFromHex(userid)
		if e != nil {
			log.Error(ctx, "[DelUsers] userid:", userid, "format wrong:", e)
			return nil, ecode.ErrReq
		}
		undup[obj] = nil
	}
	userids := make([]primitive.ObjectID, 0, len(undup))
	for userid := range undup {
		userids = append(userids, userid)
	}
	if e := s.userDao.MongoDelUsers(ctx, userids); e != nil {
		log.Error(ctx, "[DelUsers]", req.UserIds, e)
		return nil, ecode.ErrSystem
	}
	return &api.DelUsersResp{}, nil
}
func (s *Service) CreateRole(ctx context.Context, req *api.CreateRoleReq) (*api.CreateRoleResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[CreateRole] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.RoleControlNodeId)
	if e != nil {
		log.Error(ctx, "[CreateRole] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	if e := s.userDao.MongoCreateRole(ctx, req.RoleName, req.Comment); e != nil {
		log.Error(ctx, "[CreateRole]", req.RoleName, e)
		return nil, ecode.ErrSystem
	}
	return &api.CreateRoleResp{}, nil
}
func (s *Service) GetRoles(ctx context.Context, req *api.GetRolesReq) (*api.GetRolesResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[GetRoles] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.RoleControlNodeId)
	if e != nil {
		log.Error(ctx, "[GetRoles] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !canread && !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	roles, e := s.userDao.MongoGetRoles(ctx, req.RoleNames)
	if e != nil {
		log.Error(ctx, "[GetRoles]", req.RoleNames, e)
		return nil, ecode.ErrSystem
	}
	resp := &api.GetRolesResp{
		Roles: make([]*api.RoleInfo, 0, len(roles)),
	}
	for _, role := range roles {
		resp.Roles = append(resp.Roles, &api.RoleInfo{
			RoleName: role.RoleName,
			Comment:  role.Comment,
			Ctime:    role.Ctime,
		})
	}
	return resp, nil
}
func (s *Service) SearchRoles(ctx context.Context, req *api.SearchRolesReq) (*api.SearchRolesResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SearchRoles] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.RoleControlNodeId)
	if e != nil {
		log.Error(ctx, "[SearchRoles] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !canread && !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	var roles map[string]*model.Role
	var totalsize int64
	if req.Page == 0 {
		roles, totalsize, e = s.userDao.MongoSearchRoles(ctx, req.RoleName, 0, 0)
	} else {
		skip := int64(req.Page-1) * 20
		roles, totalsize, e = s.userDao.MongoSearchRoles(ctx, req.RoleName, 20, skip)
	}
	if e != nil {
		log.Error(ctx, "[SearchRoles]", req.RoleName, e)
		return nil, ecode.ErrSystem
	}
	resp := &api.SearchRolesResp{
		Roles:     make([]*api.RoleInfo, 0, len(roles)),
		Page:      req.Page,
		Pagesize:  20,
		Totalsize: uint32(totalsize),
	}
	for _, role := range roles {
		resp.Roles = append(resp.Roles, &api.RoleInfo{
			RoleName: role.RoleName,
			Comment:  role.Comment,
			Ctime:    role.Ctime,
		})
	}
	return resp, nil
}
func (s *Service) UpdateRole(ctx context.Context, req *api.UpdateRoleReq) (*api.UpdateRoleResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateRole] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.RoleControlNodeId)
	if e != nil {
		log.Error(ctx, "[UpdateRole] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !canwrite && !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	if e := s.userDao.MongoUpdateRole(ctx, req.RoleName, req.Comment); e != nil {
		log.Error(ctx, "[UpdateRole]", req.RoleName, e)
		return nil, ecode.ErrSystem
	}
	return &api.UpdateRoleResp{}, nil
}
func (s *Service) DelRoles(ctx context.Context, req *api.DelRolesReq) (*api.DelRolesResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelRoles] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.RoleControlNodeId)
	if e != nil {
		log.Error(ctx, "[DelRoles] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	if e := s.userDao.MongoDelRoles(ctx, req.RoleNames); e != nil {
		log.Error(ctx, "[DelRoles]", req.RoleNames, e)
		return nil, ecode.ErrSystem
	}
	return &api.DelRolesResp{}, nil
}

func (s *Service) AddUserRole(ctx context.Context, req *api.AddUserRoleReq) (*api.AddUserRoleResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[AddUserRole] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.UserControlNodeId)
	if e != nil {
		log.Error(ctx, "[AddUserRole] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !admin {
		return nil, ecode.ErrPermission
	}
	canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.RoleControlNodeId)
	if e != nil {
		log.Error(ctx, "[AddUserRole] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !canread && !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	userid, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[AddUserRole] userid:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	if e = s.userDao.MongoAddUserRole(ctx, userid, req.RoleName); e != nil {
		log.Error(ctx, "[AddUserRole]", req.UserId, req.RoleName, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.AddUserRoleResp{}, nil
}
func (s *Service) DelUserRole(ctx context.Context, req *api.DelUserRoleReq) (*api.DelUserRoleResp, error) {
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelUserRole] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.UserControlNodeId)
	if e != nil {
		log.Error(ctx, "[DelUserRole] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !admin {
		return nil, ecode.ErrPermission
	}

	//logic
	userid, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[DelUserRole] userid:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	if e = s.userDao.MongoDelUserRole(ctx, userid, req.RoleName); e != nil {
		log.Error(ctx, "[DelUserRole]", req.UserId, req.RoleName, e)
		return nil, ecode.ErrSystem
	}
	return &api.DelUserRoleResp{}, nil
}

// Stop -
func (s *Service) Stop() {

}
