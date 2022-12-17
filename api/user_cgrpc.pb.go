// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.75
// 	protoc              v3.21.11
// source: api/user.proto

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
)

var _CGrpcPathUserUserLogin = "/admin.user/user_login"
var _CGrpcPathUserInviteProject = "/admin.user/invite_project"
var _CGrpcPathUserKickProject = "/admin.user/kick_project"
var _CGrpcPathUserSearchUsers = "/admin.user/search_users"
var _CGrpcPathUserUpdateUser = "/admin.user/update_user"
var _CGrpcPathUserCreateRole = "/admin.user/create_role"
var _CGrpcPathUserSearchRoles = "/admin.user/search_roles"
var _CGrpcPathUserUpdateRole = "/admin.user/update_role"
var _CGrpcPathUserDelRoles = "/admin.user/del_roles"
var _CGrpcPathUserAddUserRole = "/admin.user/add_user_role"
var _CGrpcPathUserDelUserRole = "/admin.user/del_user_role"

type UserCGrpcClient interface {
	UserLogin(context.Context, *UserLoginReq) (*UserLoginResp, error)
	InviteProject(context.Context, *InviteProjectReq) (*InviteProjectResp, error)
	KickProject(context.Context, *KickProjectReq) (*KickProjectResp, error)
	SearchUsers(context.Context, *SearchUsersReq) (*SearchUsersResp, error)
	UpdateUser(context.Context, *UpdateUserReq) (*UpdateUserResp, error)
	CreateRole(context.Context, *CreateRoleReq) (*CreateRoleResp, error)
	SearchRoles(context.Context, *SearchRolesReq) (*SearchRolesResp, error)
	UpdateRole(context.Context, *UpdateRoleReq) (*UpdateRoleResp, error)
	DelRoles(context.Context, *DelRolesReq) (*DelRolesResp, error)
	AddUserRole(context.Context, *AddUserRoleReq) (*AddUserRoleResp, error)
	DelUserRole(context.Context, *DelUserRoleReq) (*DelUserRoleResp, error)
}

type userCGrpcClient struct {
	cc *cgrpc.CGrpcClient
}

func NewUserCGrpcClient(c *cgrpc.CGrpcClient) UserCGrpcClient {
	return &userCGrpcClient{cc: c}
}

