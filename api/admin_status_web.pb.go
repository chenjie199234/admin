// Code generated by protoc-gen-go-web. DO NOT EDIT.
// version:
// 	protoc-gen-go-web v0.0.97<br />
// 	protoc            v4.25.3<br />
// source: api/admin_status.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	pool "github.com/chenjie199234/Corelib/pool"
	common "github.com/chenjie199234/Corelib/util/common"
	web "github.com/chenjie199234/Corelib/web"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	io "io"
	http "net/http"
	strconv "strconv"
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
	query := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&query)
	// req.Timestamp
	query = pool.CheckCap(&query, len(query)+9+22)
	query = append(query, "timestamp="...)
	query = strconv.AppendInt(query, req.GetTimestamp(), 10)
	query = append(query, '&')
	if len(query) > 0 {
		// drop last &
		query = query[:len(query)-1]
	}
	querystr := common.BTS(query)
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
		if e := ctx.ParseForm(); e != nil {
			log.Error(ctx, "[/admin.status/ping] parse form failed", log.CError(e))
			ctx.Abort(cerror.ErrReq)
			return
		}
		// req.Timestamp
		if form := ctx.GetForm("timestamp"); len(form) != 0 {
			if num, e := strconv.ParseInt(form, 10, 64); e != nil {
				log.Error(ctx, "[/admin.status/ping] data format wrong", log.String("field", "timestamp"))
				ctx.Abort(cerror.ErrReq)
				return
			} else {
				req.Timestamp = num
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.status/ping] validate failed", log.String("validate", errstr))
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
			respd, _ := protojson.MarshalOptions{AllowPartial: true, UseProtoNames: true, UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(resp)
			ctx.Write("application/json", respd)
		}
	}
}
func RegisterStatusWebServer(router *web.Router, svc StatusWebServer, allmids map[string]web.OutsideHandler) {
	// avoid lint
	_ = allmids
	router.Get(_WebPathStatusPing, _Status_Ping_WebHandler(svc.Ping))
}
