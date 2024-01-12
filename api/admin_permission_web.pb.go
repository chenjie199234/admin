// Code generated by protoc-gen-go-web. DO NOT EDIT.
// version:
// 	protoc-gen-go-web v0.0.96<br />
// 	protoc            v4.25.1<br />
// source: api/admin_permission.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	web "github.com/chenjie199234/Corelib/web"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	io "io"
	http "net/http"
	strings "strings"
)

var _WebPathPermissionGetUserPermission = "/admin.permission/get_user_permission"
var _WebPathPermissionUpdateUserPermission = "/admin.permission/update_user_permission"
var _WebPathPermissionUpdateRolePermission = "/admin.permission/update_role_permission"
var _WebPathPermissionAddNode = "/admin.permission/add_node"
var _WebPathPermissionUpdateNode = "/admin.permission/update_node"
var _WebPathPermissionMoveNode = "/admin.permission/move_node"
var _WebPathPermissionDelNode = "/admin.permission/del_node"
var _WebPathPermissionListUserNode = "/admin.permission/list_user_node"
var _WebPathPermissionListRoleNode = "/admin.permission/list_role_node"
var _WebPathPermissionListProjectNode = "/admin.permission/list_project_node"

type PermissionWebClient interface {
	GetUserPermission(context.Context, *GetUserPermissionReq, http.Header) (*GetUserPermissionResp, error)
	UpdateUserPermission(context.Context, *UpdateUserPermissionReq, http.Header) (*UpdateUserPermissionResp, error)
	UpdateRolePermission(context.Context, *UpdateRolePermissionReq, http.Header) (*UpdateRolePermissionResp, error)
	AddNode(context.Context, *AddNodeReq, http.Header) (*AddNodeResp, error)
	UpdateNode(context.Context, *UpdateNodeReq, http.Header) (*UpdateNodeResp, error)
	MoveNode(context.Context, *MoveNodeReq, http.Header) (*MoveNodeResp, error)
	DelNode(context.Context, *DelNodeReq, http.Header) (*DelNodeResp, error)
	ListUserNode(context.Context, *ListUserNodeReq, http.Header) (*ListUserNodeResp, error)
	ListRoleNode(context.Context, *ListRoleNodeReq, http.Header) (*ListRoleNodeResp, error)
	ListProjectNode(context.Context, *ListProjectNodeReq, http.Header) (*ListProjectNodeResp, error)
}

type permissionWebClient struct {
	cc *web.WebClient
}

func NewPermissionWebClient(c *web.WebClient) PermissionWebClient {
	return &permissionWebClient{cc: c}
}

