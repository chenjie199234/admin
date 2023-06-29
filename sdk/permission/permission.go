package permission

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/chenjie199234/admin/api"

	"github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/util/egroup"
	"github.com/chenjie199234/Corelib/web"
)

type Sdk struct {
	client api.PermissionWebClient
}

// if tlsc is not nil,the tls will be actived
func NewPermissionSdk(selfappgroup, selfappname, serverappgroup, serverhost string, tlsc *tls.Config) (*Sdk, error) {
	di := discover.NewDirectDiscover(serverappgroup, "admin", serverhost, 9000, 10000, 8000)
	tmpclient, e := web.NewWebClient(&web.ClientConfig{
		ConnectTimeout: time.Second * 3,
		GlobalTimeout:  0,
		HeartProbe:     time.Second * 3,
		IdleTimeout:    time.Second * 10,
	}, di, selfappgroup, selfappname, serverappgroup, "admin", tlsc)
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
