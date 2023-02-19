package permission

import (
	"context"
	"github.com/chenjie199234/admin/api"

	"github.com/chenjie199234/Corelib/util/egroup"
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

func (s *Sdk) CheckMulti(ctx context.Context, userid string, readNodeIDs [][]uint32, writeNodeIDs [][]uint32, adminNodeIDs [][]uint32) (bool, error) {
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
