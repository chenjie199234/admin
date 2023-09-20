// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.81<br />
// 	protoc              v4.24.1<br />
// source: api/admin_app.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
)

var _CGrpcPathAppGetApp = "/admin.app/get_app"
var _CGrpcPathAppSetApp = "/admin.app/set_app"
var _CGrpcPathAppDelApp = "/admin.app/del_app"
var _CGrpcPathAppUpdateAppSecret = "/admin.app/update_app_secret"
var _CGrpcPathAppDelKey = "/admin.app/del_key"
var _CGrpcPathAppGetKeyConfig = "/admin.app/get_key_config"
var _CGrpcPathAppSetKeyConfig = "/admin.app/set_key_config"
var _CGrpcPathAppRollback = "/admin.app/rollback"
var _CGrpcPathAppWatch = "/admin.app/watch"
var _CGrpcPathAppSetProxy = "/admin.app/set_proxy"
var _CGrpcPathAppDelProxy = "/admin.app/del_proxy"
var _CGrpcPathAppProxy = "/admin.app/proxy"

type AppCGrpcClient interface {
	GetApp(context.Context, *GetAppReq) (*GetAppResp, error)
	SetApp(context.Context, *SetAppReq) (*SetAppResp, error)
	DelApp(context.Context, *DelAppReq) (*DelAppResp, error)
	UpdateAppSecret(context.Context, *UpdateAppSecretReq) (*UpdateAppSecretResp, error)
	DelKey(context.Context, *DelKeyReq) (*DelKeyResp, error)
	GetKeyConfig(context.Context, *GetKeyConfigReq) (*GetKeyConfigResp, error)
	SetKeyConfig(context.Context, *SetKeyConfigReq) (*SetKeyConfigResp, error)
	Rollback(context.Context, *RollbackReq) (*RollbackResp, error)
	Watch(context.Context, *WatchReq) (*WatchResp, error)
	SetProxy(context.Context, *SetProxyReq) (*SetProxyResp, error)
	DelProxy(context.Context, *DelProxyReq) (*DelProxyResp, error)
	Proxy(context.Context, *ProxyReq) (*ProxyResp, error)
}

type appCGrpcClient struct {
	cc *cgrpc.CGrpcClient
}

func NewAppCGrpcClient(c *cgrpc.CGrpcClient) AppCGrpcClient {
	return &appCGrpcClient{cc: c}
}

