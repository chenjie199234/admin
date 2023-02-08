// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.77<br />
// 	protoc              v3.21.11<br />
// source: api/proxy.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
)

var _CGrpcPathProxyTob = "/admin.proxy/tob"

type ProxyCGrpcClient interface {
	Tob(context.Context, *TobReq) (*TobResp, error)
}

type proxyCGrpcClient struct {
	cc *cgrpc.CGrpcClient
}

func NewProxyCGrpcClient(c *cgrpc.CGrpcClient) ProxyCGrpcClient {
	return &proxyCGrpcClient{cc: c}
}

func (c *proxyCGrpcClient) Tob(ctx context.Context, req *TobReq) (*TobResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(TobResp)
	if e := c.cc.Call(ctx, _CGrpcPathProxyTob, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}

type ProxyCGrpcServer interface {
	Tob(context.Context, *TobReq) (*TobResp, error)
}

func _Proxy_Tob_CGrpcHandler(handler func(context.Context, *TobReq) (*TobResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(TobReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.proxy/tob]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(TobResp)
		}
		ctx.Write(resp)
	}
}
func RegisterProxyCGrpcServer(engine *cgrpc.CGrpcServer, svc ProxyCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("admin.proxy", "tob", _Proxy_Tob_CGrpcHandler(svc.Tob))
}
