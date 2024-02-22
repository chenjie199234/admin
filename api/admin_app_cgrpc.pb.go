// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.98<br />
// 	protoc              v4.25.3<br />
// source: api/admin_app.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	grpc "google.golang.org/grpc"
)

var _CGrpcPathAppGetApp = "/admin.app/get_app"
var _CGrpcPathAppSetApp = "/admin.app/set_app"
var _CGrpcPathAppDelApp = "/admin.app/del_app"
var _CGrpcPathAppUpdateAppSecret = "/admin.app/update_app_secret"
var _CGrpcPathAppDelKey = "/admin.app/del_key"
var _CGrpcPathAppGetKeyConfig = "/admin.app/get_key_config"
var _CGrpcPathAppSetKeyConfig = "/admin.app/set_key_config"
var _CGrpcPathAppRollback = "/admin.app/rollback"
var _CGrpcPathAppWatchConfig = "/admin.app/watch_config"
var _CGrpcPathAppWatchDiscover = "/admin.app/watch_discover"
var _CGrpcPathAppGetInstances = "/admin.app/get_instances"
var _CGrpcPathAppGetInstanceInfo = "/admin.app/get_instance_info"
var _CGrpcPathAppSetProxy = "/admin.app/set_proxy"
var _CGrpcPathAppDelProxy = "/admin.app/del_proxy"
var _CGrpcPathAppProxy = "/admin.app/proxy"

type AppCGrpcClient interface {
	GetApp(context.Context, *GetAppReq, ...grpc.CallOption) (*GetAppResp, error)
	SetApp(context.Context, *SetAppReq, ...grpc.CallOption) (*SetAppResp, error)
	DelApp(context.Context, *DelAppReq, ...grpc.CallOption) (*DelAppResp, error)
	UpdateAppSecret(context.Context, *UpdateAppSecretReq, ...grpc.CallOption) (*UpdateAppSecretResp, error)
	DelKey(context.Context, *DelKeyReq, ...grpc.CallOption) (*DelKeyResp, error)
	GetKeyConfig(context.Context, *GetKeyConfigReq, ...grpc.CallOption) (*GetKeyConfigResp, error)
	SetKeyConfig(context.Context, *SetKeyConfigReq, ...grpc.CallOption) (*SetKeyConfigResp, error)
	Rollback(context.Context, *RollbackReq, ...grpc.CallOption) (*RollbackResp, error)
	WatchConfig(context.Context, *WatchConfigReq, ...grpc.CallOption) (*WatchConfigResp, error)
	WatchDiscover(context.Context, *WatchDiscoverReq, ...grpc.CallOption) (*WatchDiscoverResp, error)
	GetInstances(context.Context, *GetInstancesReq, ...grpc.CallOption) (*GetInstancesResp, error)
	GetInstanceInfo(context.Context, *GetInstanceInfoReq, ...grpc.CallOption) (*GetInstanceInfoResp, error)
	SetProxy(context.Context, *SetProxyReq, ...grpc.CallOption) (*SetProxyResp, error)
	DelProxy(context.Context, *DelProxyReq, ...grpc.CallOption) (*DelProxyResp, error)
	Proxy(context.Context, *ProxyReq, ...grpc.CallOption) (*ProxyResp, error)
}

type appCGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewAppCGrpcClient(cc grpc.ClientConnInterface) AppCGrpcClient {
	return &appCGrpcClient{cc: cc}
}

