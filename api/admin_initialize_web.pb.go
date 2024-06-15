// Code generated by protoc-gen-go-web. DO NOT EDIT.
// version:
// 	protoc-gen-go-web v0.0.116<br />
// 	protoc            v5.27.0<br />
// source: api/admin_initialize.proto<br />

package api

import (
	context "context"
	cerror "github.com/chenjie199234/Corelib/cerror"
	log "github.com/chenjie199234/Corelib/log"
	metadata "github.com/chenjie199234/Corelib/metadata"
	web "github.com/chenjie199234/Corelib/web"
	protojson "google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	io "io"
	http "net/http"
	strconv "strconv"
	strings "strings"
)

var _WebPathInitializeInitStatus = "/admin.initialize/init_status"
var _WebPathInitializeInit = "/admin.initialize/init"
var _WebPathInitializeRootLogin = "/admin.initialize/root_login"
var _WebPathInitializeUpdateRootPassword = "/admin.initialize/update_root_password"
var _WebPathInitializeCreateProject = "/admin.initialize/create_project"
var _WebPathInitializeUpdateProject = "/admin.initialize/update_project"
var _WebPathInitializeListProject = "/admin.initialize/list_project"
var _WebPathInitializeDeleteProject = "/admin.initialize/delete_project"

type InitializeWebClient interface {
	// 初始化状态
	InitStatus(context.Context, *InitStatusReq, http.Header) (*InitStatusResp, error)
	// 初始化
	Init(context.Context, *InitReq, http.Header) (*InitResp, error)
	// 登录
	RootLogin(context.Context, *RootLoginReq, http.Header) (*RootLoginResp, error)
	// 更新密码
	UpdateRootPassword(context.Context, *UpdateRootPasswordReq, http.Header) (*UpdateRootPasswordResp, error)
	// 创建项目
	CreateProject(context.Context, *CreateProjectReq, http.Header) (*CreateProjectResp, error)
	// 更新项目
	UpdateProject(context.Context, *UpdateProjectReq, http.Header) (*UpdateProjectResp, error)
	// 获取项目列表
	ListProject(context.Context, *ListProjectReq, http.Header) (*ListProjectResp, error)
	// 删除项目
	DeleteProject(context.Context, *DeleteProjectReq, http.Header) (*DeleteProjectResp, error)
}

type initializeWebClient struct {
	cc *web.WebClient
}

func NewInitializeWebClient(c *web.WebClient) InitializeWebClient {
	return &initializeWebClient{cc: c}
}

