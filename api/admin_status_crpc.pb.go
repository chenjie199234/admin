// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.132<br />
// 	protoc             v5.29.0<br />
// source: api/admin_status.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	proto "google.golang.org/protobuf/proto"
	slog "log/slog"
)

var _CrpcPathStatusPing = "/admin.status/ping"

type StatusCrpcClient interface {
	// ping check server's health
	Ping(ctx context.Context, req *Pingreq) (resp *Pingresp, e error)
}

type statusCrpcClient struct {
	cc *crpc.CrpcClient
}

func NewStatusCrpcClient(c *crpc.CrpcClient) StatusCrpcClient {
	return &statusCrpcClient{cc: c}
}

func (c *statusCrpcClient) Ping(ctx context.Context, req *Pingreq) (*Pingresp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if errstr := req.Validate(); errstr != "" {
		slog.ErrorContext(ctx, "[/admin.status/ping] request validate failed", slog.String("error", errstr))
		return nil, cerror.ErrReq
	}
	reqd, _ := proto.Marshal(req)
	var respbody []byte
	if e := c.cc.Call(ctx, _CrpcPathStatusPing, reqd, func(ctx *crpc.CallContext) error {
		var e error
		if respbody, e = ctx.Recv(); e != nil {
			slog.ErrorContext(ctx, "[/admin.status/ping] read response failed", slog.String("error", e.Error()))
		}
		return e
	}); e != nil {
		return nil, e
	}
	resp := new(Pingresp)
	if len(respbody) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respbody, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type StatusCrpcServer interface {
	// ping check server's health
	// Context is *crpc.NoStreamServerContext
	Ping(context.Context, *Pingreq) (*Pingresp, error)
}

func _Status_Ping_CrpcHandler(handler func(context.Context, *Pingreq) (*Pingresp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.ServerContext) {
		reqbody, e := ctx.Recv()
		if e != nil {
			slog.ErrorContext(ctx, "[/admin.status/ping] read request failed", slog.String("error", e.Error()))
			ctx.Abort(e)
			return
		}
		req := new(Pingreq)
		if e := proto.Unmarshal(reqbody, req); e != nil {
			slog.ErrorContext(ctx, "[/admin.status/ping] request decode failed", slog.String("error", e.Error()))
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			slog.ErrorContext(ctx, "[/admin.status/ping] request validate failed", slog.String("error", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(crpc.NewNoStreamServerContext(ctx), req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(Pingresp)
		}
		respd, _ := proto.Marshal(resp)
		if e := ctx.Send(respd); e != nil {
			slog.ErrorContext(ctx, "[/admin.status/ping] send response failed", slog.String("error", e.Error()))
		}
	}
}
func RegisterStatusCrpcServer(engine *crpc.CrpcServer, svc StatusCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("admin.status", "ping", _Status_Ping_CrpcHandler(svc.Ping))
}
