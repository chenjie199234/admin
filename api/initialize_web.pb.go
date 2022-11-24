// Code generated by protoc-gen-go-web. DO NOT EDIT.
// version:
// 	protoc-gen-go-web v0.0.1
// 	protoc            v3.21.1
// source: api/initialize.proto

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

var _WebPathInitializeInit = "/admin.initialize/init"
var _WebPathInitializeRootLogin = "/admin.initialize/root_login"
var _WebPathInitializeRootPassword = "/admin.initialize/root_password"
var _WebPathInitializeCreateProject = "/admin.initialize/create_project"
var _WebPathInitializeUpdateProject = "/admin.initialize/update_project"
var _WebPathInitializeListProject = "/admin.initialize/list_project"
var _WebPathInitializeDeleteProject = "/admin.initialize/delete_project"

type InitializeWebClient interface {
	// 初始化
	Init(context.Context, *InitReq, http.Header) (*InitResp, error)
	// 登录
	RootLogin(context.Context, *RootLoginReq, http.Header) (*RootLoginResp, error)
	// 更新密码
	RootPassword(context.Context, *RootPasswordReq, http.Header) (*RootPasswordResp, error)
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
	} else if e := protojson.Unmarshal(data, resp); e != nil {
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
	} else if e := protojson.Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}
func (c *initializeWebClient) RootPassword(ctx context.Context, req *RootPasswordReq, header http.Header) (*RootPasswordResp, error) {
	if req == nil {
		return nil, cerror.ErrReq
	}
	if header == nil {
		header = make(http.Header)
	}
	header.Set("Content-Type", "application/x-protobuf")
	header.Set("Accept", "application/x-protobuf")
	reqd, _ := proto.Marshal(req)
	r, e := c.cc.Post(ctx, _WebPathInitializeRootPassword, "", header, metadata.GetMetadata(ctx), reqd)
	if e != nil {
		return nil, e
	}
	data, e := io.ReadAll(r.Body)
	r.Body.Close()
	if e != nil {
		return nil, cerror.ConvertStdError(e)
	}
	resp := new(RootPasswordResp)
	if len(data) == 0 {
		return resp, nil
	}
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/x-protobuf") {
		if e := proto.Unmarshal(data, resp); e != nil {
			return nil, cerror.ErrResp
		}
	} else if e := protojson.Unmarshal(data, resp); e != nil {
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
	} else if e := protojson.Unmarshal(data, resp); e != nil {
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
	} else if e := protojson.Unmarshal(data, resp); e != nil {
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
	} else if e := protojson.Unmarshal(data, resp); e != nil {
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
	} else if e := protojson.Unmarshal(data, resp); e != nil {
		return nil, cerror.ErrResp
	}
	return resp, nil
}

type InitializeWebServer interface {
	// 初始化
	Init(context.Context, *InitReq) (*InitResp, error)
	// 登录
	RootLogin(context.Context, *RootLoginReq) (*RootLoginResp, error)
	// 更新密码
	RootPassword(context.Context, *RootPasswordReq) (*RootPasswordResp, error)
	// 创建项目
	CreateProject(context.Context, *CreateProjectReq) (*CreateProjectResp, error)
	// 更新项目
	UpdateProject(context.Context, *UpdateProjectReq) (*UpdateProjectResp, error)
	// 获取项目列表
	ListProject(context.Context, *ListProjectReq) (*ListProjectResp, error)
	// 删除项目
	DeleteProject(context.Context, *DeleteProjectReq) (*DeleteProjectResp, error)
}