func (c *initializeWebClient) InitStatus(ctx context.Context, req *InitStatusReq, header http.Header) (*InitStatusResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathInitializeInitStatus, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(InitStatusResp)
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
func (c *initializeWebClient) Init(ctx context.Context, req *InitReq, header http.Header) (*InitResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathInitializeInit, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(InitResp)
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
func (c *initializeWebClient) RootLogin(ctx context.Context, req *RootLoginReq, header http.Header) (*RootLoginResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathInitializeRootLogin, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(RootLoginResp)
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
func (c *initializeWebClient) UpdateRootPassword(ctx context.Context, req *UpdateRootPasswordReq, header http.Header) (*UpdateRootPasswordResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathInitializeUpdateRootPassword, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(UpdateRootPasswordResp)
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
func (c *initializeWebClient) CreateProject(ctx context.Context, req *CreateProjectReq, header http.Header) (*CreateProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathInitializeCreateProject, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(CreateProjectResp)
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
func (c *initializeWebClient) UpdateProject(ctx context.Context, req *UpdateProjectReq, header http.Header) (*UpdateProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathInitializeUpdateProject, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(UpdateProjectResp)
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
func (c *initializeWebClient) ListProject(ctx context.Context, req *ListProjectReq, header http.Header) (*ListProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathInitializeListProject, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(ListProjectResp)
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
func (c *initializeWebClient) DeleteProject(ctx context.Context, req *DeleteProjectReq, header http.Header) (*DeleteProjectResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathInitializeDeleteProject, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(DeleteProjectResp)
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

type InitializeWebServer interface {
	// 初始化状态
	InitStatus(context.Context, *InitStatusReq) (*InitStatusResp, error)
	// 初始化
	Init(context.Context, *InitReq) (*InitResp, error)
	// 登录
	RootLogin(context.Context, *RootLoginReq) (*RootLoginResp, error)
	// 更新密码
	UpdateRootPassword(context.Context, *UpdateRootPasswordReq) (*UpdateRootPasswordResp, error)
	// 创建项目
	CreateProject(context.Context, *CreateProjectReq) (*CreateProjectResp, error)
	// 更新项目
	UpdateProject(context.Context, *UpdateProjectReq) (*UpdateProjectResp, error)
	// 获取项目列表
	ListProject(context.Context, *ListProjectReq) (*ListProjectResp, error)
	// 删除项目
	DeleteProject(context.Context, *DeleteProjectReq) (*DeleteProjectResp, error)
}

func _Initialize_InitStatus_WebHandler(handler func(context.Context, *InitStatusReq) (*InitStatusResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(InitStatusReq)
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(InitStatusResp)
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
func _Initialize_Init_WebHandler(handler func(context.Context, *InitReq) (*InitResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(InitReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/init] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/init] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/init] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/init] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				log.Error(ctx, "[/admin.initialize/init] parse form failed", log.CError(e))
				ctx.Abort(cerror.ErrReq)
				return
			}
			// req.Password
			if form := ctx.GetForm("password"); len(form) != 0 {
				req.Password = form
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/init] validate failed", log.String("validate", errstr))
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
			resp = new(InitResp)
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
func _Initialize_RootLogin_WebHandler(handler func(context.Context, *RootLoginReq) (*RootLoginResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(RootLoginReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/root_login] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/root_login] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/root_login] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/root_login] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				log.Error(ctx, "[/admin.initialize/root_login] parse form failed", log.CError(e))
				ctx.Abort(cerror.ErrReq)
				return
			}
			// req.Password
			if form := ctx.GetForm("password"); len(form) != 0 {
				req.Password = form
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/root_login] validate failed", log.String("validate", errstr))
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
			resp = new(RootLoginResp)
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
func _Initialize_UpdateRootPassword_WebHandler(handler func(context.Context, *UpdateRootPasswordReq) (*UpdateRootPasswordResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateRootPasswordReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/update_root_password] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/update_root_password] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/update_root_password] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/update_root_password] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				log.Error(ctx, "[/admin.initialize/update_root_password] parse form failed", log.CError(e))
				ctx.Abort(cerror.ErrReq)
				return
			}
			// req.OldPassword
			if form := ctx.GetForm("old_password"); len(form) != 0 {
				req.OldPassword = form
			}
			// req.NewPassword
			if form := ctx.GetForm("new_password"); len(form) != 0 {
				req.NewPassword = form
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/update_root_password] validate failed", log.String("validate", errstr))
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
			resp = new(UpdateRootPasswordResp)
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
func _Initialize_CreateProject_WebHandler(handler func(context.Context, *CreateProjectReq) (*CreateProjectResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(CreateProjectReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/create_project] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/create_project] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/create_project] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/create_project] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				log.Error(ctx, "[/admin.initialize/create_project] parse form failed", log.CError(e))
				ctx.Abort(cerror.ErrReq)
				return
			}
			// req.ProjectName
			if form := ctx.GetForm("project_name"); len(form) != 0 {
				req.ProjectName = form
			}
			// req.ProjectData
			if form := ctx.GetForm("project_data"); len(form) != 0 {
				req.ProjectData = form
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/create_project] validate failed", log.String("validate", errstr))
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
			resp = new(CreateProjectResp)
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
func _Initialize_UpdateProject_WebHandler(handler func(context.Context, *UpdateProjectReq) (*UpdateProjectResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(UpdateProjectReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/update_project] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/update_project] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/update_project] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/update_project] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				log.Error(ctx, "[/admin.initialize/update_project] parse form failed", log.CError(e))
				ctx.Abort(cerror.ErrReq)
				return
			}
			// req.ProjectId
			if forms := ctx.GetForms("project_id"); len(forms) > 0 {
				req.ProjectId = make([]uint32, 0, len(forms))
				for _, form := range forms {
					if num, e := strconv.ParseUint(form, 10, 32); e != nil {
						log.Error(ctx, "[/admin.initialize/update_project] data format wrong", log.String("field", "project_id"))
						ctx.Abort(cerror.ErrReq)
						return
					} else {
						req.ProjectId = append(req.ProjectId, uint32(num))
					}
				}
			}
			// req.NewProjectName
			if form := ctx.GetForm("new_project_name"); len(form) != 0 {
				req.NewProjectName = form
			}
			// req.NewProjectData
			if form := ctx.GetForm("new_project_data"); len(form) != 0 {
				req.NewProjectData = form
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/update_project] validate failed", log.String("validate", errstr))
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
			resp = new(UpdateProjectResp)
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
func _Initialize_ListProject_WebHandler(handler func(context.Context, *ListProjectReq) (*ListProjectResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(ListProjectReq)
		resp, e := handler(ctx, req)
		ee := cerror.ConvertStdError(e)
		if ee != nil {
			ctx.Abort(ee)
			return
		}
		if resp == nil {
			resp = new(ListProjectResp)
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
func _Initialize_DeleteProject_WebHandler(handler func(context.Context, *DeleteProjectReq) (*DeleteProjectResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(DeleteProjectReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/delete_project] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := (protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}).Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/delete_project] unmarshal json body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else if strings.HasPrefix(ctx.GetContentType(), "application/x-protobuf") {
			data, e := ctx.GetBody()
			if e != nil {
				log.Error(ctx, "[/admin.initialize/delete_project] get body failed", log.CError(e))
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				if e := proto.Unmarshal(data, req); e != nil {
					log.Error(ctx, "[/admin.initialize/delete_project] unmarshal proto body failed", log.CError(e))
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		} else {
			if e := ctx.ParseForm(); e != nil {
				log.Error(ctx, "[/admin.initialize/delete_project] parse form failed", log.CError(e))
				ctx.Abort(cerror.ErrReq)
				return
			}
			// req.ProjectId
			if forms := ctx.GetForms("project_id"); len(forms) > 0 {
				req.ProjectId = make([]uint32, 0, len(forms))
				for _, form := range forms {
					if num, e := strconv.ParseUint(form, 10, 32); e != nil {
						log.Error(ctx, "[/admin.initialize/delete_project] data format wrong", log.String("field", "project_id"))
						ctx.Abort(cerror.ErrReq)
						return
					} else {
						req.ProjectId = append(req.ProjectId, uint32(num))
					}
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/delete_project] validate failed", log.String("validate", errstr))
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
			resp = new(DeleteProjectResp)
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
func RegisterInitializeWebServer(router *web.Router, svc InitializeWebServer, allmids map[string]web.OutsideHandler) {
	// avoid lint
	_ = allmids
	router.Post(_WebPathInitializeInitStatus, _Initialize_InitStatus_WebHandler(svc.InitStatus))
	{
		requiredMids := []string{"accesskey"}
		mids := make([]web.OutsideHandler, 0, 2)
		for _, v := range requiredMids {
			if mid, ok := allmids[v]; ok {
				mids = append(mids, mid)
			} else {
				panic("missing midware:" + v)
			}
		}
		mids = append(mids, _Initialize_Init_WebHandler(svc.Init))
		router.Post(_WebPathInitializeInit, mids...)
	}
	router.Post(_WebPathInitializeRootLogin, _Initialize_RootLogin_WebHandler(svc.RootLogin))
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
		mids = append(mids, _Initialize_UpdateRootPassword_WebHandler(svc.UpdateRootPassword))
		router.Post(_WebPathInitializeUpdateRootPassword, mids...)
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
		mids = append(mids, _Initialize_CreateProject_WebHandler(svc.CreateProject))
		router.Post(_WebPathInitializeCreateProject, mids...)
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
		mids = append(mids, _Initialize_UpdateProject_WebHandler(svc.UpdateProject))
		router.Post(_WebPathInitializeUpdateProject, mids...)
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
		mids = append(mids, _Initialize_ListProject_WebHandler(svc.ListProject))
		router.Post(_WebPathInitializeListProject, mids...)
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
		mids = append(mids, _Initialize_DeleteProject_WebHandler(svc.DeleteProject))
		router.Post(_WebPathInitializeDeleteProject, mids...)
	}
}
