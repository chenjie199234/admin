// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.77<br />
// 	protoc             v4.22.3<br />
// source: api/app.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
)

var _CrpcPathAppGetApp = "/admin.app/get_app"
var _CrpcPathAppGetAppInstances = "/admin.app/get_app_instances"
var _CrpcPathAppGetAppInstanceCmd = "/admin.app/get_app_instance_cmd"
var _CrpcPathAppCreateApp = "/admin.app/create_app"
var _CrpcPathAppDelApp = "/admin.app/del_app"
var _CrpcPathAppUpdateAppSecret = "/admin.app/update_app_secret"
var _CrpcPathAppDelKey = "/admin.app/del_key"
var _CrpcPathAppGetKeyConfig = "/admin.app/get_key_config"
var _CrpcPathAppSetKeyConfig = "/admin.app/set_key_config"
var _CrpcPathAppRollback = "/admin.app/rollback"
var _CrpcPathAppWatch = "/admin.app/watch"
var _CrpcPathAppSetProxy = "/admin.app/set_proxy"
var _CrpcPathAppDelProxy = "/admin.app/del_proxy"
var _CrpcPathAppProxy = "/admin.app/proxy"

type AppCrpcClient interface {
	GetApp(context.Context, *GetAppReq) (*GetAppResp, error)
	GetAppInstances(context.Context, *GetAppInstancesReq) (*GetAppInstancesResp, error)
	GetAppInstanceCmd(context.Context, *GetAppInstanceCmdReq) (*GetAppInstanceCmdResp, error)
	CreateApp(context.Context, *CreateAppReq) (*CreateAppResp, error)
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

type appCrpcClient struct {
	cc *crpc.CrpcClient
}

func NewAppCrpcClient(c *crpc.CrpcClient) AppCrpcClient {
	return &appCrpcClient{cc: c}
}

func (c *appCrpcClient) GetApp(ctx context.Context, req *GetAppReq) (*GetAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppGetApp, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(GetAppResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) GetAppInstances(ctx context.Context, req *GetAppInstancesReq) (*GetAppInstancesResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppGetAppInstances, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(GetAppInstancesResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) GetAppInstanceCmd(ctx context.Context, req *GetAppInstanceCmdReq) (*GetAppInstanceCmdResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppGetAppInstanceCmd, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(GetAppInstanceCmdResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) CreateApp(ctx context.Context, req *CreateAppReq) (*CreateAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppCreateApp, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(CreateAppResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) DelApp(ctx context.Context, req *DelAppReq) (*DelAppResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppDelApp, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(DelAppResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) UpdateAppSecret(ctx context.Context, req *UpdateAppSecretReq) (*UpdateAppSecretResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppUpdateAppSecret, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(UpdateAppSecretResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) DelKey(ctx context.Context, req *DelKeyReq) (*DelKeyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppDelKey, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(DelKeyResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) GetKeyConfig(ctx context.Context, req *GetKeyConfigReq) (*GetKeyConfigResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppGetKeyConfig, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(GetKeyConfigResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) SetKeyConfig(ctx context.Context, req *SetKeyConfigReq) (*SetKeyConfigResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppSetKeyConfig, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(SetKeyConfigResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) Rollback(ctx context.Context, req *RollbackReq) (*RollbackResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppRollback, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(RollbackResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) Watch(ctx context.Context, req *WatchReq) (*WatchResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppWatch, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(WatchResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) SetProxy(ctx context.Context, req *SetProxyReq) (*SetProxyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppSetProxy, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(SetProxyResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) DelProxy(ctx context.Context, req *DelProxyReq) (*DelProxyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppDelProxy, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(DelProxyResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *appCrpcClient) Proxy(ctx context.Context, req *ProxyReq) (*ProxyResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathAppProxy, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(ProxyResp)
	if len(respd) == 0 {
		return resp, nil
	}
	if len(respd) >= 2 && respd[0] == '{' && respd[len(respd)-1] == '}' {
		if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(respd, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type AppCrpcServer interface {
	GetApp(context.Context, *GetAppReq) (*GetAppResp, error)
	GetAppInstances(context.Context, *GetAppInstancesReq) (*GetAppInstancesResp, error)
	GetAppInstanceCmd(context.Context, *GetAppInstanceCmdReq) (*GetAppInstanceCmdResp, error)
	CreateApp(context.Context, *CreateAppReq) (*CreateAppResp, error)
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

func _App_GetApp_CrpcHandler(handler func(context.Context, *GetAppReq) (*GetAppResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(GetAppReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/get_app] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/get_app] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/get_app]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_GetAppInstances_CrpcHandler(handler func(context.Context, *GetAppInstancesReq) (*GetAppInstancesResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(GetAppInstancesReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/get_app_instances] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/get_app_instances] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/get_app_instances]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetAppInstancesResp)
		}
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_GetAppInstanceCmd_CrpcHandler(handler func(context.Context, *GetAppInstanceCmdReq) (*GetAppInstanceCmdResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(GetAppInstanceCmdReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/get_app_instance_cmd] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/get_app_instance_cmd] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/get_app_instance_cmd]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetAppInstanceCmdResp)
		}
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_CreateApp_CrpcHandler(handler func(context.Context, *CreateAppReq) (*CreateAppResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(CreateAppReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/create_app] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/create_app] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/create_app]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_DelApp_CrpcHandler(handler func(context.Context, *DelAppReq) (*DelAppResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(DelAppReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/del_app] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/del_app] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/del_app]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_UpdateAppSecret_CrpcHandler(handler func(context.Context, *UpdateAppSecretReq) (*UpdateAppSecretResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(UpdateAppSecretReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/update_app_secret] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/update_app_secret] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/update_app_secret]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_DelKey_CrpcHandler(handler func(context.Context, *DelKeyReq) (*DelKeyResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(DelKeyReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/del_key] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/del_key] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/del_key]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_GetKeyConfig_CrpcHandler(handler func(context.Context, *GetKeyConfigReq) (*GetKeyConfigResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(GetKeyConfigReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/get_key_config] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/get_key_config] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/get_key_config]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_SetKeyConfig_CrpcHandler(handler func(context.Context, *SetKeyConfigReq) (*SetKeyConfigResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(SetKeyConfigReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/set_key_config] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/set_key_config] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/set_key_config]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_Rollback_CrpcHandler(handler func(context.Context, *RollbackReq) (*RollbackResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(RollbackReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/rollback] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/rollback] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/rollback]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_Watch_CrpcHandler(handler func(context.Context, *WatchReq) (*WatchResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(WatchReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/watch] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/watch] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/watch]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_SetProxy_CrpcHandler(handler func(context.Context, *SetProxyReq) (*SetProxyResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(SetProxyReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/set_proxy] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/set_proxy] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/set_proxy]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_DelProxy_CrpcHandler(handler func(context.Context, *DelProxyReq) (*DelProxyResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(DelProxyReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/del_proxy] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/del_proxy] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/del_proxy]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func _App_Proxy_CrpcHandler(handler func(context.Context, *ProxyReq) (*ProxyResp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		var preferJSON bool
		req := new(ProxyReq)
		reqbody := ctx.GetBody()
		if len(reqbody) >= 2 && reqbody[0] == '{' && reqbody[len(reqbody)-1] == '}' {
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				req.Reset()
				if e := proto.Unmarshal(reqbody, req); e != nil {
					log.Error(ctx, "[/admin.app/proxy] json and proto format decode both failed")
					ctx.Abort(cerror.ErrReq)
					return
				}
			} else {
				preferJSON = true
			}
		} else if e := proto.Unmarshal(reqbody, req); e != nil {
			req.Reset()
			if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(reqbody, req); e != nil {
				log.Error(ctx, "[/admin.app/proxy] json and proto format decode both failed")
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				preferJSON = true
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.app/proxy]", errstr)
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
		if preferJSON {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write(respd)
		} else {
			respd, _ := proto.Marshal(resp)
			ctx.Write(respd)
		}
	}
}
func RegisterAppCrpcServer(engine *crpc.CrpcServer, svc AppCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler(_CrpcPathAppGetApp, _App_GetApp_CrpcHandler(svc.GetApp))
	engine.RegisterHandler(_CrpcPathAppGetAppInstances, _App_GetAppInstances_CrpcHandler(svc.GetAppInstances))
	engine.RegisterHandler(_CrpcPathAppGetAppInstanceCmd, _App_GetAppInstanceCmd_CrpcHandler(svc.GetAppInstanceCmd))
	engine.RegisterHandler(_CrpcPathAppCreateApp, _App_CreateApp_CrpcHandler(svc.CreateApp))
	engine.RegisterHandler(_CrpcPathAppDelApp, _App_DelApp_CrpcHandler(svc.DelApp))
	engine.RegisterHandler(_CrpcPathAppUpdateAppSecret, _App_UpdateAppSecret_CrpcHandler(svc.UpdateAppSecret))
	engine.RegisterHandler(_CrpcPathAppDelKey, _App_DelKey_CrpcHandler(svc.DelKey))
	engine.RegisterHandler(_CrpcPathAppGetKeyConfig, _App_GetKeyConfig_CrpcHandler(svc.GetKeyConfig))
	engine.RegisterHandler(_CrpcPathAppSetKeyConfig, _App_SetKeyConfig_CrpcHandler(svc.SetKeyConfig))
	engine.RegisterHandler(_CrpcPathAppRollback, _App_Rollback_CrpcHandler(svc.Rollback))
	engine.RegisterHandler(_CrpcPathAppWatch, _App_Watch_CrpcHandler(svc.Watch))
	engine.RegisterHandler(_CrpcPathAppSetProxy, _App_SetProxy_CrpcHandler(svc.SetProxy))
	engine.RegisterHandler(_CrpcPathAppDelProxy, _App_DelProxy_CrpcHandler(svc.DelProxy))
	engine.RegisterHandler(_CrpcPathAppProxy, _App_Proxy_CrpcHandler(svc.Proxy))
}
