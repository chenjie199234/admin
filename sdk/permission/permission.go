package permission

import (
	"github.com/chenjie199234/admin/api"
	"golang.org/x/net/context"

	"github.com/chenjie199234/Corelib/web"
)

type Sdk struct {
	client api.PermissionWebClient
}

func NewPermissionSdk(selfgroup, selfname, servergroup, serverhost string) (*Sdk, error) {
	tmpclient, e := web.NewWebClient(&web.ClientConfig{}, selfgroup, selfname, servergroup, "admin", serverhost)
	if e != nil {
		return nil, e
	}
	return &Sdk{client: api.NewPermissionWebClient(tmpclient)}, nil
}

func (s *Sdk) CheckAdmin(ctx context.Context, userid string, nodeid []uint32) (bool, error) {
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

func (s *Sdk) CheckRead(ctx context.Context, userid string, nodeid []uint32) (bool, error) {
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

func (s *Sdk) CheckWrite(ctx context.Context, userid string, nodeid []uint32) (bool, error) {
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
