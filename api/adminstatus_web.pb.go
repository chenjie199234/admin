// Code generated by protoc-gen-go-web. DO NOT EDIT.
// version:
// 	protoc-gen-go-web v0.0.76
// 	protoc            v3.21.11
// source: api/adminstatus.proto

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	pool "github.com/chenjie199234/Corelib/pool"
	web "github.com/chenjie199234/Corelib/web"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	io "io"
	http "net/http"
	strings "strings"
)

var _WebPathStatusPing = "/admin.status/ping"

type StatusWebClient interface {
	// ping check server's health
	Ping(context.Context, *Pingreq, http.Header) (*Pingresp, error)
}

type statusWebClient struct {
	cc *web.WebClient
}

func NewStatusWebClient(c *web.WebClient) StatusWebClient {
	return &statusWebClient{cc: c}
}

func (c *statusWebClient) Ping(ctx context.Context, req *Pingreq, header http.Header) (*Pingresp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	header.Set("Accept", "application/x-protobuf")
	query := pool.GetBuffer()
	defer pool.PutBuffer(query)
	if req.GetTimestamp() != 0 {
		query.AppendString("timestamp=")
		query.AppendInt64(req.GetTimestamp())
		query.AppendByte('&')
	}
	querystr := query.String()
	if len(querystr) > 0 {
		// drop last &
		querystr = querystr[:len(querystr)-1]
	}
	r, e := c.cc.Get(ctx, _WebPathStatusPing, querystr, header, metadata.GetMetadata(ctx))
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(Pingresp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type StatusWebServer interface {
	// ping check server's health
	Ping(context.Context, *Pingreq) (*Pingresp, error)
}

func _Status_Ping_WebHandler(handler func(context.Context, *Pingreq) (*Pingresp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(Pingreq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					ctx.Abort(cerror.ErrReq)
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
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				ctx.Abort(cerror.ErrReq)
				return
			}
			data := pool.GetBuffer()
			defer pool.PutBuffer(data)
			data.AppendByte('{')
			if form := ctx.GetForm("timestamp"); len(form) != 0 {
				data.AppendString("\"timestamp\":")
				data.AppendString(form)
				data.AppendByte(',')
			}
			if data.Len() == 1 {
				data.AppendByte('}')
			} else {
				data.Bytes()[data.Len()-1] = '}'
			}
			if data.Len() > 2 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data.Bytes(), req); e != nil {
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.status/ping]", errstr)
			ctx.Abort(cerror.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(Pingresp)
		}
		if strings.HasPrefix(ctx.GetAcceptType(), "application/x-protobuf") {
			respd, _ := proto.Marshal(resp)
			ctx.Write("application/x-protobuf", respd)
		} else {
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func RegisterStatusWebServer(engine *web.WebServer, svc StatusWebServer, allmids map[string]web.OutsideHandler) {
	// avoid lint
	_ = allmids
	{
		requiredMids := []string{"accesskey", "rate"}
		mids := make([]web.OutsideHandler, 0, 3)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Status_Ping_WebHandler(svc.Ping))
		engine.Get(_WebPathStatusPing, mids...)
	}
}
