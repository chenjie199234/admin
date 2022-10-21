// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.1
// 	protoc             v3.21.1
// source: api/config.proto

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	proto "google.golang.org/protobuf/proto"
)

var _CrpcPathConfigGroups = "/admin.config/groups"
var _CrpcPathConfigDelGroup = "/admin.config/del_group"
var _CrpcPathConfigApps = "/admin.config/apps"
var _CrpcPathConfigDelApp = "/admin.config/del_app"
var _CrpcPathConfigKeys = "/admin.config/keys"
var _CrpcPathConfigDelKey = "/admin.config/del_key"
var _CrpcPathConfigCreate = "/admin.config/create"
var _CrpcPathConfigUpdatecipher = "/admin.config/updatecipher"
var _CrpcPathConfigGet = "/admin.config/get"
var _CrpcPathConfigSet = "/admin.config/set"
var _CrpcPathConfigRollback = "/admin.config/rollback"
var _CrpcPathConfigWatch = "/admin.config/watch"

type ConfigCrpcClient interface {
	// get all groups
	Groups(context.Context, *GroupsReq) (*GroupsResp, error)
	// del one specific group
	DelGroup(context.Context, *DelGroupReq) (*DelGroupResp, error)
	// get all apps in one specific group
	Apps(context.Context, *AppsReq) (*AppsResp, error)
	// del one specific app in one specific group
	DelApp(context.Context, *DelAppReq) (*DelAppResp, error)
	// get all config's keys in one specific app
	Keys(context.Context, *KeysReq) (*KeysResp, error)
	// del one specific key in one specific app
	DelKey(context.Context, *DelKeyReq) (*DelKeyResp, error)
	// create one specific app
	Create(context.Context, *CreateReq) (*CreateResp, error)
	// update one specific app's cipher
	Updatecipher(context.Context, *UpdatecipherReq) (*UpdatecipherResp, error)
	// get config
	Get(context.Context, *GetReq) (*GetResp, error)
	// set config
	Set(context.Context, *SetReq) (*SetResp, error)
	// rollback config
	Rollback(context.Context, *RollbackReq) (*RollbackResp, error)
	// watch config
	Watch(context.Context, *WatchReq) (*WatchResp, error)
}

type configCrpcClient struct {
	cc *crpc.CrpcClient
}

func NewConfigCrpcClient(c *crpc.CrpcClient) ConfigCrpcClient {
	return &configCrpcClient{cc: c}
}

