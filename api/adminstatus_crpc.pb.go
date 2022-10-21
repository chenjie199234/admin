// Code generated by protoc-gen-go-crpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-crpc v0.0.1
// 	protoc             v3.21.1
// source: api/adminstatus.proto

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	crpc "github.com/chenjie199234/Corelib/crpc"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	proto "google.golang.org/protobuf/proto"
)

var _CrpcPathStatusPing = "/admin.status/ping"

type StatusCrpcClient interface {
	// ping check server's health
	Ping(context.Context, *Pingreq) (*Pingresp, error)
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
	reqd, _ := proto.Marshal(req)
	respd, e := c.cc.Call(ctx, _CrpcPathStatusPing, reqd, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	resp := new(Pingresp)
	if len(respd) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(respd, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type StatusCrpcServer interface {
	// ping check server's health
	Ping(context.Context, *Pingreq) (*Pingresp, error)
}

func _Status_Ping_CrpcHandler(handler func(context.Context, *Pingreq) (*Pingresp, error)) crpc.OutsideHandler {
	return func(ctx *crpc.Context) {
		req := new(Pingreq)
		if e := proto.Unmarshal(ctx.GetBody(), req); e != nil {
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
		respd, _ := proto.Marshal(resp)
		ctx.Write(respd)
	}
}
func RegisterStatusCrpcServer(engine *crpc.CrpcServer, svc StatusCrpcServer, allmids map[string]crpc.OutsideHandler) {
	// avoid lint
	_ = allmids
	{
		requiredMids := []string{"accesskey", "rate"}
		mids := make([]crpc.OutsideHandler, 0, 3)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Status_Ping_CrpcHandler(svc.Ping))
		engine.RegisterHandler(_CrpcPathStatusPing, mids...)
	}
}
