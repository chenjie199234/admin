// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.1
// 	protoc             v3.21.1
// source: api/permission.proto

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	proto "google.golang.org/protobuf/proto"
)

var _CrpcPathPermissionGetUserPermission = "/admin.permission/get_user_permission"
var _CrpcPathPermissionUpdateUserPermission = "/admin.permission/update_user_permission"
var _CrpcPathPermissionUpdateRolePermission = "/admin.permission/update_role_permission"
var _CrpcPathPermissionAddNode = "/admin.permission/add_node"
var _CrpcPathPermissionUpdateNode = "/admin.permission/update_node"
var _CrpcPathPermissionMoveNode = "/admin.permission/move_node"
var _CrpcPathPermissionDelNode = "/admin.permission/del_node"
var _CrpcPathPermissionListUserNode = "/admin.permission/list_user_node"
var _CrpcPathPermissionListRoleNode = "/admin.permission/list_role_node"
var _CrpcPathPermissionListAllNode = "/admin.permission/list_all_node"

type PermissionCrpcClient interface {
	GetUserPermission(context.Context, *GetUserPermissionReq) (*GetUserPermissionResp, error)
	UpdateUserPermission(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)
	UpdateRolePermission(context.Context, *UpdateRolePermissionReq) (*UpdateRolePermissionResp, error)
	AddNode(context.Context, *AddNodeReq) (*AddNodeResp, error)
	UpdateNode(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)
	MoveNode(context.Context, *MoveNodeReq) (*MoveNodeResp, error)
	DelNode(context.Context, *DelNodeReq) (*DelNodeResp, error)
	ListUserNode(context.Context, *ListUserNodeReq) (*ListUserNodeResp, error)
	ListRoleNode(context.Context, *ListRoleNodeReq) (*ListRoleNodeResp, error)
	ListAllNode(context.Context, *ListAllNodeReq) (*ListAllNodeResp, error)
}

type permissionCrpcClient struct {
	cc *crpc.CrpcClient
}

func NewPermissionCrpcClient(c *crpc.CrpcClient) PermissionCrpcClient {
	return &permissionCrpcClient{cc: c}
}