func (c *configCrpcClient) Groups(ctx context.Context, req *GroupsReq) (*GroupsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigGroups, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(GroupsResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) DelGroup(ctx context.Context, req *DelGroupReq) (*DelGroupResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigDelGroup, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(DelGroupResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Apps(ctx context.Context, req *AppsReq) (*AppsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigApps, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(AppsResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) DelApp(ctx context.Context, req *DelAppReq) (*DelAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigDelApp, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(DelAppResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Keys(ctx context.Context, req *KeysReq) (*KeysResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigKeys, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(KeysResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) DelKey(ctx context.Context, req *DelKeyReq) (*DelKeyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigDelKey, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(DelKeyResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Create(ctx context.Context, req *CreateReq) (*CreateResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigCreate, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(CreateResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Updatecipher(ctx context.Context, req *UpdatecipherReq) (*UpdatecipherResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigUpdatecipher, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdatecipherResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Get(ctx context.Context, req *GetReq) (*GetResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigGet, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(GetResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Set(ctx context.Context, req *SetReq) (*SetResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigSet, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(SetResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Rollback(ctx context.Context, req *RollbackReq) (*RollbackResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigRollback, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(RollbackResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *configCrpcClient) Watch(ctx context.Context, req *WatchReq) (*WatchResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathConfigWatch, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(WatchResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type ConfigCrpcServer interface {
	// get all groups
	Groups(context.Context, *GroupsReq) (*GroupsResp, error)
	// del one specific group
	DelGroup(context.Context, *DelGroupReq) (*DelGroupResp, error)
	// get all apps in one specific group
	Apps(context.Context, *AppsReq) (*AppsResp, error)
	// del one specific app in one specific group
	DelApp(context.Context, *DelAppReq) (*DelAppResp, error)
	// get all config's keys in one specific app
	Keys(context.Context, *KeysReq) (*KeysResp, error)
	// del one specific key in one specific app
	DelKey(context.Context, *DelKeyReq) (*DelKeyResp, error)
	// create one specific app
	Create(context.Context, *CreateReq) (*CreateResp, error)
	// update one specific app's cipher
	Updatecipher(context.Context, *UpdatecipherReq) (*UpdatecipherResp, error)
	// get config
	Get(context.Context, *GetReq) (*GetResp, error)
	// set config
	Set(context.Context, *SetReq) (*SetResp, error)
	// rollback config
	Rollback(context.Context, *RollbackReq) (*RollbackResp, error)
	// watch config
	Watch(context.Context, *WatchReq) (*WatchResp, error)
}

func _Config_Groups_CrpcHandler(handler func(context.Context, *GroupsReq) (*GroupsResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(GroupsReq)
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
			resp = new(GroupsResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_DelGroup_CrpcHandler(handler func(context.Context, *DelGroupReq) (*DelGroupResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(DelGroupReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/del_group]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelGroupResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Apps_CrpcHandler(handler func(context.Context, *AppsReq) (*AppsResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(AppsReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/apps]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(AppsResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_DelApp_CrpcHandler(handler func(context.Context, *DelAppReq) (*DelAppResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(DelAppReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/del_app]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelAppResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Keys_CrpcHandler(handler func(context.Context, *KeysReq) (*KeysResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(KeysReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/keys]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(KeysResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_DelKey_CrpcHandler(handler func(context.Context, *DelKeyReq) (*DelKeyResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(DelKeyReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/del_key]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelKeyResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Create_CrpcHandler(handler func(context.Context, *CreateReq) (*CreateResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(CreateReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/create]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(CreateResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Updatecipher_CrpcHandler(handler func(context.Context, *UpdatecipherReq) (*UpdatecipherResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(UpdatecipherReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/updatecipher]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdatecipherResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Get_CrpcHandler(handler func(context.Context, *GetReq) (*GetResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(GetReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/get]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Set_CrpcHandler(handler func(context.Context, *SetReq) (*SetResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(SetReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/set]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SetResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Rollback_CrpcHandler(handler func(context.Context, *RollbackReq) (*RollbackResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(RollbackReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/rollback]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(RollbackResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func _Config_Watch_CrpcHandler(handler func(context.Context, *WatchReq) (*WatchResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(WatchReq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/watch]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(WatchResp)
		}
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func RegisterConfigCrpcServer(engine *crpc.CrpcServer, svc ConfigCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler(_CrpcPathConfigGroups, _Config_Groups_CrpcHandler(svc.Groups))
	engine.RegisterHandler(_CrpcPathConfigDelGroup, _Config_DelGroup_CrpcHandler(svc.DelGroup))
	engine.RegisterHandler(_CrpcPathConfigApps, _Config_Apps_CrpcHandler(svc.Apps))
	engine.RegisterHandler(_CrpcPathConfigDelApp, _Config_DelApp_CrpcHandler(svc.DelApp))
	engine.RegisterHandler(_CrpcPathConfigKeys, _Config_Keys_CrpcHandler(svc.Keys))
	engine.RegisterHandler(_CrpcPathConfigDelKey, _Config_DelKey_CrpcHandler(svc.DelKey))
	engine.RegisterHandler(_CrpcPathConfigCreate, _Config_Create_CrpcHandler(svc.Create))
	engine.RegisterHandler(_CrpcPathConfigUpdatecipher, _Config_Updatecipher_CrpcHandler(svc.Updatecipher))
	engine.RegisterHandler(_CrpcPathConfigGet, _Config_Get_CrpcHandler(svc.Get))
	engine.RegisterHandler(_CrpcPathConfigSet, _Config_Set_CrpcHandler(svc.Set))
	engine.RegisterHandler(_CrpcPathConfigRollback, _Config_Rollback_CrpcHandler(svc.Rollback))
	engine.RegisterHandler(_CrpcPathConfigWatch, _Config_Watch_CrpcHandler(svc.Watch))
}
