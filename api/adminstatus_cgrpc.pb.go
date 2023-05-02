// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.77<br />
// 	protoc              v4.22.3<br />
// source: api/adminstatus.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
)

var _CGrpcPathStatusPing = "/admin.status/ping"

type StatusCGrpcClient interface {
	// ping check server's health
	Ping(context.Context, *Pingreq) (*Pingresp, error)
}

type statusCGrpcClient struct {
	cc *cgrpc.CGrpcClient
}

func NewStatusCGrpcClient(c *cgrpc.CGrpcClient) StatusCGrpcClient {
	return &statusCGrpcClient{cc: c}
}

func (c *statusCGrpcClient) Ping(ctx context.Context, req *Pingreq) (*Pingresp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	resp := new(Pingresp)
	if e := c.cc.Call(ctx, _CGrpcPathStatusPing, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}

type StatusCGrpcServer interface {
	// ping check server's health
	Ping(context.Context, *Pingreq) (*Pingresp, error)
}

func _Status_Ping_CGrpcHandler(handler func(context.Context, *Pingreq) (*Pingresp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(Pingreq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.status/ping]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(Pingresp)
		}
		ctx.Write(resp)
	}
}
func RegisterStatusCGrpcServer(engine *cgrpc.CGrpcServer, svc StatusCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	{
		requiredMids := []string{"accesskey", "rate"}
		mids := make([]cgrpc.OutsideHandler, 0, 3)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Status_Ping_CGrpcHandler(svc.Ping))
		engine.RegisterHandler("admin.status", "ping", mids...)
	}
}
