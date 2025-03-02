package permission

import (
	"context"
	"crypto/tls"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/ecode"

	"github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/util/egroup"
	"github.com/chenjie199234/Corelib/util/name"
	"github.com/chenjie199234/Corelib/web"
)

type PermissionSdk struct {
	client    api.PermissionWebClient
	accesskey string
}

var (
	ErrMissingEnvProject   = errors.New("missing env ADMIN_SERVICE_PROJECT")
	ErrMissingEnvGroup     = errors.New("missing env ADMIN_SERVICE_GROUP")
	ErrMissingEnvHost      = errors.New("missing env ADMIN_SERVICE_WEB_HOST")
	ErrMissingEnvAccessKey = errors.New("missing env ADMIN_SERVICE_PERMISSION_ACCESS_KEY")
	ErrWrongEnvPort        = errors.New("env ADMIN_SERVICE_WEB_PORT must be number <= 65535")
)

// if tlsc is not nil,the tls will be actived
// required env:
// ADMIN_SERVICE_PROJECT
// ADMIN_SERVICE_GROUP
// ADMIN_SERVICE_WEB_HOST
// ADMIN_SERVICE_WEB_PORT
// ADMIN_SERVICE_PERMISSION_ACCESS_KEY
func NewPermissionSdk(tlsc *tls.Config) (*PermissionSdk, error) {
	if e := name.HasSelfFullName(); e != nil {
		slog.Error("new admin permission sdk failed,please call github.com/chenjie199234/admin/sdk.Init() first")
		return nil, e
	}
	project, group, host, port, accesskey, e := env()
	if e != nil {
		return nil, e
	}
	di, e := discover.NewStaticDiscover(project, group, "admin", []string{host}, 0, 0, port)
	if e != nil {
		return nil, e
	}
	tmpclient, e := web.NewWebClient(nil, di, project, group, "admin", tlsc)
	if e != nil {
		return nil, e
	}
	return &PermissionSdk{client: api.NewPermissionWebClient(tmpclient), accesskey: accesskey}, nil
}
func env() (projectname, group string, host string, port int, accesskey string, e error) {
	if str, ok := os.LookupEnv("ADMIN_SERVICE_PROJECT"); ok && str != "<ADMIN_SERVICE_PROJECT>" && str != "" {
		projectname = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvProject
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_GROUP"); ok && str != "<ADMIN_SERVICE_GROUP>" && str != "" {
		group = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvGroup
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_WEB_HOST"); ok && str != "<ADMIN_SERVICE_WEB_HOST>" && str != "" {
		host = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvHost
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_WEB_PORT"); ok && str != "<ADMIN_SERVICE_WEB_PORT>" && str != "" {
		var e error
		port, e = strconv.Atoi(str)
		if e != nil || port < 0 || port > 65535 {
			return "", "", "", 0, "", ErrWrongEnvPort
		}
	}
	if str, ok := os.LookupEnv("ADMIN_SERVICE_PERMISSION_ACCESS_KEY"); ok && str != "<ADMIN_SERVICE_PERMISSION_ACCESS_KEY>" && str != "" {
		accesskey = str
	} else {
		return "", "", "", 0, "", ErrMissingEnvAccessKey
	}
	return
}

// if pass will return nil
// if not pass will reutrn ecode.ErrPermission(use cerror.Equal to check)(https://github.com/chenjie199234/Corelib/tree/main/cerror)
func (s *PermissionSdk) CheckMulti(ctx context.Context, userid string, readNodeIDs [][]uint32, writeNodeIDs [][]uint32, adminNodeIDs [][]uint32) error {
	eg := egroup.GetGroup(ctx)
	for _, v := range readNodeIDs {
		nodeid := v
		eg.Go(func(gctx context.Context) error {
			return s.CheckRead(gctx, userid, nodeid)
		})
	}
	for _, v := range writeNodeIDs {
		nodeid := v
		eg.Go(func(gctx context.Context) error {
			return s.CheckWrite(gctx, userid, nodeid)
		})
	}
	for _, v := range adminNodeIDs {
		nodeid := v
		eg.Go(func(gctx context.Context) error {
			return s.CheckAdmin(gctx, userid, nodeid)
		})
	}
	return egroup.PutGroup(eg)
}

// if pass will return nil
// if not pass will reutrn ecode.ErrPermission(use cerror.Equal to check)(https://github.com/chenjie199234/Corelib/tree/main/cerror)
func (s *PermissionSdk) CheckAdmin(ctx context.Context, userid string, nodeid []uint32) error {
	req := &api.GetUserPermissionReq{
		UserId: userid,
		NodeId: nodeid,
	}
	header := make(http.Header)
	header.Set("Access-Key", s.accesskey)
	resp, e := s.client.GetUserPermission(ctx, req, header)
	if e != nil {
		return e
	}
	if resp.GetAdmin() {
		return nil
	}
	return ecode.ErrPermission
}

// if pass will return nil
// if not pass will reutrn ecode.ErrPermission(use cerror.Equal to check)(https://github.com/chenjie199234/Corelib/tree/main/cerror)
func (s *PermissionSdk) CheckRead(ctx context.Context, userid string, nodeid []uint32) error {
	req := &api.GetUserPermissionReq{
		UserId: userid,
		NodeId: nodeid,
	}
	header := make(http.Header)
	header.Set("Access-Key", s.accesskey)
	resp, e := s.client.GetUserPermission(ctx, req, header)
	if e != nil {
		return e
	}
	if resp.GetCanread() || resp.GetAdmin() {
		return nil
	}
	return ecode.ErrPermission
}

// if pass will return nil
// if not pass will reutrn ecode.ErrPermission(use cerror.Equal to check)(https://github.com/chenjie199234/Corelib/tree/main/cerror)
func (s *PermissionSdk) CheckWrite(ctx context.Context, userid string, nodeid []uint32) error {
	req := &api.GetUserPermissionReq{
		UserId: userid,
		NodeId: nodeid,
	}
	header := make(http.Header)
	header.Set("Access-Key", s.accesskey)
	resp, e := s.client.GetUserPermission(ctx, req, header)
	if e != nil {
		return e
	}
	if resp.GetCanread() && resp.GetCanwrite() {
		return nil
	}
	return ecode.ErrPermission
}
