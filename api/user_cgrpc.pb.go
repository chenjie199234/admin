// Code generated by protoc-gen-go-cgrpc. DO NOT EDIT.
// version:
// 	protoc-gen-go-cgrpc v0.0.1
// 	protoc              v3.21.1
// source: api/user.proto

package api

import (
	context "context"
	cgrpc "github.com/chenjie199234/Corelib/cgrpc"
	error1 "github.com/chenjie199234/Corelib/error"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
)

var _CGrpcPathUserLogin = "/admin.user/login"
var _CGrpcPathUserGetUsers = "/admin.user/get_users"
var _CGrpcPathUserSearchUsers = "/admin.user/search_users"

type UserCGrpcClient interface {
	Login(context.Context, *LoginReq) (*LoginResp, error)
	GetUsers(context.Context, *GetUsersReq) (*GetUsersResp, error)
	SearchUsers(context.Context, *SearchUsersReq) (*SearchUsersResp, error)
}

type userCGrpcClient struct {
	cc *cgrpc.CGrpcClient
}

func NewUserCGrpcClient(c *cgrpc.CGrpcClient) UserCGrpcClient {
	return &userCGrpcClient{cc: c}
}

func (c *userCGrpcClient) Login(ctx context.Context, req *LoginReq) (*LoginResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(LoginResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserLogin, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) GetUsers(ctx context.Context, req *GetUsersReq) (*GetUsersResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(GetUsersResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserGetUsers, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}
func (c *userCGrpcClient) SearchUsers(ctx context.Context, req *SearchUsersReq) (*SearchUsersResp, error) {
	if req == nil {
		return nil, error1.ErrReq
	}
	resp := new(SearchUsersResp)
	if e := c.cc.Call(ctx, _CGrpcPathUserSearchUsers, req, resp, metadata.GetMetadata(ctx)); e != nil {
		return nil, e
	}
	return resp, nil
}

type UserCGrpcServer interface {
	Login(context.Context, *LoginReq) (*LoginResp, error)
	GetUsers(context.Context, *GetUsersReq) (*GetUsersResp, error)
	SearchUsers(context.Context, *SearchUsersReq) (*SearchUsersResp, error)
}

func _User_Login_CGrpcHandler(handler func(context.Context, *LoginReq) (*LoginResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(LoginReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(LoginResp)
		}
		ctx.Write(resp)
	}
}
func _User_GetUsers_CGrpcHandler(handler func(context.Context, *GetUsersReq) (*GetUsersResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(GetUsersReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/get_users]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(GetUsersResp)
		}
		ctx.Write(resp)
	}
}
func _User_SearchUsers_CGrpcHandler(handler func(context.Context, *SearchUsersReq) (*SearchUsersResp, error)) cgrpc.OutsideHandler {
	return func(ctx *cgrpc.Context) {
		req := new(SearchUsersReq)
		if ctx.DecodeReq(req) != nil {
			ctx.Abort(error1.ErrReq)
			return
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.user/search_users]", errstr)
			ctx.Abort(error1.ErrReq)
			return
		}
		resp, e := handler(ctx, req)
		if e != nil {
			ctx.Abort(e)
			return
		}
		if resp == nil {
			resp = new(SearchUsersResp)
		}
		ctx.Write(resp)
	}
}
func RegisterUserCGrpcServer(engine *cgrpc.CGrpcServer, svc UserCGrpcServer, allmids map[string]cgrpc.OutsideHandler) {
	//avoid lint
	_ = allmids
	engine.RegisterHandler("admin.user", "login", _User_Login_CGrpcHandler(svc.Login))
	engine.RegisterHandler("admin.user", "get_users", _User_GetUsers_CGrpcHandler(svc.GetUsers))
	engine.RegisterHandler("admin.user", "search_users", _User_SearchUsers_CGrpcHandler(svc.SearchUsers))
}