func (c *permissionCrpcClient) GetUserPermission(ctx context.Context, req *GetUserPermissionReq) (*GetUserPermissionResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathPermissionGetUserPermission, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(GetUserPermissionResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionCrpcClient) UpdateUserPermission(ctx context.Context, req *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathPermissionUpdateUserPermission, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdateUserPermissionResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionCrpcClient) UpdateRolePermission(ctx context.Context, req *UpdateRolePermissionReq) (*UpdateRolePermissionResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathPermissionUpdateRolePermission, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdateRolePermissionResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionCrpcClient) AddNode(ctx context.Context, req *AddNodeReq) (*AddNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathPermissionAddNode, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(AddNodeResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionCrpcClient) UpdateNode(ctx context.Context, req *UpdateNodeReq) (*UpdateNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathPermissionUpdateNode, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdateNodeResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionCrpcClient) MoveNode(ctx context.Context, req *MoveNodeReq) (*MoveNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathPermissionMoveNode, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(MoveNodeResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionCrpcClient) DelNode(ctx context.Context, req *DelNodeReq) (*DelNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathPermissionDelNode, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(DelNodeResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionCrpcClient) ListUserNode(ctx context.Context, req *ListUserNodeReq) (*ListUserNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathPermissionListUserNode, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(ListUserNodeResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionCrpcClient) ListRoleNode(ctx context.Context, req *ListRoleNodeReq) (*ListRoleNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathPermissionListRoleNode, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(ListRoleNodeResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionCrpcClient) ListAllNode(ctx context.Context, req *ListAllNodeReq) (*ListAllNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathPermissionListAllNode, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(ListAllNodeResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type PermissionCrpcServer interface {
	GetUserPermission(context.Context, *GetUserPermissionReq) (*GetUserPermissionResp, error)
	UpdateUserPermission(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)
	UpdateRolePermission(context.Context, *UpdateRolePermissionReq) (*UpdateRolePermissionResp, error)
	AddNode(context.Context, *AddNodeReq) (*AddNodeResp, error)
	UpdateNode(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)
	MoveNode(context.Context, *MoveNodeReq) (*MoveNodeResp, error)
	DelNode(context.Context, *DelNodeReq) (*DelNodeResp, error)
	ListUserNode(context.Context, *ListUserNodeReq) (*ListUserNodeResp, error)
	ListRoleNode(context.Context, *ListRoleNodeReq) (*ListRoleNodeResp, error)
	ListAllNode(context.Context, *ListAllNodeReq) (*ListAllNodeResp, error)
}

func _Permission_GetUserPermission_CrpcHandler(handler func(context.Context, *GetUserPermissionReq) (*GetUserPermissionResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(GetUserPermissionReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/get_user_permission]", errstr)
			ctx.Abort(cerror.ErrReq)
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
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Permission_UpdateUserPermission_CrpcHandler(handler func(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(UpdateUserPermissionReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/update_user_permission]", errstr)
			ctx.Abort(cerror.ErrReq)
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
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Permission_UpdateRolePermission_CrpcHandler(handler func(context.Context, *UpdateRolePermissionReq) (*UpdateRolePermissionResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(UpdateRolePermissionReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/update_role_permission]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateRolePermissionResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Permission_AddNode_CrpcHandler(handler func(context.Context, *AddNodeReq) (*AddNodeResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(AddNodeReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/add_node]", errstr)
			ctx.Abort(cerror.ErrReq)
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
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Permission_UpdateNode_CrpcHandler(handler func(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(UpdateNodeReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/update_node]", errstr)
			ctx.Abort(cerror.ErrReq)
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
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Permission_MoveNode_CrpcHandler(handler func(context.Context, *MoveNodeReq) (*MoveNodeResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(MoveNodeReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/move_node]", errstr)
			ctx.Abort(cerror.ErrReq)
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
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Permission_DelNode_CrpcHandler(handler func(context.Context, *DelNodeReq) (*DelNodeResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(DelNodeReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/del_node]", errstr)
			ctx.Abort(cerror.ErrReq)
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
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Permission_ListUserNode_CrpcHandler(handler func(context.Context, *ListUserNodeReq) (*ListUserNodeResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(ListUserNodeReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
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
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Permission_ListRoleNode_CrpcHandler(handler func(context.Context, *ListRoleNodeReq) (*ListRoleNodeResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(ListRoleNodeReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/list_role_node]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(ListRoleNodeResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Permission_ListAllNode_CrpcHandler(handler func(context.Context, *ListAllNodeReq) (*ListAllNodeResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(ListAllNodeReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(ListAllNodeResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func RegisterPermissionCrpcServer(engine *crpc.CrpcServer, svc PermissionCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler(_CrpcPathPermissionGetUserPermission, _Permission_GetUserPermission_CrpcHandler(svc.GetUserPermission))
	engine.RegisterHandler(_CrpcPathPermissionUpdateUserPermission, _Permission_UpdateUserPermission_CrpcHandler(svc.UpdateUserPermission))
	engine.RegisterHandler(_CrpcPathPermissionUpdateRolePermission, _Permission_UpdateRolePermission_CrpcHandler(svc.UpdateRolePermission))
	engine.RegisterHandler(_CrpcPathPermissionAddNode, _Permission_AddNode_CrpcHandler(svc.AddNode))
	engine.RegisterHandler(_CrpcPathPermissionUpdateNode, _Permission_UpdateNode_CrpcHandler(svc.UpdateNode))
	engine.RegisterHandler(_CrpcPathPermissionMoveNode, _Permission_MoveNode_CrpcHandler(svc.MoveNode))
	engine.RegisterHandler(_CrpcPathPermissionDelNode, _Permission_DelNode_CrpcHandler(svc.DelNode))
	engine.RegisterHandler(_CrpcPathPermissionListUserNode, _Permission_ListUserNode_CrpcHandler(svc.ListUserNode))
	engine.RegisterHandler(_CrpcPathPermissionListRoleNode, _Permission_ListRoleNode_CrpcHandler(svc.ListRoleNode))
	engine.RegisterHandler(_CrpcPathPermissionListAllNode, _Permission_ListAllNode_CrpcHandler(svc.ListAllNode))
}
