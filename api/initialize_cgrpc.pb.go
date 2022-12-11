// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.75-dev
// 	protoc              v3.21.11
// source: api/initialize.proto

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
)

var _CGrpcPathInitializeInit = "/admin.initialize/init"
var _CGrpcPathInitializeRootLogin = "/admin.initialize/root_login"
var _CGrpcPathInitializeRootPassword = "/admin.initialize/root_password"
var _CGrpcPathInitializeCreateProject = "/admin.initialize/create_project"
var _CGrpcPathInitializeUpdateProject = "/admin.initialize/update_project"
var _CGrpcPathInitializeListProject = "/admin.initialize/list_project"
var _CGrpcPathInitializeDeleteProject = "/admin.initialize/delete_project"

type InitializeCGrpcClient interface {
	// 初始化
	Init(context.Context, *InitReq) (*InitResp, error)
	// 登录
	RootLogin(context.Context, *RootLoginReq) (*RootLoginResp, error)
	// 更新密码
	RootPassword(context.Context, *RootPasswordReq) (*RootPasswordResp, error)
	// 创建项目
	CreateProject(context.Context, *CreateProjectReq) (*CreateProjectResp, error)
	// 更新项目
	UpdateProject(context.Context, *UpdateProjectReq) (*UpdateProjectResp, error)
	// 获取项目列表
	ListProject(context.Context, *ListProjectReq) (*ListProjectResp, error)
	// 删除项目
	DeleteProject(context.Context, *DeleteProjectReq) (*DeleteProjectResp, error)
}

type initializeCGrpcClient struct {
	cc *cgrpc.CGrpcClient
}

func NewInitializeCGrpcClient(c *cgrpc.CGrpcClient) InitializeCGrpcClient {
	return &initializeCGrpcClient{cc: c}
}

func (c *initializeCGrpcClient) Init(ctx context.Context, req *InitReq) (*InitResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(InitResp)
	if e := c.cc.Call(ctx, _CGrpcPathInitializeInit, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) RootLogin(ctx context.Context, req *RootLoginReq) (*RootLoginResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(RootLoginResp)
	if e := c.cc.Call(ctx, _CGrpcPathInitializeRootLogin, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) RootPassword(ctx context.Context, req *RootPasswordReq) (*RootPasswordResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(RootPasswordResp)
	if e := c.cc.Call(ctx, _CGrpcPathInitializeRootPassword, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) CreateProject(ctx context.Context, req *CreateProjectReq) (*CreateProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(CreateProjectResp)
	if e := c.cc.Call(ctx, _CGrpcPathInitializeCreateProject, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) UpdateProject(ctx context.Context, req *UpdateProjectReq) (*UpdateProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateProjectResp)
	if e := c.cc.Call(ctx, _CGrpcPathInitializeUpdateProject, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) ListProject(ctx context.Context, req *ListProjectReq) (*ListProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(ListProjectResp)
	if e := c.cc.Call(ctx, _CGrpcPathInitializeListProject, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) DeleteProject(ctx context.Context, req *DeleteProjectReq) (*DeleteProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DeleteProjectResp)
	if e := c.cc.Call(ctx, _CGrpcPathInitializeDeleteProject, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}

type InitializeCGrpcServer interface {
	// 初始化
	Init(context.Context, *InitReq) (*InitResp, error)
	// 登录
	RootLogin(context.Context, *RootLoginReq) (*RootLoginResp, error)
	// 更新密码
	RootPassword(context.Context, *RootPasswordReq) (*RootPasswordResp, error)
	// 创建项目
	CreateProject(context.Context, *CreateProjectReq) (*CreateProjectResp, error)
	// 更新项目
	UpdateProject(context.Context, *UpdateProjectReq) (*UpdateProjectResp, error)
	// 获取项目列表
	ListProject(context.Context, *ListProjectReq) (*ListProjectResp, error)
	// 删除项目
	DeleteProject(context.Context, *DeleteProjectReq) (*DeleteProjectResp, error)
}

func _Initialize_Init_CGrpcHandler(handler func(context.Context, *InitReq) (*InitResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(InitReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/init]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(InitResp)
		}
		ctx.Write(resp)
	}
}
func _Initialize_RootLogin_CGrpcHandler(handler func(context.Context, *RootLoginReq) (*RootLoginResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(RootLoginReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/root_login]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(RootLoginResp)
		}
		ctx.Write(resp)
	}
}
func _Initialize_RootPassword_CGrpcHandler(handler func(context.Context, *RootPasswordReq) (*RootPasswordResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(RootPasswordReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/root_password]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(RootPasswordResp)
		}
		ctx.Write(resp)
	}
}
func _Initialize_CreateProject_CGrpcHandler(handler func(context.Context, *CreateProjectReq) (*CreateProjectResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(CreateProjectReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/create_project]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(CreateProjectResp)
		}
		ctx.Write(resp)
	}
}
func _Initialize_UpdateProject_CGrpcHandler(handler func(context.Context, *UpdateProjectReq) (*UpdateProjectResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateProjectReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/update_project]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateProjectResp)
		}
		ctx.Write(resp)
	}
}
func _Initialize_ListProject_CGrpcHandler(handler func(context.Context, *ListProjectReq) (*ListProjectResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(ListProjectReq)
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
			resp = new(ListProjectResp)
		}
		ctx.Write(resp)
	}
}
func _Initialize_DeleteProject_CGrpcHandler(handler func(context.Context, *DeleteProjectReq) (*DeleteProjectResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DeleteProjectReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/delete_project]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DeleteProjectResp)
		}
		ctx.Write(resp)
	}
}
func RegisterInitializeCGrpcServer(engine *cgrpc.CGrpcServer, svc InitializeCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("admin.initialize", "init", _Initialize_Init_CGrpcHandler(svc.Init))
	engine.RegisterHandler("admin.initialize", "root_login", _Initialize_RootLogin_CGrpcHandler(svc.RootLogin))
	engine.RegisterHandler("admin.initialize", "root_password", _Initialize_RootPassword_CGrpcHandler(svc.RootPassword))
	engine.RegisterHandler("admin.initialize", "create_project", _Initialize_CreateProject_CGrpcHandler(svc.CreateProject))
	engine.RegisterHandler("admin.initialize", "update_project", _Initialize_UpdateProject_CGrpcHandler(svc.UpdateProject))
	engine.RegisterHandler("admin.initialize", "list_project", _Initialize_ListProject_CGrpcHandler(svc.ListProject))
	engine.RegisterHandler("admin.initialize", "delete_project", _Initialize_DeleteProject_CGrpcHandler(svc.DeleteProject))
}
