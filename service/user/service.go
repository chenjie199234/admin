package user

import (
	"context"
	"sort"
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
func (s *Service) LoginInfo(ctx context.Context, req *api.LoginInfoReq) (*api.LoginInfoResp, error) {
	md := metadata.GetMetadata(ctx)
	//get other's info,need check permission
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[LoginInfo] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if operator.IsZero() {
		log.Error(ctx, "[LoginInfo] operator:", md["Token-Data"], "is the superadmin,shouldn't send this request")
		return &api.LoginInfoResp{User: nil}, nil
	}
	users, e := s.userDao.MongoGetUsers(ctx, []primitive.ObjectID{operator})
	if e != nil {
		log.Error(ctx, "[LoginInfo] operator:", md["Token-Data"], e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	user, ok := users[operator]
	if !ok {
		log.Error(ctx, "[LoginInfo] operator:", md["Token-Data"], "missing")
		return nil, ecode.ErrSystem
	}
	tmp := make(map[string]*api.ProjectRoles)
	for _, role := range user.Roles {
		index := strings.Index(role, ":")
		roleproject := role[:index]
		roleprojectid, e := util.ParseNodeIDstr(roleproject)
		if e != nil {
			log.Error(ctx, "[LoginInfo] operator:", md["Token-Data"], "role:", role, "projectid format wrong:", e)
			return nil, ecode.ErrSystem
		}
		rolename := role[index+1:]
		tmp[roleproject] = &api.ProjectRoles{
			ProjectId: roleprojectid,
		}
		tmp[roleproject].Roles = append(tmp[roleproject].Roles, rolename)
	}
	for _, project := range user.Projects {
		if _, ok := tmp[project]; ok {
			continue
		}
		id, e := util.ParseNodeIDstr(project)
		if e != nil {
			log.Error(ctx, "[LoginInfo] operator:", md["Token-Data"], "projectid:", project, "format wrong:", e)
			return nil, ecode.ErrSystem
		}
		tmp[project] = &api.ProjectRoles{
			ProjectId: id,
			Roles:     make([]string, 0),
		}
	}
	respuser := &api.UserInfo{
		UserId:       user.ID.Hex(),
		UserName:     user.UserName,
		Department:   user.Department,
		Ctime:        user.Ctime,
		ProjectRoles: make([]*api.ProjectRoles, 0, len(user.Projects)),
	}
	for _, v := range tmp {
		respuser.ProjectRoles = append(respuser.ProjectRoles, v)
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
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()
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
	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[InviteProject] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoInvite(ctx, operator, projectid, target); e != nil {
		log.Error(ctx, "[InviteProject] operator:", md["Token-Data"], "project:", projectid, "target:", req.UserId, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.InviteProjectResp{}, nil
}
func (s *Service) KickProject(ctx context.Context, req *api.KickProjectReq) (*api.KickProjectResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
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
	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, target, projectid, true)
		if e != nil {
			log.Error(ctx, "[KickProject] target:", req.UserId, "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if admin {
			//target is admin in this project,only root can kick this target from this project
			return nil, ecode.ErrPermission
		}
		//target is not admin in this project
		_, _, admin, e = s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[KickProject] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoKick(ctx, operator, projectid, target); e != nil {
		log.Error(ctx, "[InviteProject] operator:", md["Token-Data"], "project:", projectid, "target:", req.UserId, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.KickProjectResp{}, nil
}
func (s *Service) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
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
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SearchUsers] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[SearchUsers] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
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
	var users map[primitive.ObjectID]*model.User
	var totalsize int64
	if req.Page == 0 {
		if req.OnlyProject {
			users, totalsize, e = s.userDao.MongoSearchUsers(ctx, projectid, req.UserName, 0, 0)
		} else {
			users, totalsize, e = s.userDao.MongoSearchUsers(ctx, "", req.UserName, 0, 0)
		}
	} else {
		skip := int64(req.Page-1) * 20
		if req.OnlyProject {
			users, totalsize, e = s.userDao.MongoSearchUsers(ctx, projectid, req.UserName, 20, skip)
		} else {
			users, totalsize, e = s.userDao.MongoSearchUsers(ctx, "", req.UserName, 20, skip)
		}
	}
	if e != nil {
		log.Error(ctx, "[SearchUsers] operator:", md["Token-Data"], "project:", projectid, "username:", req.UserName, e)
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
	//if search user only in the required project,the role should be only in the project too
	for _, user := range users {
		if user.ID.IsZero() {
			//jump the superadmin
			continue
		}
		tmp := make(map[string]*api.ProjectRoles)
		for _, role := range user.Roles {
			index := strings.Index(role, ":")
			roleproject := role[:index]
			if req.OnlyProject && projectid != roleproject {
				continue
			}
			id, e := util.ParseNodeIDstr(roleproject)
			if e != nil {
				log.Error(ctx, "[SearchUsers] operator:", md["Token-Data"], "role:", role, "projectid format wrong:", e)
				return nil, ecode.ErrSystem
			}
			rolename := role[index+1:]
			tmp[roleproject] = &api.ProjectRoles{
				ProjectId: id,
			}
			tmp[roleproject].Roles = append(tmp[roleproject].Roles, rolename)
		}
		for _, project := range user.Projects {
			if _, ok := tmp[project]; ok {
				continue
			}
			if req.OnlyProject && projectid != project {
				continue
			}
			id, e := util.ParseNodeIDstr(project)
			if e != nil {
				log.Error(ctx, "[SearchUsers] operator:", md["Token-Data"], "projectid:", project, "format wrong:", e)
				return nil, ecode.ErrSystem
			}
			tmp[project] = &api.ProjectRoles{
				ProjectId: id,
				Roles:     make([]string, 0),
			}
		}
		respuser := &api.UserInfo{
			UserId:       user.ID.Hex(),
			UserName:     user.UserName,
			Department:   user.Department,
			Ctime:        user.Ctime,
			ProjectRoles: make([]*api.ProjectRoles, 0, len(user.Projects)),
		}
		for _, v := range tmp {
			respuser.ProjectRoles = append(respuser.ProjectRoles, v)
		}
		sort.Slice(respuser.ProjectRoles, func(i, j int) bool {
			return respuser.ProjectRoles[i].ProjectId[1] < respuser.ProjectRoles[j].ProjectId[1]
		})
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
func (s *Service) UpdateUser(ctx context.Context, req *api.UpdateUserReq) (*api.UpdateUserResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateUser] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[UpdateUser] target:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	if !operator.IsZero() {
		//permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.AdminProjectID+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[UpdateUser] operator:", md["Token-Data"], "project:", model.AdminProjectID, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoUpdateUser(ctx, target, req.NewUserName, req.NewDepartment); e != nil {
		log.Error(ctx, "[UpdateUser] operator:", md["Token-Data"], "target:", req.UserId, "new user name:", req.NewUserName, "new department:", req.NewDepartment, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.UpdateUserResp{}, nil
}
func (s *Service) CreateRole(ctx context.Context, req *api.CreateRoleReq) (*api.CreateRoleResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}
	req.RoleName = strings.TrimSpace(req.RoleName)
	if req.RoleName == "" {
		return nil, ecode.ErrReq
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
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[CreateRole] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[CreateRole] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoCreateRole(ctx, projectid, req.RoleName, req.Comment); e != nil {
		log.Error(ctx, "[CreateRole] operator:", md["Token-Data"], "project:", projectid, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.CreateRoleResp{}, nil
}
func (s *Service) SearchRoles(ctx context.Context, req *api.SearchRolesReq) (*api.SearchRolesResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
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
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SearchRoles] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[SearchRoles] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
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
		roles, totalsize, e = s.userDao.MongoSearchRoles(ctx, projectid, req.RoleName, 0, 0)
	} else {
		skip := int64(req.Page-1) * 20
		roles, totalsize, e = s.userDao.MongoSearchRoles(ctx, projectid, req.RoleName, 20, skip)
	}
	if e != nil {
		log.Error(ctx, "[SearchRoles] operator:", md["Token-Data"], "project:", projectid, "rolename:", req.RoleName, e)
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
			ProjectId: req.ProjectId,
			RoleName:  role.RoleName,
			Comment:   role.Comment,
			Ctime:     role.Ctime,
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
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.ProjectId {
		buf.AppendUint32(v)
		if i != len(req.ProjectId)-1 {
			buf.AppendByte(',')
		}
	}
	projectid := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateRole] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[UpdateRole] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoUpdateRole(ctx, projectid, req.RoleName, req.NewComment); e != nil {
		log.Error(ctx, "[UpdateRole] operator:", md["Token-Data"], "project:", projectid, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.UpdateRoleResp{}, nil
}
func (s *Service) DelRoles(ctx context.Context, req *api.DelRolesReq) (*api.DelRolesResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
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
	//permission check
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelRoles] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[DelRoles] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoDelRoles(ctx, projectid, req.RoleNames); e != nil {
		log.Error(ctx, "[DelRoles] operator:", md["Token-Data"], "project:", projectid, "rolenames:", req.RoleNames, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.DelRolesResp{}, nil
}
func (s *Service) AddUserRole(ctx context.Context, req *api.AddUserRoleReq) (*api.AddUserRoleResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
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
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[AddUserRole] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e = s.userDao.MongoAddUserRole(ctx, target, projectid, req.RoleName); e != nil {
		log.Error(ctx, "[AddUserRole] operator:", md["Token-Data"], "project:", projectid, "target:", req.UserId, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.AddUserRoleResp{}, nil
}
func (s *Service) DelUserRole(ctx context.Context, req *api.DelUserRoleReq) (*api.DelUserRoleResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
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
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[DelUserRole] operator:", md["Token-Data"], "project:", projectid, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e = s.userDao.MongoDelUserRole(ctx, target, projectid, req.RoleName); e != nil {
		log.Error(ctx, "[DelUserRole] operator:", md["Token-Data"], "project:", projectid, "target:", req.UserId, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.DelUserRoleResp{}, nil
}

// Stop -
func (s *Service) Stop() {

}