func (c *permissionWebClient) GetUserPermission(ctx context.Context, req *GetUserPermissionReq, header http.Header) (*GetUserPermissionResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathPermissionGetUserPermission, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(GetUserPermissionResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionWebClient) UpdateUserPermission(ctx context.Context, req *UpdateUserPermissionReq, header http.Header) (*UpdateUserPermissionResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathPermissionUpdateUserPermission, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(UpdateUserPermissionResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionWebClient) UpdateRolePermission(ctx context.Context, req *UpdateRolePermissionReq, header http.Header) (*UpdateRolePermissionResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathPermissionUpdateRolePermission, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(UpdateRolePermissionResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionWebClient) AddNode(ctx context.Context, req *AddNodeReq, header http.Header) (*AddNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathPermissionAddNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(AddNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionWebClient) UpdateNode(ctx context.Context, req *UpdateNodeReq, header http.Header) (*UpdateNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathPermissionUpdateNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(UpdateNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionWebClient) MoveNode(ctx context.Context, req *MoveNodeReq, header http.Header) (*MoveNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathPermissionMoveNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(MoveNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionWebClient) DelNode(ctx context.Context, req *DelNodeReq, header http.Header) (*DelNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathPermissionDelNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(DelNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionWebClient) ListUserNode(ctx context.Context, req *ListUserNodeReq, header http.Header) (*ListUserNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathPermissionListUserNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(ListUserNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionWebClient) ListRoleNode(ctx context.Context, req *ListRoleNodeReq, header http.Header) (*ListRoleNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathPermissionListRoleNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(ListRoleNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *permissionWebClient) ListProjectNode(ctx context.Context, req *ListProjectNodeReq, header http.Header) (*ListProjectNodeResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathPermissionListProjectNode, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(ListProjectNodeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type PermissionWebServer interface {
	GetUserPermission(context.Context, *GetUserPermissionReq) (*GetUserPermissionResp, error)
	UpdateUserPermission(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)
	UpdateRolePermission(context.Context, *UpdateRolePermissionReq) (*UpdateRolePermissionResp, error)
	AddNode(context.Context, *AddNodeReq) (*AddNodeResp, error)
	UpdateNode(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)
	MoveNode(context.Context, *MoveNodeReq) (*MoveNodeResp, error)
	DelNode(context.Context, *DelNodeReq) (*DelNodeResp, error)
	ListUserNode(context.Context, *ListUserNodeReq) (*ListUserNodeResp, error)
	ListRoleNode(context.Context, *ListRoleNodeReq) (*ListRoleNodeResp, error)
	ListProjectNode(context.Context, *ListProjectNodeReq) (*ListProjectNodeResp, error)
}

func _Permission_GetUserPermission_WebHandler(handler func(context.Context, *GetUserPermissionReq) (*GetUserPermissionResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(GetUserPermissionReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/get_user_permission] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/get_user_permission] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/get_user_permission] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/get_user_permission] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/admin.permission/get_user_permission] Content-Type unknown,must be application/json or application/x-protobuf")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/get_user_permission] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(GetUserPermissionResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Permission_UpdateUserPermission_WebHandler(handler func(context.Context, *UpdateUserPermissionReq) (*UpdateUserPermissionResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateUserPermissionReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/update_user_permission] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/update_user_permission] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/update_user_permission] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/update_user_permission] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/admin.permission/update_user_permission] Content-Type unknown,must be application/json or application/x-protobuf")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/update_user_permission] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(UpdateUserPermissionResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Permission_UpdateRolePermission_WebHandler(handler func(context.Context, *UpdateRolePermissionReq) (*UpdateRolePermissionResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateRolePermissionReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/update_role_permission] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/update_role_permission] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/update_role_permission] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/update_role_permission] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/admin.permission/update_role_permission] Content-Type unknown,must be application/json or application/x-protobuf")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/update_role_permission] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(UpdateRolePermissionResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Permission_AddNode_WebHandler(handler func(context.Context, *AddNodeReq) (*AddNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(AddNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/add_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/add_node] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/add_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/add_node] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/admin.permission/add_node] Content-Type unknown,must be application/json or application/x-protobuf")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/add_node] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(AddNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Permission_UpdateNode_WebHandler(handler func(context.Context, *UpdateNodeReq) (*UpdateNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/update_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/update_node] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/update_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/update_node] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/admin.permission/update_node] Content-Type unknown,must be application/json or application/x-protobuf")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/update_node] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(UpdateNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Permission_MoveNode_WebHandler(handler func(context.Context, *MoveNodeReq) (*MoveNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(MoveNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/move_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/move_node] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/move_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/move_node] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/admin.permission/move_node] Content-Type unknown,must be application/json or application/x-protobuf")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/move_node] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(MoveNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Permission_DelNode_WebHandler(handler func(context.Context, *DelNodeReq) (*DelNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(DelNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/del_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/del_node] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/del_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/del_node] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/admin.permission/del_node] Content-Type unknown,must be application/json or application/x-protobuf")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/del_node] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(DelNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Permission_ListUserNode_WebHandler(handler func(context.Context, *ListUserNodeReq) (*ListUserNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(ListUserNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/list_user_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/list_user_node] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/list_user_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/list_user_node] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/admin.permission/list_user_node] Content-Type unknown,must be application/json or application/x-protobuf")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/list_user_node] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(ListUserNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Permission_ListRoleNode_WebHandler(handler func(context.Context, *ListRoleNodeReq) (*ListRoleNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(ListRoleNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/list_role_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/list_role_node] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/list_role_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/list_role_node] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/admin.permission/list_role_node] Content-Type unknown,must be application/json or application/x-protobuf")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/list_role_node] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(ListRoleNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func _Permission_ListProjectNode_WebHandler(handler func(context.Context, *ListProjectNodeReq) (*ListProjectNodeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(ListProjectNodeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/list_project_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/list_project_node] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.permission/list_project_node] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.permission/list_project_node] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			log.Error(ctx, "[/admin.permission/list_project_node] Content-Type unknown,must be application/json or application/x-protobuf")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.permission/list_project_node] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(ListProjectNodeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func RegisterPermissionWebServer(router *web.Router, svc PermissionWebServer, allmids map[string]web.OutsideHandler) {
	// avoid lint
	_ = allmids
	{
		requiredMids := []string{"accesskey"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Permission_GetUserPermission_WebHandler(svc.GetUserPermission))
		router.Post(_WebPathPermissionGetUserPermission, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Permission_UpdateUserPermission_WebHandler(svc.UpdateUserPermission))
		router.Post(_WebPathPermissionUpdateUserPermission, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Permission_UpdateRolePermission_WebHandler(svc.UpdateRolePermission))
		router.Post(_WebPathPermissionUpdateRolePermission, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Permission_AddNode_WebHandler(svc.AddNode))
		router.Post(_WebPathPermissionAddNode, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Permission_UpdateNode_WebHandler(svc.UpdateNode))
		router.Post(_WebPathPermissionUpdateNode, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Permission_MoveNode_WebHandler(svc.MoveNode))
		router.Post(_WebPathPermissionMoveNode, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Permission_DelNode_WebHandler(svc.DelNode))
		router.Post(_WebPathPermissionDelNode, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Permission_ListUserNode_WebHandler(svc.ListUserNode))
		router.Post(_WebPathPermissionListUserNode, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Permission_ListRoleNode_WebHandler(svc.ListRoleNode))
		router.Post(_WebPathPermissionListRoleNode, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Permission_ListProjectNode_WebHandler(svc.ListProjectNode))
		router.Post(_WebPathPermissionListProjectNode, mids...)
	}
}
