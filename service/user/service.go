package user

import (
	"context"
	"strings"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/pool"
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

func (s *Service) UserLogin(ctx context.Context, req *api.UserLoginReq) (*api.UserLoginResp, error) {
	var userid string
	//TODO login and get userid
	tokenstr := publicmids.MakeToken(ctx, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, userid)
	return &api.UserLoginResp{Token: tokenstr}, nil
}
func (s *Service) InviteProject(ctx context.Context, req *api.InviteProjectReq) (*api.InviteProjectResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[InviteProject] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[InviteProject] target:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	//user control permission check
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.UserControl, true)
		if e != nil {
			log.Error(ctx, "[InviteProject] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoInvite(ctx, operator, project, target); e != nil {
		log.Error(ctx, "[InviteProject] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.InviteProjectResp{}, nil
}
func (s *Service) KickProject(ctx context.Context, req *api.KickProjectReq) (*api.KickProjectResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[KickProject] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[KickProject] target:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	//user control permission check
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, target, project, true)
		if e != nil {
			log.Error(ctx, "[KickProject] target:", req.UserId, "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if admin {
			//target is admin in this project,only root can kick this target from this project
			return nil, ecode.ErrPermission
		}
		//target is not admin in this project
		_, _, admin, e = s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.UserControl, true)
		if e != nil {
			log.Error(ctx, "[KickProject] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoKick(ctx, operator, project, target); e != nil {
		log.Error(ctx, "[InviteProject] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.KickProjectResp{}, nil
}
func (s *Service) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SearchUsers] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.UserControl, true)
		if e != nil {
			log.Error(ctx, "[SearchUsers] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if req.OnlyProject {
			if !canread && !admin {
				return nil, ecode.ErrPermission
			}
		} else {
			if !admin {
				return nil, ecode.ErrPermission
			}
		}
	}

	//logic
	var users map[primitive.ObjectID]*model.User
	var totalsize int64
	if req.Page == 0 {
		users, totalsize, e = s.userDao.MongoSearchUsers(ctx, project, req.UserName, 0, 0)
	} else {
		skip := int64(req.Page-1) * 20
		users, totalsize, e = s.userDao.MongoSearchUsers(ctx, project, req.UserName, 20, skip)
	}
	if e != nil {
		log.Error(ctx, "[SearchUsers] operator:", md["Token-Data"], "project:", project, "username:", req.UserName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.SearchUsersResp{
		Users:     make([]*api.UserInfo, 0, len(users)),
		Page:      req.Page,
		Pagesize:  20,
		Totalsize: uint32(totalsize),
	}
	if resp.Page == 0 {
		resp.Pagesize = resp.Totalsize
	}
	for _, user := range users {
		undup := make(map[string]*api.UserRoleInfo, len(user.Roles))
		for _, role := range user.Roles {
			index := strings.Index(role, ":")
			roleproject := role[:index]
			rolename := role[index+1:]
			if !req.OnlyProject || roleproject == project {
				if _, ok := undup[roleproject]; !ok {
					undup[roleproject] = &api.UserRoleInfo{
						Project:   roleproject,
						RoleNames: make([]string, 0, 10),
					}
				}
				undup[roleproject].RoleNames = append(undup[roleproject].RoleNames, rolename)
			}
		}
		roles := make([]*api.UserRoleInfo, 0, len(undup))
		for _, v := range undup {
			roles = append(roles, v)
		}
		resp.Users = append(resp.Users, &api.UserInfo{
			UserId:     user.ID.Hex(),
			UserName:   user.UserName,
			Department: user.Department,
			Ctime:      user.Ctime,
			Roles:      roles,
		})
	}
	return resp, nil
}
func (s *Service) CreateRole(ctx context.Context, req *api.CreateRoleReq) (*api.CreateRoleResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[CreateRole] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.RoleControl, true)
		if e != nil {
			log.Error(ctx, "[CreateRole] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoCreateRole(ctx, project, req.RoleName, req.Comment); e != nil {
		log.Error(ctx, "[CreateRole] operator:", md["Token-Data"], "project:", project, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.CreateRoleResp{}, nil
}
func (s *Service) SearchRoles(ctx context.Context, req *api.SearchRolesReq) (*api.SearchRolesResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SearchRoles] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.RoleControl, true)
		if e != nil {
			log.Error(ctx, "[SearchRoles] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	var roles map[string]*model.Role
	var totalsize int64
	if req.Page == 0 {
		roles, totalsize, e = s.userDao.MongoSearchRoles(ctx, project, req.RoleName, 0, 0)
	} else {
		skip := int64(req.Page-1) * 20
		roles, totalsize, e = s.userDao.MongoSearchRoles(ctx, project, req.RoleName, 20, skip)
	}
	if e != nil {
		log.Error(ctx, "[SearchRoles] operator:", md["Token-Data"], "project:", project, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.SearchRolesResp{
		Roles:     make([]*api.RoleInfo, 0, len(roles)),
		Page:      req.Page,
		Pagesize:  20,
		Totalsize: uint32(totalsize),
	}
	if resp.Page == 0 {
		resp.Pagesize = resp.Totalsize
	}
	for _, role := range roles {
		resp.Roles = append(resp.Roles, &api.RoleInfo{
			Project:  req.Project,
			RoleName: role.RoleName,
			Comment:  role.Comment,
			Ctime:    role.Ctime,
		})
	}
	return resp, nil
}
func (s *Service) UpdateRole(ctx context.Context, req *api.UpdateRoleReq) (*api.UpdateRoleResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateRole] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.RoleControl, true)
		if e != nil {
			log.Error(ctx, "[UpdateRole] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoUpdateRole(ctx, project, req.RoleName, req.Comment); e != nil {
		log.Error(ctx, "[UpdateRole] operator:", md["Token-Data"], "project:", project, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.UpdateRoleResp{}, nil
}
func (s *Service) DelRoles(ctx context.Context, req *api.DelRolesReq) (*api.DelRolesResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelRoles] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.RoleControl, true)
		if e != nil {
			log.Error(ctx, "[DelRoles] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoDelRoles(ctx, project, req.RoleNames); e != nil {
		log.Error(ctx, "[DelRoles] operator:", md["Token-Data"], "project:", project, "rolenames:", req.RoleNames, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.DelRolesResp{}, nil
}
func (s *Service) AddUserRole(ctx context.Context, req *api.AddUserRoleReq) (*api.AddUserRoleResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[AddUserRole] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[AddUserRole] target:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.UserControl, true)
		if e != nil {
			log.Error(ctx, "[AddUserRole] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.RoleControl, true)
		if e != nil {
			log.Error(ctx, "[AddUserRole] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e = s.userDao.MongoAddUserRole(ctx, target, project, req.RoleName); e != nil {
		log.Error(ctx, "[AddUserRole] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.AddUserRoleResp{}, nil
}
func (s *Service) DelUserRole(ctx context.Context, req *api.DelUserRoleReq) (*api.DelUserRoleResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelUserRole] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[DelUserRole] target:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.UserControl, true)
		if e != nil {
			log.Error(ctx, "[DelUserRole] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e = s.userDao.MongoDelUserRole(ctx, target, project, req.RoleName); e != nil {
		log.Error(ctx, "[DelUserRole] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.DelUserRoleResp{}, nil
}

// Stop -
func (s *Service) Stop() {

}
