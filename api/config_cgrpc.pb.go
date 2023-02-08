// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.77<br />
// 	protoc              v3.21.11<br />
// source: api/config.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
)

var _CGrpcPathConfigGroups = "/admin.config/groups"
var _CGrpcPathConfigApps = "/admin.config/apps"
var _CGrpcPathConfigCreateApp = "/admin.config/create_app"
var _CGrpcPathConfigDelApp = "/admin.config/del_app"
var _CGrpcPathConfigUpdateAppSecret = "/admin.config/update_app_secret"
var _CGrpcPathConfigKeys = "/admin.config/keys"
var _CGrpcPathConfigDelKey = "/admin.config/del_key"
var _CGrpcPathConfigGetKeyConfig = "/admin.config/get_key_config"
var _CGrpcPathConfigSetKeyConfig = "/admin.config/set_key_config"
var _CGrpcPathConfigRollback = "/admin.config/rollback"
var _CGrpcPathConfigWatch = "/admin.config/watch"

type ConfigCGrpcClient interface {
	// get all groups
	Groups(context.Context, *GroupsReq) (*GroupsResp, error)
	// get all apps in one specific group
	Apps(context.Context, *AppsReq) (*AppsResp, error)
	// create one specific app
	CreateApp(context.Context, *CreateAppReq) (*CreateAppResp, error)
	// del one specific app in one specific group
	DelApp(context.Context, *DelAppReq) (*DelAppResp, error)
	// update one specific app's secret
	UpdateAppSecret(context.Context, *UpdateAppSecretReq) (*UpdateAppSecretResp, error)
	// get all config's keys in one specific app
	Keys(context.Context, *KeysReq) (*KeysResp, error)
	// del one specific key in one specific app
	DelKey(context.Context, *DelKeyReq) (*DelKeyResp, error)
	// get config
	GetKeyConfig(context.Context, *GetKeyConfigReq) (*GetKeyConfigResp, error)
	// set config
	SetKeyConfig(context.Context, *SetKeyConfigReq) (*SetKeyConfigResp, error)
	// rollback config
	Rollback(context.Context, *RollbackReq) (*RollbackResp, error)
	// watch config
	Watch(context.Context, *WatchReq) (*WatchResp, error)
}

type configCGrpcClient struct {
	cc *cgrpc.CGrpcClient
}

func NewConfigCGrpcClient(c *cgrpc.CGrpcClient) ConfigCGrpcClient {
	return &configCGrpcClient{cc: c}
}

