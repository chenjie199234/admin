// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.1
// 	protoc              v3.21.1
// source: api/permission.proto

package api

import (
	context "context"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	error1 "github.com/chenjie199234/Corelib/error"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
)

var _CGrpcPathPermissionGetUserPermission = "/admin.permission/get_user_permission"
var _CGrpcPathPermissionUpdateUserPermission = "/admin.permission/update_user_permission"
var _CGrpcPathPermissionAddNode = "/admin.permission/add_node"
var _CGrpcPathPermissionUpdateNode = "/admin.permission/update_node"
var _CGrpcPathPermissionMoveNode = "/admin.permission/move_node"
var _CGrpcPathPermissionDelNode = "/admin.permission/del_node"
var _CGrpcPathPermissionListUserNode = "/admin.permission/list_user_node"
var _CGrpcPathPermissionListNodeUser = "/admin.permission/list_node_user"

type PermissionCGrpcClient interface {
	GetUserPermission(context.Context, *GetUserPermissionReq) (*GetUserPermissionResp, error)
	UpdateUserPermission(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)
	AddNode(context.Context, *AddNodeReq) (*AddNodeResp, error)
	UpdateNode(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)
	MoveNode(context.Context, *MoveNodeReq) (*MoveNodeResp, error)
	DelNode(context.Context, *DelNodeReq) (*DelNodeResp, error)
	ListUserNode(context.Context, *ListUserNodeReq) (*ListUserNodeResp, error)
	ListNodeUser(context.Context, *ListNodeUserReq) (*ListNodeUserResp, error)
}

type permissionCGrpcClient struct {
	cc *cgrpc.CGrpcClient
}

func NewPermissionCGrpcClient(c *cgrpc.CGrpcClient) PermissionCGrpcClient {
	return &permissionCGrpcClient{cc: c}
}

func (c *permissionCGrpcClient) GetUserPermission(ctx context.Context, req *GetUserPermissionReq) (*GetUserPermissionResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(GetUserPermissionResp)
	if e := c.cc.Call(ctx, _CGrpcPathPermissionGetUserPermission, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *permissionCGrpcClient) UpdateUserPermission(ctx context.Context, req *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(UpdateUserPermissionResp)
	if e := c.cc.Call(ctx, _CGrpcPathPermissionUpdateUserPermission, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *permissionCGrpcClient) AddNode(ctx context.Context, req *AddNodeReq) (*AddNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(AddNodeResp)
	if e := c.cc.Call(ctx, _CGrpcPathPermissionAddNode, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *permissionCGrpcClient) UpdateNode(ctx context.Context, req *UpdateNodeReq) (*UpdateNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(UpdateNodeResp)
	if e := c.cc.Call(ctx, _CGrpcPathPermissionUpdateNode, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *permissionCGrpcClient) MoveNode(ctx context.Context, req *MoveNodeReq) (*MoveNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(MoveNodeResp)
	if e := c.cc.Call(ctx, _CGrpcPathPermissionMoveNode, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *permissionCGrpcClient) DelNode(ctx context.Context, req *DelNodeReq) (*DelNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(DelNodeResp)
	if e := c.cc.Call(ctx, _CGrpcPathPermissionDelNode, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *permissionCGrpcClient) ListUserNode(ctx context.Context, req *ListUserNodeReq) (*ListUserNodeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(ListUserNodeResp)
	if e := c.cc.Call(ctx, _CGrpcPathPermissionListUserNode, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *permissionCGrpcClient) ListNodeUser(ctx context.Context, req *ListNodeUserReq) (*ListNodeUserResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(ListNodeUserResp)
	if e := c.cc.Call(ctx, _CGrpcPathPermissionListNodeUser, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}

type PermissionCGrpcServer interface {
	GetUserPermission(context.Context, *GetUserPermissionReq) (*GetUserPermissionResp, error)
	UpdateUserPermission(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)
	AddNode(context.Context, *AddNodeReq) (*AddNodeResp, error)
	UpdateNode(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)
	MoveNode(context.Context, *MoveNodeReq) (*MoveNodeResp, error)
	DelNode(context.Context, *DelNodeReq) (*DelNodeResp, error)
	ListUserNode(context.Context, *ListUserNodeReq) (*ListUserNodeResp, error)
	ListNodeUser(context.Context, *ListNodeUserReq) (*ListNodeUserResp, error)
}

func _Permission_GetUserPermission_CGrpcHandler(handler func(context.Context, *GetUserPermissionReq) (*GetUserPermissionResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetUserPermissionReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/get_user_permission]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetUserPermissionResp)
		}
		ctx.Write(resp)
	}
}
func _Permission_UpdateUserPermission_CGrpcHandler(handler func(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateUserPermissionReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/update_user_permission]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateUserPermissionResp)
		}
		ctx.Write(resp)
	}
}
func _Permission_AddNode_CGrpcHandler(handler func(context.Context, *AddNodeReq) (*AddNodeResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(AddNodeReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/add_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(AddNodeResp)
		}
		ctx.Write(resp)
	}
}
func _Permission_UpdateNode_CGrpcHandler(handler func(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateNodeReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/update_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateNodeResp)
		}
		ctx.Write(resp)
	}
}
func _Permission_MoveNode_CGrpcHandler(handler func(context.Context, *MoveNodeReq) (*MoveNodeResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(MoveNodeReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/move_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(MoveNodeResp)
		}
		ctx.Write(resp)
	}
}
func _Permission_DelNode_CGrpcHandler(handler func(context.Context, *DelNodeReq) (*DelNodeResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelNodeReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/del_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelNodeResp)
		}
		ctx.Write(resp)
	}
}
func _Permission_ListUserNode_CGrpcHandler(handler func(context.Context, *ListUserNodeReq) (*ListUserNodeResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(ListUserNodeReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/list_user_node]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(ListUserNodeResp)
		}
		ctx.Write(resp)
	}
}
func _Permission_ListNodeUser_CGrpcHandler(handler func(context.Context, *ListNodeUserReq) (*ListNodeUserResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(ListNodeUserReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/list_node_user]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(ListNodeUserResp)
		}
		ctx.Write(resp)
	}
}
func RegisterPermissionCGrpcServer(engine *cgrpc.CGrpcServer, svc PermissionCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	//avoid lint
	_ = allmids
	engine.RegisterHandler("admin.permission", "get_user_permission", _Permission_GetUserPermission_CGrpcHandler(svc.GetUserPermission))
	engine.RegisterHandler("admin.permission", "update_user_permission", _Permission_UpdateUserPermission_CGrpcHandler(svc.UpdateUserPermission))
	engine.RegisterHandler("admin.permission", "add_node", _Permission_AddNode_CGrpcHandler(svc.AddNode))
	engine.RegisterHandler("admin.permission", "update_node", _Permission_UpdateNode_CGrpcHandler(svc.UpdateNode))
	engine.RegisterHandler("admin.permission", "move_node", _Permission_MoveNode_CGrpcHandler(svc.MoveNode))
	engine.RegisterHandler("admin.permission", "del_node", _Permission_DelNode_CGrpcHandler(svc.DelNode))
	engine.RegisterHandler("admin.permission", "list_user_node", _Permission_ListUserNode_CGrpcHandler(svc.ListUserNode))
	engine.RegisterHandler("admin.permission", "list_node_user", _Permission_ListNodeUser_CGrpcHandler(svc.ListNodeUser))
}
