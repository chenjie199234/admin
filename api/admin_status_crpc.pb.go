// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.134<br />
// 	protoc             v6.30.1<br />
// source: api/admin_status.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	protojson "google.golang.org/protobuf/encoding/protojson"
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
	var encoder crpc.Encoder
	if e := c.cc.Call(ctx, _CrpcPathStatusPing, reqd, crpc.Encoder_Protobuf, func(ctx *crpc.CallContext) error {
		var e error
		if respbody, encoder, e = ctx.Recv(); e != nil {
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
	switch encoder {
	case crpc.Encoder_Protobuf:
		if e := proto.Unmarshal(respbody, resp); e != nil {
			return nil, cerror.ErrResp
		}
	case crpc.Encoder_Json:
		if e := protojson.Unmarshal(respbody, resp); e != nil {
			return nil, cerror.ErrResp
		}
	default:
		slog.ErrorContext(ctx, "[/admin.status/ping] unknown response encoder")
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
		reqbody, encoder, e := ctx.Recv()
		if e != nil {
			slog.ErrorContext(ctx, "[/admin.status/ping] read request failed", slog.String("error", e.Error()))
			ctx.Abort(e)
			return
		}
		req := new(Pingreq)
		switch encoder {
		case crpc.Encoder_Protobuf:
			if e := proto.Unmarshal(reqbody, req); e != nil {
				slog.ErrorContext(ctx, "[/admin.status/ping] request decode failed", slog.String("error", e.Error()))
				ctx.Abort(cerror.ErrReq)
				return
			}
		case crpc.Encoder_Json:
			if e := protojson.Unmarshal(reqbody, req); e != nil {
				slog.ErrorContext(ctx, "[/admin.status/ping] request decode failed", slog.String("error", e.Error()))
				ctx.Abort(cerror.ErrReq)
				return
			}
		default:
			slog.ErrorContext(ctx, "[/admin.status/ping] request encoder unknown")
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
		var respd []byte
		switch encoder {
		case crpc.Encoder_Protobuf:
			respd, _ = proto.Marshal(resp)
		case crpc.Encoder_Json:
			respd, _ = protojson.Marshal(resp)
		}
		if e := ctx.Send(respd, encoder); e != nil {
			slog.ErrorContext(ctx, "[/admin.status/ping] send response failed", slog.String("error", e.Error()))
		}
	}
}
func RegisterStatusCrpcServer(engine *crpc.CrpcServer, svc StatusCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.RegisterHandler("admin.status", "ping", _Status_Ping_CrpcHandler(svc.Ping))
}
