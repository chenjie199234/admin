// Code generated by protoc-gen-go-web. DO NOT EDIT.
// version:
// 	protoc-gen-go-web v0.0.1
// 	protoc            v3.21.1
// source: api/initialize.proto

package api

import (
	context "context"
	error1 "github.com/chenjie199234/Corelib/error"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	pool "github.com/chenjie199234/Corelib/pool"
	web "github.com/chenjie199234/Corelib/web"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	http "net/http"
	strings "strings"
)

var _WebPathInitializeInitialize = "/admin.initialize/initialize"

type InitializeWebClient interface {
	//初始化
	Initialize(context.Context, *InitializeReq, http.Header) (*InitializeResp, error)
}

type initializeWebClient struct {
	cc *web.WebClient
}

func NewInitializeWebClient(c *web.WebClient) InitializeWebClient {
	return &initializeWebClient{cc: c}
}

func (c *initializeWebClient) Initialize(ctx context.Context, req *InitializeReq, header http.Header) (*InitializeResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	data, e := c.cc.Post(ctx, _WebPathInitializeInitialize, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	resp := new(InitializeResp)
	if len(data) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(data, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}

type InitializeWebServer interface {
	//初始化
	Initialize(context.Context, *InitializeReq) (*InitializeResp, error)
}

func _Initialize_Initialize_WebHandler(handler func(context.Context, *InitializeReq) (*InitializeResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(InitializeReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				ctx.Abort(error1.ErrReq)
				return
			}
			data := pool.GetBuffer()
			defer pool.PutBuffer(data)
			data.AppendByte('{')
			data.AppendString("\"super_admin_password\":")
			if form := ctx.GetForm("super_admin_password"); len(form) == 0 {
				data.AppendString("\"\"")
			} else if len(form) < 2 || form[0] != '"' || form[len(form)-1] != '"' {
				data.AppendByte('"')
				data.AppendString(form)
				data.AppendByte('"')
			} else {
				data.AppendString(form)
			}
			data.AppendByte('}')
			if data.Len() > 2 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data.Bytes(), req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/initialize]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := error1.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(InitializeResp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func RegisterInitializeWebServer(engine *web.WebServer, svc InitializeWebServer, allmids map[string]web.OutsideHandler) {
	//avoid lint
	_ = allmids
	engine.Post(_WebPathInitializeInitialize, _Initialize_Initialize_WebHandler(svc.Initialize))
}
