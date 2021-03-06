// Code generated by protoc-gen-go-web. DO NOT EDIT.
// version:
// 	protoc-gen-go-web v0.0.1
// 	protoc            v3.21.1
// source: api/user.proto

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

var _WebPathUserLogin = "/admin.user/login"
var _WebPathUserGetUsers = "/admin.user/get_users"
var _WebPathUserSearchUsers = "/admin.user/search_users"

type UserWebClient interface {
	Login(context.Context, *LoginReq, http.Header) (*LoginResp, error)
	GetUsers(context.Context, *GetUsersReq, http.Header) (*GetUsersResp, error)
	SearchUsers(context.Context, *SearchUsersReq, http.Header) (*SearchUsersResp, error)
}

type userWebClient struct {
	cc *web.WebClient
}

func NewUserWebClient(c *web.WebClient) UserWebClient {
	return &userWebClient{cc: c}
}

func (c *userWebClient) Login(ctx context.Context, req *LoginReq, header http.Header) (*LoginResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	data, e := c.cc.Post(ctx, _WebPathUserLogin, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	resp := new(LoginResp)
	if len(data) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(data, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *userWebClient) GetUsers(ctx context.Context, req *GetUsersReq, header http.Header) (*GetUsersResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	data, e := c.cc.Post(ctx, _WebPathUserGetUsers, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	resp := new(GetUsersResp)
	if len(data) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(data, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}
func (c *userWebClient) SearchUsers(ctx context.Context, req *SearchUsersReq, header http.Header) (*SearchUsersResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	data, e := c.cc.Post(ctx, _WebPathUserSearchUsers, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	resp := new(SearchUsersResp)
	if len(data) == 0 {
		return resp, nil
	}
	if e := proto.Unmarshal(data, resp); e != nil {
		return nil, error1.ErrResp
	}
	return resp, nil
}

type UserWebServer interface {
	Login(context.Context, *LoginReq) (*LoginResp, error)
	GetUsers(context.Context, *GetUsersReq) (*GetUsersResp, error)
	SearchUsers(context.Context, *SearchUsersReq) (*SearchUsersResp, error)
}

func _User_Login_WebHandler(handler func(context.Context, *LoginReq) (*LoginResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(LoginReq)
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
			data.AppendByte('}')
			if data.Len() > 2 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data.Bytes(), req)
				if e != nil {
					ctx.Abort(error1.ErrReq)
					return
				}
			}
		}
		resp, e := handler(ctx, req)
		ee := error1.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(LoginResp)
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
func _User_GetUsers_WebHandler(handler func(context.Context, *GetUsersReq) (*GetUsersResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(GetUsersReq)
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
			data.AppendString("\"user_ids\":")
			if forms := ctx.GetForms("user_ids"); len(forms) == 0 {
				data.AppendString("null")
			} else {
				data.AppendByte('[')
				for _, form := range forms {
					if len(form) == 0 {
						data.AppendString("\"\"")
					} else if len(form) < 2 || form[0] != '"' || form[len(form)-1] != '"' {
						data.AppendByte('"')
						data.AppendString(form)
						data.AppendByte('"')
					} else {
						data.AppendString(form)
					}
					data.AppendByte(',')
				}
				data.Bytes()[data.Len()-1] = ']'
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
			log.Error(ctx, "[/admin.user/get_users]", errstr)
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
			resp = new(GetUsersResp)
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
func _User_SearchUsers_WebHandler(handler func(context.Context, *SearchUsersReq) (*SearchUsersResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(SearchUsersReq)
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
			data.AppendString("\"user_name\":")
			if form := ctx.GetForm("user_name"); len(form) == 0 {
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
			log.Error(ctx, "[/admin.user/search_users]", errstr)
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
			resp = new(SearchUsersResp)
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
func RegisterUserWebServer(engine *web.WebServer, svc UserWebServer, allmids map[string]web.OutsideHandler) {
	//avoid lint
	_ = allmids
	engine.Post(_WebPathUserLogin, _User_Login_WebHandler(svc.Login))
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _User_GetUsers_WebHandler(svc.GetUsers))
		engine.Post(_WebPathUserGetUsers, mids...)
	}
	{
		requiredMids := []string{"token"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _User_SearchUsers_WebHandler(svc.SearchUsers))
		engine.Post(_WebPathUserSearchUsers, mids...)
	}
}
