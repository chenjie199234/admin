package permission

import (
	"context"
	"crypto/tls"
	"errors"
	"os"
	"strconv"

	"github.com/chenjie199234/admin/api"

	"github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/util/egroup"
	"github.com/chenjie199234/Corelib/web"
)

type PermissionSdk struct {
	client api.PermissionWebClient
}

var (
	ErrMissingEnvProject = errors.New("missing env PERMISSION_SERVICE_PROJECT")
	ErrMissingEnvGroup   = errors.New("missing env PERMISSION_SERVICE_GROUP")
	ErrMissingEnvHost    = errors.New("missing env PERMISSION_SERVICE_WEB_HOST")
	ErrWrongEnvPort      = errors.New("env PERMISSION_SERVICE_WEB_PORT must be number <= 65535")
)

// if tlsc is not nil,the tls will be actived
// must set below env:
// PERMISSION_SERVICE_PROJECT
// PERMISSION_SERVICE_GROUP
// PERMISSION_SERVICE_WEB_HOST
// PERMISSION_SERVICE_WEB_PORT
func NewPermissionSdk(selfproject, selfgroup, selfapp string, tlsc *tls.Config) (*PermissionSdk, error) {
	project, group, host, port, e := env()
	if e != nil {
		return nil, e
	}
	di, e := discover.NewStaticDiscover(project, group, "admin", []string{host}, 0, 0, port)
	if e != nil {
		return nil, e
	}
	tmpclient, e := web.NewWebClient(nil, di, selfproject, selfgroup, selfapp, project, group, "admin", tlsc)
	if e != nil {
		return nil, e
	}
	return &PermissionSdk{client: api.NewPermissionWebClient(tmpclient)}, nil
}
func env() (projectname, group string, host string, port int, e error) {
	if str, ok := os.LookupEnv("PERMISSION_SERVICE_PROJECT"); ok && str != "<PERMISSION_SERVICE_PROJECT>" && str != "" {
		projectname = str
	} else {
		return "", "", "", 0, ErrMissingEnvProject
	}
	if str, ok := os.LookupEnv("PERMISSION_SERVICE_GROUP"); ok && str != "<PERMISSION_SERVICE_GROUP>" && str != "" {
		group = str
	} else {
		return "", "", "", 0, ErrMissingEnvGroup
	}
	if str, ok := os.LookupEnv("PERMISSION_SERVICE_WEB_HOST"); ok && str != "<PERMISSION_SERVICE_WEB_HOST>" && str != "" {
		host = str
	} else {
		return "", "", "", 0, ErrMissingEnvHost
	}
	if str, ok := os.LookupEnv("PERMISSION_SERVICE_WEB_PORT"); ok && str != "<PERMISSION_SERVICE_WEB_PORT>" && str != "" {
		var e error
		port, e = strconv.Atoi(str)
		if e != nil || port < 0 || port > 65535 {
			return "", "", "", 0, ErrWrongEnvPort
		}
	}
	return
}

func (s *PermissionSdk) CheckMulti(ctx context.Context, userid string, readNodeIDs [][]uint32, writeNodeIDs [][]uint32, adminNodeIDs [][]uint32) (bool, error) {
	pass := true
	eg := egroup.GetGroup(ctx)
	for _, v := range readNodeIDs {
		nodeid := v
		eg.Go(func(gctx context.Context) error {
			singlepass, e := s.CheckRead(gctx, userid, nodeid)
			if e != nil {
				return e
			}
			if !singlepass {
				pass = false
			}
			return nil
		})
	}
	for _, v := range writeNodeIDs {
		nodeid := v
		eg.Go(func(gctx context.Context) error {
			singlepass, e := s.CheckWrite(gctx, userid, nodeid)
			if e != nil {
				return e
			}
			if !singlepass {
				pass = false
			}
			return nil
		})
	}
	for _, v := range adminNodeIDs {
		nodeid := v
		eg.Go(func(gctx context.Context) error {
			singlepass, e := s.CheckAdmin(gctx, userid, nodeid)
			if e != nil {
				return e
			}
			if !singlepass {
				pass = false
			}
			return nil
		})
	}
	e := egroup.PutGroup(eg)
	return pass, e
}

func (s *PermissionSdk) CheckAdmin(ctx context.Context, userid string, nodeid []uint32) (bool, error) {
	req := &api.GetUserPermissionReq{
		UserId: userid,
		NodeId: nodeid,
	}
	resp, e := s.client.GetUserPermission(ctx, req, nil)
	if e != nil {
		return false, e
	}
	return resp.GetAdmin(), nil
}

func (s *PermissionSdk) CheckRead(ctx context.Context, userid string, nodeid []uint32) (bool, error) {
	req := &api.GetUserPermissionReq{
		UserId: userid,
		NodeId: nodeid,
	}
	resp, e := s.client.GetUserPermission(ctx, req, nil)
	if e != nil {
		return false, e
	}
	return resp.GetCanread() || resp.GetAdmin(), nil
}

func (s *PermissionSdk) CheckWrite(ctx context.Context, userid string, nodeid []uint32) (bool, error) {
	req := &api.GetUserPermissionReq{
		UserId: userid,
		NodeId: nodeid,
	}
	resp, e := s.client.GetUserPermission(ctx, req, nil)
	if e != nil {
		return false, e
	}
	return resp.GetCanread() && resp.GetCanwrite(), nil
}