func _Initialize_Init_WebHandler(handler func(context.Context, *InitReq) (*InitResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(InitReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
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
			data.AppendString("\"password\":")
			if form := ctx.GetForm("password"); len(form) == 0 {
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
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/init]", errstr)
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
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
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
			data.AppendString("\"password\":")
			if form := ctx.GetForm("password"); len(form) == 0 {
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
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/root_login]", errstr)
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
func _Initialize_RootPassword_WebHandler(handler func(context.Context, *RootPasswordReq) (*RootPasswordResp, error)) web.OutsideHandler {
	return func(ctx *web.Context) {
		req := new(RootPasswordReq)
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
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
			data.AppendString("\"old_password\":")
			if form := ctx.GetForm("old_password"); len(form) == 0 {
				data.AppendString("\"\"")
			} else if len(form) < 2 || form[0] != '"' || form[len(form)-1] != '"' {
				data.AppendByte('"')
				data.AppendString(form)
				data.AppendByte('"')
			} else {
				data.AppendString(form)
			}
			data.AppendByte(',')
			data.AppendString("\"new_password\":")
			if form := ctx.GetForm("new_password"); len(form) == 0 {
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
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/root_password]", errstr)
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
			resp = new(RootPasswordResp)
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
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
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
			data.AppendString("\"project_name\":")
			if form := ctx.GetForm("project_name"); len(form) == 0 {
				data.AppendString("\"\"")
			} else if len(form) < 2 || form[0] != '"' || form[len(form)-1] != '"' {
				data.AppendByte('"')
				data.AppendString(form)
				data.AppendByte('"')
			} else {
				data.AppendString(form)
			}
			data.AppendByte(',')
			data.AppendString("\"project_data\":")
			if form := ctx.GetForm("project_data"); len(form) == 0 {
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
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/create_project]", errstr)
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
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
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
			data.AppendString("\"project_id\":")
			if forms := ctx.GetForms("project_id"); len(forms) == 0 {
				data.AppendString("null")
			} else {
				data.AppendByte('[')
				for _, form := range forms {
					if len(form) == 0 {
						data.AppendString("0")
					} else {
						data.AppendString(form)
					}
					data.AppendByte(',')
				}
				data.Bytes()[data.Len()-1] = ']'
			}
			data.AppendByte(',')
			data.AppendString("\"new_project_name\":")
			if form := ctx.GetForm("new_project_name"); len(form) == 0 {
				data.AppendString("\"\"")
			} else if len(form) < 2 || form[0] != '"' || form[len(form)-1] != '"' {
				data.AppendByte('"')
				data.AppendString(form)
				data.AppendByte('"')
			} else {
				data.AppendString(form)
			}
			data.AppendByte(',')
			data.AppendString("\"new_project_data\":")
			if form := ctx.GetForm("new_project_data"); len(form) == 0 {
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
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/update_project]", errstr)
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
		if strings.HasPrefix(ctx.GetContentType(), "application/json") {
			data, e := ctx.GetBody()
			if e != nil {
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
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
			data := pool.GetBuffer()
			defer pool.PutBuffer(data)
			data.AppendByte('{')
			data.AppendByte('}')
			if data.Len() > 2 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data.Bytes(), req)
				if e != nil {
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		}
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
				ctx.Abort(e)
				return
			}
			if len(data) > 0 {
				e := protojson.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}.Unmarshal(data, req)
				if e != nil {
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
			data.AppendString("\"project_id\":")
			if forms := ctx.GetForms("project_id"); len(forms) == 0 {
				data.AppendString("null")
			} else {
				data.AppendByte('[')
				for _, form := range forms {
					if len(form) == 0 {
						data.AppendString("0")
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
					ctx.Abort(cerror.ErrReq)
					return
				}
			}
		}
		if errstr := req.Validate(); errstr != "" {
			log.Error(ctx, "[/admin.initialize/delete_project]", errstr)
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
func RegisterInitializeWebServer(engine *web.WebServer, svc InitializeWebServer, allmids map[string]web.OutsideHandler) {
	// avoid lint
	_ = allmids
	engine.Post(_WebPathInitializeInit, _Initialize_Init_WebHandler(svc.Init))
	engine.Post(_WebPathInitializeRootLogin, _Initialize_RootLogin_WebHandler(svc.RootLogin))
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
		mids = append(mids, _Initialize_RootPassword_WebHandler(svc.RootPassword))
		engine.Post(_WebPathInitializeRootPassword, mids...)
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
		engine.Post(_WebPathInitializeCreateProject, mids...)
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
		engine.Post(_WebPathInitializeUpdateProject, mids...)
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
		engine.Post(_WebPathInitializeListProject, mids...)
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
		engine.Post(_WebPathInitializeDeleteProject, mids...)
	}
}