func (c *configCGrpcClient) Groups(ctx context.Context, req *GroupsReq) (*GroupsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GroupsResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigGroups, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *configCGrpcClient) Apps(ctx context.Context, req *AppsReq) (*AppsResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(AppsResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigApps, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *configCGrpcClient) CreateApp(ctx context.Context, req *CreateAppReq) (*CreateAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(CreateAppResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigCreateApp, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *configCGrpcClient) DelApp(ctx context.Context, req *DelAppReq) (*DelAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelAppResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigDelApp, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *configCGrpcClient) UpdateAppSecret(ctx context.Context, req *UpdateAppSecretReq) (*UpdateAppSecretResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateAppSecretResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigUpdateAppSecret, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *configCGrpcClient) Keys(ctx context.Context, req *KeysReq) (*KeysResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(KeysResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigKeys, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *configCGrpcClient) DelKey(ctx context.Context, req *DelKeyReq) (*DelKeyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelKeyResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigDelKey, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *configCGrpcClient) GetKeyConfig(ctx context.Context, req *GetKeyConfigReq) (*GetKeyConfigResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetKeyConfigResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigGetKeyConfig, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *configCGrpcClient) SetKeyConfig(ctx context.Context, req *SetKeyConfigReq) (*SetKeyConfigResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SetKeyConfigResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigSetKeyConfig, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *configCGrpcClient) Rollback(ctx context.Context, req *RollbackReq) (*RollbackResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(RollbackResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigRollback, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *configCGrpcClient) Watch(ctx context.Context, req *WatchReq) (*WatchResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(WatchResp)
	if e := c.cc.Call(ctx, _CGrpcPathConfigWatch, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}

type ConfigCGrpcServer interface {
	// get all groups
	Groups(context.Context, *GroupsReq) (*GroupsResp, error)
	// get all apps in one specific group
	Apps(context.Context, *AppsReq) (*AppsResp, error)
	// create one specific app
	CreateApp(context.Context, *CreateAppReq) (*CreateAppResp, error)
	// del one specific app in one specific group
	DelApp(context.Context, *DelAppReq) (*DelAppResp, error)
	// update one specific app's secret
	UpdateAppSecret(context.Context, *UpdateAppSecretReq) (*UpdateAppSecretResp, error)
	// get all config's keys in one specific app
	Keys(context.Context, *KeysReq) (*KeysResp, error)
	// del one specific key in one specific app
	DelKey(context.Context, *DelKeyReq) (*DelKeyResp, error)
	// get config
	GetKeyConfig(context.Context, *GetKeyConfigReq) (*GetKeyConfigResp, error)
	// set config
	SetKeyConfig(context.Context, *SetKeyConfigReq) (*SetKeyConfigResp, error)
	// rollback config
	Rollback(context.Context, *RollbackReq) (*RollbackResp, error)
	// watch config
	Watch(context.Context, *WatchReq) (*WatchResp, error)
}

func _Config_Groups_CGrpcHandler(handler func(context.Context, *GroupsReq) (*GroupsResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GroupsReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/groups]", errstr)
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
		ctx.Write(resp)
	}
}
func _Config_Apps_CGrpcHandler(handler func(context.Context, *AppsReq) (*AppsResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(AppsReq)
		if ctx.DecodeReq(req) != nil {
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
		ctx.Write(resp)
	}
}
func _Config_CreateApp_CGrpcHandler(handler func(context.Context, *CreateAppReq) (*CreateAppResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(CreateAppReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/create_app]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(CreateAppResp)
		}
		ctx.Write(resp)
	}
}
func _Config_DelApp_CGrpcHandler(handler func(context.Context, *DelAppReq) (*DelAppResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelAppReq)
		if ctx.DecodeReq(req) != nil {
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
		ctx.Write(resp)
	}
}
func _Config_UpdateAppSecret_CGrpcHandler(handler func(context.Context, *UpdateAppSecretReq) (*UpdateAppSecretResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateAppSecretReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/update_app_secret]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(UpdateAppSecretResp)
		}
		ctx.Write(resp)
	}
}
func _Config_Keys_CGrpcHandler(handler func(context.Context, *KeysReq) (*KeysResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(KeysReq)
		if ctx.DecodeReq(req) != nil {
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
		ctx.Write(resp)
	}
}
func _Config_DelKey_CGrpcHandler(handler func(context.Context, *DelKeyReq) (*DelKeyResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelKeyReq)
		if ctx.DecodeReq(req) != nil {
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
		ctx.Write(resp)
	}
}
func _Config_GetKeyConfig_CGrpcHandler(handler func(context.Context, *GetKeyConfigReq) (*GetKeyConfigResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetKeyConfigReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/get_key_config]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetKeyConfigResp)
		}
		ctx.Write(resp)
	}
}
func _Config_SetKeyConfig_CGrpcHandler(handler func(context.Context, *SetKeyConfigReq) (*SetKeyConfigResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SetKeyConfigReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.config/set_key_config]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SetKeyConfigResp)
		}
		ctx.Write(resp)
	}
}
func _Config_Rollback_CGrpcHandler(handler func(context.Context, *RollbackReq) (*RollbackResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(RollbackReq)
		if ctx.DecodeReq(req) != nil {
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
		ctx.Write(resp)
	}
}
func _Config_Watch_CGrpcHandler(handler func(context.Context, *WatchReq) (*WatchResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(WatchReq)
		if ctx.DecodeReq(req) != nil {
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
		ctx.Write(resp)
	}
}
func RegisterConfigCGrpcServer(engine *cgrpc.CGrpcServer, svc ConfigCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("admin.config", "groups", _Config_Groups_CGrpcHandler(svc.Groups))
	engine.RegisterHandler("admin.config", "apps", _Config_Apps_CGrpcHandler(svc.Apps))
	engine.RegisterHandler("admin.config", "create_app", _Config_CreateApp_CGrpcHandler(svc.CreateApp))
	engine.RegisterHandler("admin.config", "del_app", _Config_DelApp_CGrpcHandler(svc.DelApp))
	engine.RegisterHandler("admin.config", "update_app_secret", _Config_UpdateAppSecret_CGrpcHandler(svc.UpdateAppSecret))
	engine.RegisterHandler("admin.config", "keys", _Config_Keys_CGrpcHandler(svc.Keys))
	engine.RegisterHandler("admin.config", "del_key", _Config_DelKey_CGrpcHandler(svc.DelKey))
	engine.RegisterHandler("admin.config", "get_key_config", _Config_GetKeyConfig_CGrpcHandler(svc.GetKeyConfig))
	engine.RegisterHandler("admin.config", "set_key_config", _Config_SetKeyConfig_CGrpcHandler(svc.SetKeyConfig))
	engine.RegisterHandler("admin.config", "rollback", _Config_Rollback_CGrpcHandler(svc.Rollback))
	engine.RegisterHandler("admin.config", "watch", _Config_Watch_CGrpcHandler(svc.Watch))
}
