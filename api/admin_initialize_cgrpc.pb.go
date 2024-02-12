// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.97<br />
// 	protoc              v4.25.1<br />
// source: api/admin_initialize.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	grpc "google.golang.org/grpc"
)

var _CGrpcPathInitializeInitStatus = "/admin.initialize/init_status"
var _CGrpcPathInitializeInit = "/admin.initialize/init"
var _CGrpcPathInitializeRootLogin = "/admin.initialize/root_login"
var _CGrpcPathInitializeUpdateRootPassword = "/admin.initialize/update_root_password"
var _CGrpcPathInitializeCreateProject = "/admin.initialize/create_project"
var _CGrpcPathInitializeUpdateProject = "/admin.initialize/update_project"
var _CGrpcPathInitializeListProject = "/admin.initialize/list_project"
var _CGrpcPathInitializeDeleteProject = "/admin.initialize/delete_project"

type InitializeCGrpcClient interface {
	// 初始化状态
	InitStatus(context.Context, *InitStatusReq, ...grpc.CallOption) (*InitStatusResp, error)
	// 初始化
	Init(context.Context, *InitReq, ...grpc.CallOption) (*InitResp, error)
	// 登录
	RootLogin(context.Context, *RootLoginReq, ...grpc.CallOption) (*RootLoginResp, error)
	// 更新密码
	UpdateRootPassword(context.Context, *UpdateRootPasswordReq, ...grpc.CallOption) (*UpdateRootPasswordResp, error)
	// 创建项目
	CreateProject(context.Context, *CreateProjectReq, ...grpc.CallOption) (*CreateProjectResp, error)
	// 更新项目
	UpdateProject(context.Context, *UpdateProjectReq, ...grpc.CallOption) (*UpdateProjectResp, error)
	// 获取项目列表
	ListProject(context.Context, *ListProjectReq, ...grpc.CallOption) (*ListProjectResp, error)
	// 删除项目
	DeleteProject(context.Context, *DeleteProjectReq, ...grpc.CallOption) (*DeleteProjectResp, error)
}

type initializeCGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewInitializeCGrpcClient(cc grpc.ClientConnInterface) InitializeCGrpcClient {
	return &initializeCGrpcClient{cc: cc}
}

func (c *initializeCGrpcClient) InitStatus(ctx context.Context, req *InitStatusReq, opts ...grpc.CallOption) (*InitStatusResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(InitStatusResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathInitializeInitStatus, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) Init(ctx context.Context, req *InitReq, opts ...grpc.CallOption) (*InitResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(InitResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathInitializeInit, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) RootLogin(ctx context.Context, req *RootLoginReq, opts ...grpc.CallOption) (*RootLoginResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(RootLoginResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathInitializeRootLogin, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) UpdateRootPassword(ctx context.Context, req *UpdateRootPasswordReq, opts ...grpc.CallOption) (*UpdateRootPasswordResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateRootPasswordResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathInitializeUpdateRootPassword, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) CreateProject(ctx context.Context, req *CreateProjectReq, opts ...grpc.CallOption) (*CreateProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(CreateProjectResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathInitializeCreateProject, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) UpdateProject(ctx context.Context, req *UpdateProjectReq, opts ...grpc.CallOption) (*UpdateProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateProjectResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathInitializeUpdateProject, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) ListProject(ctx context.Context, req *ListProjectReq, opts ...grpc.CallOption) (*ListProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(ListProjectResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathInitializeListProject, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *initializeCGrpcClient) DeleteProject(ctx context.Context, req *DeleteProjectReq, opts ...grpc.CallOption) (*DeleteProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DeleteProjectResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathInitializeDeleteProject, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}

type InitializeCGrpcServer interface {
	// 初始化状态
	InitStatus(context.Context, *InitStatusReq) (*InitStatusResp, error)
	// 初始化
	Init(context.Context, *InitReq) (*InitResp, error)
	// 登录
	RootLogin(context.Context, *RootLoginReq) (*RootLoginResp, error)
	// 更新密码
	UpdateRootPassword(context.Context, *UpdateRootPasswordReq) (*UpdateRootPasswordResp, error)
	// 创建项目
	CreateProject(context.Context, *CreateProjectReq) (*CreateProjectResp, error)
	// 更新项目
	UpdateProject(context.Context, *UpdateProjectReq) (*UpdateProjectResp, error)
	// 获取项目列表
	ListProject(context.Context, *ListProjectReq) (*ListProjectResp, error)
	// 删除项目
	DeleteProject(context.Context, *DeleteProjectReq) (*DeleteProjectResp, error)
}

func _Initialize_InitStatus_CGrpcHandler(handler func(context.Context, *InitStatusReq) (*InitStatusResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(InitStatusReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.initialize/init_status] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(InitStatusResp)
		}
		ctx.Write(resp)
	}
}
func _Initialize_Init_CGrpcHandler(handler func(context.Context, *InitReq) (*InitResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(InitReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.initialize/init] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/init] validate failed", log.String("validate", errstr))
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
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.initialize/root_login] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/root_login] validate failed", log.String("validate", errstr))
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
func _Initialize_UpdateRootPassword_CGrpcHandler(handler func(context.Context, *UpdateRootPasswordReq) (*UpdateRootPasswordResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateRootPasswordReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.initialize/update_root_password] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/update_root_password] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateRootPasswordResp)
		}
		ctx.Write(resp)
	}
}
func _Initialize_CreateProject_CGrpcHandler(handler func(context.Context, *CreateProjectReq) (*CreateProjectResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(CreateProjectReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.initialize/create_project] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/create_project] validate failed", log.String("validate", errstr))
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
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.initialize/update_project] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/update_project] validate failed", log.String("validate", errstr))
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
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.initialize/list_project] decode failed")
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
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.initialize/delete_project] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/delete_project] validate failed", log.String("validate", errstr))
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
	engine.RegisterHandler("admin.initialize", "init_status", _Initialize_InitStatus_CGrpcHandler(svc.InitStatus))
	engine.RegisterHandler("admin.initialize", "init", _Initialize_Init_CGrpcHandler(svc.Init))
	engine.RegisterHandler("admin.initialize", "root_login", _Initialize_RootLogin_CGrpcHandler(svc.RootLogin))
	engine.RegisterHandler("admin.initialize", "update_root_password", _Initialize_UpdateRootPassword_CGrpcHandler(svc.UpdateRootPassword))
	engine.RegisterHandler("admin.initialize", "create_project", _Initialize_CreateProject_CGrpcHandler(svc.CreateProject))
	engine.RegisterHandler("admin.initialize", "update_project", _Initialize_UpdateProject_CGrpcHandler(svc.UpdateProject))
	engine.RegisterHandler("admin.initialize", "list_project", _Initialize_ListProject_CGrpcHandler(svc.ListProject))
	engine.RegisterHandler("admin.initialize", "delete_project", _Initialize_DeleteProject_CGrpcHandler(svc.DeleteProject))
}
