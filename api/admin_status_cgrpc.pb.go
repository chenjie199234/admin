// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.134<br />
// 	protoc              v6.30.2<br />
// source: api/admin_status.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	grpc "google.golang.org/grpc"
	slog "log/slog"
)

var _CGrpcPathStatusPing = "/admin.status/ping"

type StatusCGrpcClient interface {
	// ping check server's health
	Ping(context.Context, *Pingreq) (*Pingresp, error)
}

type statusCGrpcClient struct {
	cc grpc.ClientConnInterface
}

func NewStatusCGrpcClient(cc grpc.ClientConnInterface) StatusCGrpcClient {
	return &statusCGrpcClient{cc: cc}
}

func (c *statusCGrpcClient) Ping(ctx context.Context, req *Pingreq) (*Pingresp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if errstr := req.Validate(); errstr != "" {
		slog.ErrorContext(ctx, "[/admin.status/ping] validate failed", slog.String("error", errstr))
		return nil, cerror.ErrReq
	}
	resp := new(Pingresp)
	if e := c.cc.Invoke(ctx, _CGrpcPathStatusPing, req, resp); e != nil {
		return nil, e
	}
	return resp, nil
}

type StatusCGrpcServer interface {
	// ping check server's health
	// Context is *cgrpc.NoStreamServerContext
	Ping(context.Context, *Pingreq) (*Pingresp, error)
}

func _Status_Ping_CGrpcHandler(handler func(context.Context, *Pingreq) (*Pingresp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.ServerContext) {
		req := new(Pingreq)
		if e := ctx.Read(req); e != nil {
			slog.ErrorContext(ctx, "[/admin.status/ping] decode failed", slog.String("error", e.Error()))
			ctx.Abort(cerror.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			slog.ErrorContext(ctx, "[/admin.status/ping] validate failed", slog.String("error", errstr))
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(cgrpc.NewNoStreamServerContext(ctx), req)
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
	engine.RegisterHandler("admin.status", "ping", false, false, _Status_Ping_CGrpcHandler(svc.Ping))
}
