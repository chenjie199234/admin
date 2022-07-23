package permission

import (
	"context"
	"sort"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	cerror "github.com/chenjie199234/Corelib/error"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/util/egroup"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
)

//Service subservice for permission business
type Service struct {
	userDao       *userdao.Dao
	permissionDao *permissiondao.Dao
}

//Start -
func Start() *Service {
	return &Service{
		userDao:       userdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}
}
func (s *Service) GetUserPermission(ctx context.Context, req *api.GetUserPermissionReq) (*api.GetUserPermissionResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	obj, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[GetUserPermission] userid:", req.UserId, "format error:", e)
		return nil, ecode.ErrReq
	}
	canread, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, obj, req.NodeId)
	if e != nil {
		log.Error(ctx, "[GetUserPermission] userid:", req.UserId, "nodeid:", req.NodeId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.GetUserPermissionResp{Canread: canread, Canwrite: canwrite, Admin: admin}, nil
}
func (s *Service) UpdateUserPermission(ctx context.Context, req *api.UpdateUserPermissionReq) (*api.UpdateUserPermissionResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	if !req.Admin && req.Canwrite && !req.Canread {
		return nil, ecode.ErrReq
	}
	targetobj, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[UpdateUserPermission] target userid:", req.UserId, "format error:", e)
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[UpdateUserPermission] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.permissionDao.MongoUpdateUserPermission(ctx, obj, targetobj, req.NodeId, req.Admin, req.Canread, req.Canwrite); e != nil {
		log.Error(ctx, "[UpdateUserPermission] userid:", userid, "target userid:", req.UserId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.UpdateUserPermissionResp{}, nil
}
func (s *Service) AddNode(ctx context.Context, req *api.AddNodeReq) (*api.AddNodeResp, error) {
	if req.PnodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[AddNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.permissionDao.MongoAddNode(ctx, obj, req.PnodeId, req.NodeName, req.NodeData); e != nil {
		log.Error(ctx, "[AddNode] userid:", userid, "name:", req.NodeName, "data:", req.NodeData, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.AddNodeResp{}, nil
}
func (s *Service) UpdateNode(ctx context.Context, req *api.UpdateNodeReq) (*api.UpdateNodeResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[UpdateNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.permissionDao.MongoUpdateNode(ctx, obj, req.NodeId, req.NodeName, req.NodeData); e != nil {
		log.Error(ctx, "[UpdateNode] userid:", userid, "nodeid:", req.NodeId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.UpdateNodeResp{}, nil
}
func (s *Service) MoveNode(ctx context.Context, req *api.MoveNodeReq) (*api.MoveNodeResp, error) {
	if req.NodeId[0] != 0 || req.PnodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[MoveNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	if e := s.permissionDao.MongoMoveNode(ctx, obj, req.NodeId, req.PnodeId); e != nil {
		log.Error(ctx, "[MoveNode] userid:", userid, "old nodeid:", req.NodeId, "new parent:", req.PnodeId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.MoveNodeResp{}, nil
}
func (s *Service) DelNode(ctx context.Context, req *api.DelNodeReq) (*api.DelNodeResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[DelNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.permissionDao.MongoDeleteNode(ctx, obj, req.NodeId); e != nil {
		log.Error(ctx, "[DelNode] userid:", userid, "nodeid:", req.NodeId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.DelNodeResp{}, nil
}
func (s *Service) ListUserNode(ctx context.Context, req *api.ListUserNodeReq) (*api.ListUserNodeResp, error) {
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[ListUserNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	usernodes, e := s.permissionDao.MongoGetUserNodes(ctx, obj, nil)
	if e != nil {
		log.Error(ctx, "[ListUserNode] userid:", userid, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	sort.Slice(usernodes, func(i, j int) bool {
		return len(usernodes[i].NodeId) < len(usernodes[j].NodeId)
	})
	adminnodeids := make([][]uint32, 0, len(usernodes))
	nodeids := make([][]uint32, 0, len(usernodes))
	for _, usernode := range usernodes {
		//该node是否在某个admin node下面,如果是在下面，那么直接跳过
		jump := false
		for _, adminnodeid := range adminnodeids {
			belowadmin := true
			for i := range adminnodeid {
				if adminnodeid[i] != usernode.NodeId[i] {
					belowadmin = false
					break
				}
			}
			if belowadmin {
				jump = true
				break
			}
		}
		if jump {
			continue
		}
		if usernode.X {
			adminnodeids = append(adminnodeids, usernode.NodeId)
		}
		nodeids = append(nodeids, usernode.NodeId)
	}
	usernodeinfos, e := s.permissionDao.MongoGetNodes(ctx, nodeids)
	if e != nil {
		log.Error(ctx, "[ListUserNode] userid:", userid, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	treenodes := make([]*api.NodeInfo, 0, len(nodeids))
	for _, nodeid := range nodeids {
		var usernode *model.UserNode
		for _, v := range usernodes {
			if len(nodeid) != len(v.NodeId) {
				continue
			}
			same := true
			for i := range nodeid {
				if nodeid[i] != v.NodeId[i] {
					same = false
					break
				}
			}
			if same {
				usernode = v
				break
			}
		}
		var nodeinfo *model.Node
		for _, v := range usernodeinfos {
			if len(nodeid) != len(v.NodeId) {
				continue
			}
			same := true
			for i := range nodeid {
				if nodeid[i] != v.NodeId[i] {
					same = false
					break
				}
			}
			if same {
				nodeinfo = v
				break
			}
		}
		if nodeinfo == nil {
			continue
		}
		tmp := &api.NodeInfo{
			NodeId:   nodeid,
			NodeName: nodeinfo.NodeName,
			NodeData: nodeinfo.NodeData,
		}
		if usernode.X {
			tmp.Canread = 1
			tmp.Canwrite = 1
			tmp.Admin = 1
		} else if usernode.R {
			tmp.Canread = 1
			if usernode.W {
				tmp.Canwrite = 1
			}
		}
		treenodes = append(treenodes, tmp)
	}
	g := egroup.GetGroup(ctx)
	for _, v := range treenodes {
		if v.Admin != 1 {
			continue
		}
		treenode := v
		g.Go(func(gctx context.Context) error {
			nodes, e := s.permissionDao.MongoListNode(ctx, treenode.NodeId)
			if e != nil {
				log.Error(gctx, "[ListUserNode] userid:", userid, "nodeid:", treenode.NodeId, e)
				return e
			}
			sort.Slice(nodes, func(i, j int) bool {
				return len(nodes[i].NodeId) < len(nodes[j].NodeId)
			})
			for _, node := range nodes {
				addTreeNode(treenode, &api.NodeInfo{
					NodeId:   node.NodeId,
					NodeName: node.NodeName,
					NodeData: node.NodeData,
					Canread:  1,
					Canwrite: 1,
					Admin:    1,
					Children: make([]*api.NodeInfo, 0),
				})
			}
			return nil
		})
	}
	if e := egroup.PutGroup(g); e != nil {
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	root := &api.NodeInfo{
		NodeId: []uint32{},
	}
	for _, v := range treenodes {
		addTreeNode(root, v)
	}
	return &api.ListUserNodeResp{
		Nodes: root.Children,
	}, nil
}
func (s *Service) ListAllNode(ctx context.Context, req *api.ListAllNodeReq) (*api.ListAllNodeResp, error) {
	nodes, e := s.permissionDao.MongoListNode(ctx, nil)
	if e != nil {
		log.Error(ctx, "[ListAllNode]", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	root := &api.NodeInfo{}
	//find the root
	for _, node := range nodes {
		if len(node.NodeId) == 1 {
			root.NodeId = node.NodeId
			root.NodeName = node.NodeName
			break
		}
	}
	if len(root.NodeId) == 0 {
		return nil, ecode.ErrNotInited
	}
	for _, node := range nodes {
		if len(node.NodeId) > 1 {
			addTreeNode(root, &api.NodeInfo{
				NodeId:   node.NodeId,
				NodeName: node.NodeName,
				Children: make([]*api.NodeInfo, 0),
			})
		}
	}
	return &api.ListAllNodeResp{Nodes: root}, nil
}
func addTreeNode(root, node *api.NodeInfo) bool {
	if len(root.NodeId) > len(node.NodeId) {
		return false
	}
	isprefix := true
	for i := range root.NodeId {
		if root.NodeId[i] != node.NodeId[i] {
			isprefix = false
			break
		}
	}
	if !isprefix {
		return false
	}
	if len(root.NodeId) == len(node.NodeId) {
		return true
	}
	for _, child := range root.Children {
		if addTreeNode(child, node) {
			return true
		}
	}
	root.Children = append(root.Children, node)
	return true
}
func (s *Service) ListNodeUser(ctx context.Context, req *api.ListNodeUserReq) (*api.ListNodeUserResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[ListNodeUser] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	_, _, x, e := s.permissionDao.MongoGetUserPermission(ctx, obj, req.NodeId)
	if e != nil {
		log.Error(ctx, "[ListNodeUser] userid:", userid, "nodeid:", req.NodeId, "get user permission error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	if !x {
		log.Error(ctx, "[ListNodeUser] userid:", userid, "nodeid:", req.NodeId, "missing permission")
		return nil, ecode.ErrPermission
	}
	nodeusers, e := s.permissionDao.MongoGetNodeUsers(ctx, req.NodeId, nil)
	if e != nil {
		log.Error(ctx, "[ListNodeUser] userid:", userid, "nodeid:", req.NodeId, "get node users error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	undup := make(map[primitive.ObjectID]*struct{}, 100)
	for _, r := range nodeusers.R {
		undup[r] = nil
	}
	for _, w := range nodeusers.W {
		undup[w] = nil
	}
	for _, x := range nodeusers.X {
		undup[x] = nil
	}
	userids := make([]primitive.ObjectID, 0, len(undup))
	for userid := range undup {
		userids = append(userids, userid)
	}
	users, e := s.userDao.MongoGetUsers(ctx, userids)
	if e != nil {
		log.Error(ctx, "[ListNodeUser] userid:", userid, "nodeid:", req.NodeId, "get all node users' info error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	resp := &api.ListNodeUserResp{
		Canread:  make(map[string]*api.UserInfo, len(nodeusers.R)),
		Canwrite: make(map[string]*api.UserInfo, len(nodeusers.W)),
		Admin:    make(map[string]*api.UserInfo, len(nodeusers.X)),
	}
	for _, r := range nodeusers.R {
		if userinfo, ok := users[r]; ok {
			resp.Canread[userinfo.ID.Hex()] = &api.UserInfo{
				UserId:     userinfo.ID.Hex(),
				UserName:   userinfo.Name,
				Department: userinfo.Department,
				Ctime:      userinfo.Ctime,
			}
		}
	}
	for _, w := range nodeusers.W {
		if userinfo, ok := users[w]; ok {
			resp.Canwrite[userinfo.ID.Hex()] = &api.UserInfo{
				UserId:     userinfo.ID.Hex(),
				UserName:   userinfo.Name,
				Department: userinfo.Department,
				Ctime:      userinfo.Ctime,
			}
		}
	}
	for _, x := range nodeusers.X {
		if userinfo, ok := users[x]; ok {
			resp.Admin[userinfo.ID.Hex()] = &api.UserInfo{
				UserId:     userinfo.ID.Hex(),
				UserName:   userinfo.Name,
				Department: userinfo.Department,
				Ctime:      userinfo.Ctime,
			}
		}
	}
	return resp, nil
}

func (s *Service) ListAdmin(ctx context.Context, req *api.ListAdminReq) (*api.ListAdminResp, error) {
	users, e := s.permissionDao.MongoListAdmin(ctx, req.NodeId)
	if e != nil {
		log.Error(ctx, "[ListAdmin] node_id:", req.NodeId, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	resp := &api.ListAdminResp{
		Users: make([]*api.UserInfo, 0, len(users)),
	}
	for _, user := range users {
		resp.Users = append(resp.Users, &api.UserInfo{
			UserId:     user.ID.Hex(),
			UserName:   user.Name,
			Department: user.Department,
		})
	}
	return resp, nil
}

//Stop -
func (s *Service) Stop() {

}