func (c *appCGrpcClient) GetApp(ctx context.Context, req *GetAppReq) (*GetAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetAppResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppGetApp, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) SetApp(ctx context.Context, req *SetAppReq) (*SetAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SetAppResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppSetApp, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) DelApp(ctx context.Context, req *DelAppReq) (*DelAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelAppResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppDelApp, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) UpdateAppSecret(ctx context.Context, req *UpdateAppSecretReq) (*UpdateAppSecretResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateAppSecretResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppUpdateAppSecret, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) DelKey(ctx context.Context, req *DelKeyReq) (*DelKeyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelKeyResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppDelKey, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) GetKeyConfig(ctx context.Context, req *GetKeyConfigReq) (*GetKeyConfigResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetKeyConfigResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppGetKeyConfig, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) SetKeyConfig(ctx context.Context, req *SetKeyConfigReq) (*SetKeyConfigResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SetKeyConfigResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppSetKeyConfig, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) Rollback(ctx context.Context, req *RollbackReq) (*RollbackResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(RollbackResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppRollback, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) Watch(ctx context.Context, req *WatchReq) (*WatchResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(WatchResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppWatch, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) SetProxy(ctx context.Context, req *SetProxyReq) (*SetProxyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SetProxyResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppSetProxy, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) DelProxy(ctx context.Context, req *DelProxyReq) (*DelProxyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelProxyResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppDelProxy, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) Proxy(ctx context.Context, req *ProxyReq) (*ProxyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(ProxyResp)
	if e := c.cc.Call(ctx, _CGrpcPathAppProxy, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}

type AppCGrpcServer interface {
	GetApp(context.Context, *GetAppReq) (*GetAppResp, error)
	SetApp(context.Context, *SetAppReq) (*SetAppResp, error)
	DelApp(context.Context, *DelAppReq) (*DelAppResp, error)
	UpdateAppSecret(context.Context, *UpdateAppSecretReq) (*UpdateAppSecretResp, error)
	DelKey(context.Context, *DelKeyReq) (*DelKeyResp, error)
	GetKeyConfig(context.Context, *GetKeyConfigReq) (*GetKeyConfigResp, error)
	SetKeyConfig(context.Context, *SetKeyConfigReq) (*SetKeyConfigResp, error)
	Rollback(context.Context, *RollbackReq) (*RollbackResp, error)
	Watch(context.Context, *WatchReq) (*WatchResp, error)
	SetProxy(context.Context, *SetProxyReq) (*SetProxyResp, error)
	DelProxy(context.Context, *DelProxyReq) (*DelProxyResp, error)
	Proxy(context.Context, *ProxyReq) (*ProxyResp, error)
}

func _App_GetApp_CGrpcHandler(handler func(context.Context, *GetAppReq) (*GetAppResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetAppReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/get_app]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/get_app]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetAppResp)
		}
		ctx.Write(resp)
	}
}
func _App_SetApp_CGrpcHandler(handler func(context.Context, *SetAppReq) (*SetAppResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SetAppReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/set_app]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/set_app]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SetAppResp)
		}
		ctx.Write(resp)
	}
}
func _App_DelApp_CGrpcHandler(handler func(context.Context, *DelAppReq) (*DelAppResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelAppReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/del_app]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/del_app]", map[string]interface{}{"error": errstr})
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
func _App_UpdateAppSecret_CGrpcHandler(handler func(context.Context, *UpdateAppSecretReq) (*UpdateAppSecretResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(UpdateAppSecretReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/update_app_secret]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/update_app_secret]", map[string]interface{}{"error": errstr})
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
func _App_DelKey_CGrpcHandler(handler func(context.Context, *DelKeyReq) (*DelKeyResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelKeyReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/del_key]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/del_key]", map[string]interface{}{"error": errstr})
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
func _App_GetKeyConfig_CGrpcHandler(handler func(context.Context, *GetKeyConfigReq) (*GetKeyConfigResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetKeyConfigReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/get_key_config]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/get_key_config]", map[string]interface{}{"error": errstr})
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
func _App_SetKeyConfig_CGrpcHandler(handler func(context.Context, *SetKeyConfigReq) (*SetKeyConfigResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SetKeyConfigReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/set_key_config]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/set_key_config]", map[string]interface{}{"error": errstr})
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
func _App_Rollback_CGrpcHandler(handler func(context.Context, *RollbackReq) (*RollbackResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(RollbackReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/rollback]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/rollback]", map[string]interface{}{"error": errstr})
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
func _App_Watch_CGrpcHandler(handler func(context.Context, *WatchReq) (*WatchResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(WatchReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/watch]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/watch]", map[string]interface{}{"error": errstr})
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
func _App_SetProxy_CGrpcHandler(handler func(context.Context, *SetProxyReq) (*SetProxyResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SetProxyReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/set_proxy]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/set_proxy]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SetProxyResp)
		}
		ctx.Write(resp)
	}
}
func _App_DelProxy_CGrpcHandler(handler func(context.Context, *DelProxyReq) (*DelProxyResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(DelProxyReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/del_proxy]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/del_proxy]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(DelProxyResp)
		}
		ctx.Write(resp)
	}
}
func _App_Proxy_CGrpcHandler(handler func(context.Context, *ProxyReq) (*ProxyResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(ProxyReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/proxy]", map[string]interface{}{"error": e})
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/proxy]", map[string]interface{}{"error": errstr})
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(ProxyResp)
		}
		ctx.Write(resp)
	}
}
func RegisterAppCGrpcServer(engine *cgrpc.CGrpcServer, svc AppCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("admin.app", "get_app", _App_GetApp_CGrpcHandler(svc.GetApp))
	engine.RegisterHandler("admin.app", "set_app", _App_SetApp_CGrpcHandler(svc.SetApp))
	engine.RegisterHandler("admin.app", "del_app", _App_DelApp_CGrpcHandler(svc.DelApp))
	engine.RegisterHandler("admin.app", "update_app_secret", _App_UpdateAppSecret_CGrpcHandler(svc.UpdateAppSecret))
	engine.RegisterHandler("admin.app", "del_key", _App_DelKey_CGrpcHandler(svc.DelKey))
	engine.RegisterHandler("admin.app", "get_key_config", _App_GetKeyConfig_CGrpcHandler(svc.GetKeyConfig))
	engine.RegisterHandler("admin.app", "set_key_config", _App_SetKeyConfig_CGrpcHandler(svc.SetKeyConfig))
	engine.RegisterHandler("admin.app", "rollback", _App_Rollback_CGrpcHandler(svc.Rollback))
	engine.RegisterHandler("admin.app", "watch", _App_Watch_CGrpcHandler(svc.Watch))
	engine.RegisterHandler("admin.app", "set_proxy", _App_SetProxy_CGrpcHandler(svc.SetProxy))
	engine.RegisterHandler("admin.app", "del_proxy", _App_DelProxy_CGrpcHandler(svc.DelProxy))
	engine.RegisterHandler("admin.app", "proxy", _App_Proxy_CGrpcHandler(svc.Proxy))
}