func (c *appCGrpcClient) GetApp(ctx context.Context, req *GetAppReq, opts ...grpc.CallOption) (*GetAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetAppResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppGetApp, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) SetApp(ctx context.Context, req *SetAppReq, opts ...grpc.CallOption) (*SetAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SetAppResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppSetApp, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) DelApp(ctx context.Context, req *DelAppReq, opts ...grpc.CallOption) (*DelAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelAppResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppDelApp, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) UpdateAppSecret(ctx context.Context, req *UpdateAppSecretReq, opts ...grpc.CallOption) (*UpdateAppSecretResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(UpdateAppSecretResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppUpdateAppSecret, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) DelKey(ctx context.Context, req *DelKeyReq, opts ...grpc.CallOption) (*DelKeyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelKeyResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppDelKey, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) GetKeyConfig(ctx context.Context, req *GetKeyConfigReq, opts ...grpc.CallOption) (*GetKeyConfigResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetKeyConfigResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppGetKeyConfig, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) SetKeyConfig(ctx context.Context, req *SetKeyConfigReq, opts ...grpc.CallOption) (*SetKeyConfigResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SetKeyConfigResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppSetKeyConfig, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) Rollback(ctx context.Context, req *RollbackReq, opts ...grpc.CallOption) (*RollbackResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(RollbackResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppRollback, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) WatchConfig(ctx context.Context, req *WatchConfigReq, opts ...grpc.CallOption) (*WatchConfigResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(WatchConfigResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppWatchConfig, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) WatchDiscover(ctx context.Context, req *WatchDiscoverReq, opts ...grpc.CallOption) (*WatchDiscoverResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(WatchDiscoverResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppWatchDiscover, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) GetInstances(ctx context.Context, req *GetInstancesReq, opts ...grpc.CallOption) (*GetInstancesResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetInstancesResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppGetInstances, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) GetInstanceInfo(ctx context.Context, req *GetInstanceInfoReq, opts ...grpc.CallOption) (*GetInstanceInfoResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(GetInstanceInfoResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppGetInstanceInfo, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) SetProxy(ctx context.Context, req *SetProxyReq, opts ...grpc.CallOption) (*SetProxyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(SetProxyResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppSetProxy, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) DelProxy(ctx context.Context, req *DelProxyReq, opts ...grpc.CallOption) (*DelProxyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(DelProxyResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppDelProxy, req, resp, opts...); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *appCGrpcClient) Proxy(ctx context.Context, req *ProxyReq, opts ...grpc.CallOption) (*ProxyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(ProxyResp)
	if e := c.cc.Invoke(ctx, _CGrpcPathAppProxy, req, resp, opts...); e != nil {
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
	WatchConfig(context.Context, *WatchConfigReq) (*WatchConfigResp, error)
	WatchDiscover(context.Context, *WatchDiscoverReq) (*WatchDiscoverResp, error)
	GetInstances(context.Context, *GetInstancesReq) (*GetInstancesResp, error)
	GetInstanceInfo(context.Context, *GetInstanceInfoReq) (*GetInstanceInfoResp, error)
	SetProxy(context.Context, *SetProxyReq) (*SetProxyResp, error)
	DelProxy(context.Context, *DelProxyReq) (*DelProxyResp, error)
	Proxy(context.Context, *ProxyReq) (*ProxyResp, error)
}

func _App_GetApp_CGrpcHandler(handler func(context.Context, *GetAppReq) (*GetAppResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetAppReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/get_app] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/get_app] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/admin.app/set_app] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/set_app] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/admin.app/del_app] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/del_app] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/admin.app/update_app_secret] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/update_app_secret] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/admin.app/del_key] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/del_key] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/admin.app/get_key_config] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/get_key_config] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/admin.app/set_key_config] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/set_key_config] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/admin.app/rollback] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/rollback] validate failed", log.String("validate", errstr))
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
func _App_WatchConfig_CGrpcHandler(handler func(context.Context, *WatchConfigReq) (*WatchConfigResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(WatchConfigReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/watch_config] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/watch_config] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(WatchConfigResp)
		}
		ctx.Write(resp)
	}
}
func _App_WatchDiscover_CGrpcHandler(handler func(context.Context, *WatchDiscoverReq) (*WatchDiscoverResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(WatchDiscoverReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/watch_discover] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/watch_discover] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(WatchDiscoverResp)
		}
		ctx.Write(resp)
	}
}
func _App_GetInstances_CGrpcHandler(handler func(context.Context, *GetInstancesReq) (*GetInstancesResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetInstancesReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/get_instances] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/get_instances] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetInstancesResp)
		}
		ctx.Write(resp)
	}
}
func _App_GetInstanceInfo_CGrpcHandler(handler func(context.Context, *GetInstanceInfoReq) (*GetInstanceInfoResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetInstanceInfoReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/get_instance_info] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/get_instance_info] validate failed", log.String("validate", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetInstanceInfoResp)
		}
		ctx.Write(resp)
	}
}
func _App_SetProxy_CGrpcHandler(handler func(context.Context, *SetProxyReq) (*SetProxyResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SetProxyReq)
		if e := ctx.DecodeReq(req); e != nil {
			log.Error(ctx, "[/admin.app/set_proxy] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/set_proxy] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/admin.app/del_proxy] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/del_proxy] validate failed", log.String("validate", errstr))
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
			log.Error(ctx, "[/admin.app/proxy] decode failed")
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/proxy] validate failed", log.String("validate", errstr))
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
	engine.RegisterHandler("admin.app", "watch_config", _App_WatchConfig_CGrpcHandler(svc.WatchConfig))
	engine.RegisterHandler("admin.app", "watch_discover", _App_WatchDiscover_CGrpcHandler(svc.WatchDiscover))
	engine.RegisterHandler("admin.app", "get_instances", _App_GetInstances_CGrpcHandler(svc.GetInstances))
	engine.RegisterHandler("admin.app", "get_instance_info", _App_GetInstanceInfo_CGrpcHandler(svc.GetInstanceInfo))
	engine.RegisterHandler("admin.app", "set_proxy", _App_SetProxy_CGrpcHandler(svc.SetProxy))
	engine.RegisterHandler("admin.app", "del_proxy", _App_DelProxy_CGrpcHandler(svc.DelProxy))
	engine.RegisterHandler("admin.app", "proxy", _App_Proxy_CGrpcHandler(svc.Proxy))
}
