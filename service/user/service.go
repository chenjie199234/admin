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
		log.Error(ctx, "[LoginInfo] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	if operator.IsZero() {
		log.Error(ctx, "[LoginInfo] superadmin shouldn't send this request", map[string]interface{}{"operator": md["Token-Data"]})
		return &api.LoginInfoResp{User: nil}, nil
	}
	users, e := s.userDao.MongoGetUsers(ctx, []primitive.ObjectID{operator})
	if e != nil {
		log.Error(ctx, "[LoginInfo] db op failed", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	user, ok := users[operator]
	if !ok {
		log.Error(ctx, "[LoginInfo] operator not exist", map[string]interface{}{"operator": md["Token-Data"]})
		return nil, ecode.ErrSystem
	}
	tmp := make(map[string]*api.ProjectRoles)
	//join projects
	for _, userprojectid := range user.ProjectIDs {
		if _, ok := tmp[userprojectid]; ok {
			continue
		}
		id, e := util.ParseNodeIDstr(userprojectid)
		if e != nil {
			log.Error(ctx, "[LoginInfo] operator join project's projectid format wrong", map[string]interface{}{"operator": md["Token-Data"], "project_id": userprojectid, "error": e})
			return nil, ecode.ErrSystem
		}
		tmp[userprojectid] = &api.ProjectRoles{
			ProjectId: id,
			Roles:     make([]string, 0),
		}
	}
	//roles
	for _, role := range user.Roles {
		index := strings.Index(role, ":")
		roleprojectid := role[:index]
		roleproject, e := util.ParseNodeIDstr(roleprojectid)
		if e != nil {
			log.Error(ctx, "[LoginInfo] operator's role's projectid format wrong", map[string]interface{}{"operator": md["Token-Data"], "role": role, "project_id": roleprojectid, "error": e})
			return nil, ecode.ErrSystem
		}
		rolename := role[index+1:]
		if _, ok := tmp[roleprojectid]; !ok {
			tmp[roleprojectid] = &api.ProjectRoles{
				ProjectId: roleproject,
			}
		}
		tmp[roleprojectid].Roles = append(tmp[roleprojectid].Roles, rolename)
	}
	respuser := &api.UserInfo{
		UserId:       user.ID.Hex(),
		UserName:     user.UserName,
		Department:   user.Department,
		Ctime:        uint32(user.ID.Timestamp().Unix()),
		ProjectRoles: make([]*api.ProjectRoles, 0, len(user.ProjectIDs)),
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
		log.Error(ctx, "[InviteProject] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[InviteProject] target's userid format wrong", map[string]interface{}{"target": req.UserId, "error": e})
		return nil, ecode.ErrReq
	}
	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[InviteProject] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoInvite(ctx, operator, projectid, target); e != nil {
		log.Error(ctx, "[InviteProject] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "target": req.UserId, "error": e})
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
		log.Error(ctx, "[KickProject] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[KickProject] target's userid format wrong", map[string]interface{}{"target": req.UserId, "error": e})
		return nil, ecode.ErrReq
	}
	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, target, projectid, true)
		if e != nil {
			log.Error(ctx, "[KickProject] get target's permission info failed", map[string]interface{}{"target": req.UserId, "project_id": projectid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if admin {
			//target is admin in this project,only root can kick this target from this project
			return nil, ecode.ErrPermission
		}
		//target is not admin in this project
		_, _, admin, e = s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[KickProject] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoKick(ctx, operator, projectid, target); e != nil {
		log.Error(ctx, "[InviteProject] db op failed", map[string]interface{}{"operator": md["Token-Data"], "target": req.UserId, "project_id": projectid, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.KickProjectResp{}, nil
}
func (s *Service) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[SearchUsers] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
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

	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[SearchUsers] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "error": e})
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
	users, page, totalsize, e := s.userDao.MongoSearchUsers(ctx, projectid, req.UserName, 20, int64(req.Page))
	if e != nil {
		log.Error(ctx, "[SearchUsers] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "username": req.UserName, "error": e})
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
		tmp := make(map[string]*api.ProjectRoles)
		for _, role := range user.Roles {
			index := strings.Index(role, ":")
			roleproject := role[:index]
			if projectid != roleproject {
				continue
			}
			rolename := role[index+1:]
			if _, ok := tmp[roleproject]; !ok {
				tmp[roleproject] = &api.ProjectRoles{
					ProjectId: req.ProjectId,
				}
			}
			tmp[roleproject].Roles = append(tmp[roleproject].Roles, rolename)
		}
		for _, userprojectid := range user.ProjectIDs {
			if projectid != userprojectid {
				continue
			}
			if _, ok := tmp[userprojectid]; ok {
				continue
			}
			tmp[userprojectid] = &api.ProjectRoles{
				ProjectId: req.ProjectId,
				Roles:     make([]string, 0),
			}
		}
		respuser := &api.UserInfo{
			UserId:       user.ID.Hex(),
			UserName:     user.UserName,
			Department:   user.Department,
			Ctime:        uint32(user.ID.Timestamp().Unix()),
			ProjectRoles: make([]*api.ProjectRoles, 0, len(user.ProjectIDs)),
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
		log.Error(ctx, "[UpdateUser] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[UpdateUser] target's userid format wrong", map[string]interface{}{"target": req.UserId, "error": e})
		return nil, ecode.ErrReq
	}
	if !operator.IsZero() {
		//permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.AdminProjectID+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[UpdateUser] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": model.AdminProjectID, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoUpdateUser(ctx, target, req.NewUserName, req.NewDepartment); e != nil {
		log.Error(ctx, "[UpdateUser] db op failed", map[string]interface{}{"operator": md["Token-Data"], "target": req.UserId, "new_user_name": req.NewUserName, "new_department": req.NewDepartment, "error": e})
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
		log.Error(ctx, "[CreateRole] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[CreateRole] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoCreateRole(ctx, projectid, req.RoleName, req.Comment); e != nil {
		log.Error(ctx, "[CreateRole] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "rolename": req.RoleName, "error": e})
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
		log.Error(ctx, "[SearchRoles] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[SearchRoles] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	roles, page, totalsize, e := s.userDao.MongoSearchRoles(ctx, projectid, req.RoleName, 20, int64(req.Page))
	if e != nil {
		log.Error(ctx, "[SearchRoles] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "rolename": req.RoleName, "error": e})
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
		log.Error(ctx, "[UpdateRole] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[UpdateRole] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoUpdateRole(ctx, projectid, req.RoleName, req.NewComment); e != nil {
		log.Error(ctx, "[UpdateRole] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "rolename": req.RoleName, "error": e})
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
		log.Error(ctx, "[DelRoles] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[DelRoles] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.userDao.MongoDelRoles(ctx, projectid, req.RoleNames); e != nil {
		log.Error(ctx, "[DelRoles] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "rolenames": req.RoleNames, "error": e})
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
		log.Error(ctx, "[AddUserRole] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[AddUserRole] target's userid format wrong", map[string]interface{}{"target": req.UserId, "error": e})
		return nil, ecode.ErrReq
	}
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[AddUserRole] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e = s.userDao.MongoAddUserRole(ctx, target, projectid, req.RoleName); e != nil {
		log.Error(ctx, "[AddUserRole] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "target": req.UserId, "rolename": req.RoleName, "error": e})
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
		log.Error(ctx, "[DelUserRole] operator's token format wrong", map[string]interface{}{"operator": md["Token-Data"], "error": e})
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[DelUserRole] target's userid format wrong", map[string]interface{}{"target": req.UserId, "error": e})
		return nil, ecode.ErrReq
	}
	if !operator.IsZero() {
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			log.Error(ctx, "[DelUserRole] get operator's permission info failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "error": e})
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e = s.userDao.MongoDelUserRole(ctx, target, projectid, req.RoleName); e != nil {
		log.Error(ctx, "[DelUserRole] db op failed", map[string]interface{}{"operator": md["Token-Data"], "project_id": projectid, "target": req.UserId, "rolename": req.RoleName, "error": e})
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.DelUserRoleResp{}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