func (c *userCGrpcClient) UserLogin(ctx context.Context, req *UserLoginReq) (*UserLoginResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UserLoginResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserUserLogin, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) InviteProject(ctx context.Context, req *InviteProjectReq) (*InviteProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(InviteProjectResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserInviteProject, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) KickProject(ctx context.Context, req *KickProjectReq) (*KickProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(KickProjectResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserKickProject, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) SearchUsers(ctx context.Context, req *SearchUsersReq) (*SearchUsersResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SearchUsersResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserSearchUsers, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateUser(ctx context.Context, req *UpdateUserReq) (*UpdateUserResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateUserResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserUpdateUser, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) CreateRole(ctx context.Context, req *CreateRoleReq) (*CreateRoleResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(CreateRoleResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserCreateRole, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) SearchRoles(ctx context.Context, req *SearchRolesReq) (*SearchRolesResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SearchRolesResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserSearchRoles, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) UpdateRole(ctx context.Context, req *UpdateRoleReq) (*UpdateRoleResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateRoleResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserUpdateRole, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) DelRoles(ctx context.Context, req *DelRolesReq) (*DelRolesResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelRolesResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserDelRoles, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) AddUserRole(ctx context.Context, req *AddUserRoleReq) (*AddUserRoleResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(AddUserRoleResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserAddUserRole, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) DelUserRole(ctx context.Context, req *DelUserRoleReq) (*DelUserRoleResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelUserRoleResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserDelUserRole, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}

type UserCGrpcServer interface {
	UserLogin(context.Context, *UserLoginReq) (*UserLoginResp, error)
	InviteProject(context.Context, *InviteProjectReq) (*InviteProjectResp, error)
	KickProject(context.Context, *KickProjectReq) (*KickProjectResp, error)
	SearchUsers(context.Context, *SearchUsersReq) (*SearchUsersResp, error)
	UpdateUser(context.Context, *UpdateUserReq) (*UpdateUserResp, error)
	CreateRole(context.Context, *CreateRoleReq) (*CreateRoleResp, error)
	SearchRoles(context.Context, *SearchRolesReq) (*SearchRolesResp, error)
	UpdateRole(context.Context, *UpdateRoleReq) (*UpdateRoleResp, error)
	DelRoles(context.Context, *DelRolesReq) (*DelRolesResp, error)
	AddUserRole(context.Context, *AddUserRoleReq) (*AddUserRoleResp, error)
	DelUserRole(context.Context, *DelUserRoleReq) (*DelUserRoleResp, error)
}

func _User_UserLogin_CGrpcHandler(handler func(context.Context, *UserLoginReq) (*UserLoginResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UserLoginReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UserLoginResp)
		}
		ctx.Write(resp)
	}
}
func _User_InviteProject_CGrpcHandler(handler func(context.Context, *InviteProjectReq) (*InviteProjectResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(InviteProjectReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/invite_project]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(InviteProjectResp)
		}
		ctx.Write(resp)
	}
}
func _User_KickProject_CGrpcHandler(handler func(context.Context, *KickProjectReq) (*KickProjectResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(KickProjectReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/kick_project]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(KickProjectResp)
		}
		ctx.Write(resp)
	}
}
func _User_SearchUsers_CGrpcHandler(handler func(context.Context, *SearchUsersReq) (*SearchUsersResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SearchUsersReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/search_users]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SearchUsersResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateUser_CGrpcHandler(handler func(context.Context, *UpdateUserReq) (*UpdateUserResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateUserReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/update_user]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateUserResp)
		}
		ctx.Write(resp)
	}
}
func _User_CreateRole_CGrpcHandler(handler func(context.Context, *CreateRoleReq) (*CreateRoleResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(CreateRoleReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/create_role]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(CreateRoleResp)
		}
		ctx.Write(resp)
	}
}
func _User_SearchRoles_CGrpcHandler(handler func(context.Context, *SearchRolesReq) (*SearchRolesResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SearchRolesReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/search_roles]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SearchRolesResp)
		}
		ctx.Write(resp)
	}
}
func _User_UpdateRole_CGrpcHandler(handler func(context.Context, *UpdateRoleReq) (*UpdateRoleResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateRoleReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/update_role]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateRoleResp)
		}
		ctx.Write(resp)
	}
}
func _User_DelRoles_CGrpcHandler(handler func(context.Context, *DelRolesReq) (*DelRolesResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelRolesReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/del_roles]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelRolesResp)
		}
		ctx.Write(resp)
	}
}
func _User_AddUserRole_CGrpcHandler(handler func(context.Context, *AddUserRoleReq) (*AddUserRoleResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(AddUserRoleReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/add_user_role]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(AddUserRoleResp)
		}
		ctx.Write(resp)
	}
}
func _User_DelUserRole_CGrpcHandler(handler func(context.Context, *DelUserRoleReq) (*DelUserRoleResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelUserRoleReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/del_user_role]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelUserRoleResp)
		}
		ctx.Write(resp)
	}
}
func RegisterUserCGrpcServer(engine *cgrpc.CGrpcServer, svc UserCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("admin.user", "user_login", _User_UserLogin_CGrpcHandler(svc.UserLogin))
	engine.RegisterHandler("admin.user", "invite_project", _User_InviteProject_CGrpcHandler(svc.InviteProject))
	engine.RegisterHandler("admin.user", "kick_project", _User_KickProject_CGrpcHandler(svc.KickProject))
	engine.RegisterHandler("admin.user", "search_users", _User_SearchUsers_CGrpcHandler(svc.SearchUsers))
	engine.RegisterHandler("admin.user", "update_user", _User_UpdateUser_CGrpcHandler(svc.UpdateUser))
	engine.RegisterHandler("admin.user", "create_role", _User_CreateRole_CGrpcHandler(svc.CreateRole))
	engine.RegisterHandler("admin.user", "search_roles", _User_SearchRoles_CGrpcHandler(svc.SearchRoles))
	engine.RegisterHandler("admin.user", "update_role", _User_UpdateRole_CGrpcHandler(svc.UpdateRole))
	engine.RegisterHandler("admin.user", "del_roles", _User_DelRoles_CGrpcHandler(svc.DelRoles))
	engine.RegisterHandler("admin.user", "add_user_role", _User_AddUserRole_CGrpcHandler(svc.AddUserRole))
	engine.RegisterHandler("admin.user", "del_user_role", _User_DelUserRole_CGrpcHandler(svc.DelUserRole))
}
