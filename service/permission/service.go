package permission

import (
	"context"
	"encoding/json"
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
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/egroup"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
)

// Service subservice for permission business
type Service struct {
	userDao       *userdao.Dao
	permissionDao *permissiondao.Dao
}

// Start -
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
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[GetUserPermission] userid:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	canread, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, target, req.NodeId, true)
	if e != nil {
		log.Error(ctx, "[GetUserPermission] userid:", req.UserId, "nodeid:", req.NodeId, e)
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
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateUserPermission] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[UpdateUserPermission] target:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	//permission check
	_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.UserControlNodeId, true)
	if e != nil {
		log.Error(ctx, "[UpdateUserPermission] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !admin {
		return nil, ecode.ErrPermission
	}
	//logic
	if e = s.permissionDao.MongoUpdateUserPermission(ctx, operator, target, req.NodeId, req.Admin, req.Canread, req.Canwrite); e != nil {
		log.Error(ctx, "[UpdateUserPermission] operator:", md["Token-Data"], "target:", req.UserId, "nodeid:", req.NodeId, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.UpdateUserPermissionResp{}, nil
}
func (s *Service) UpdateRolePermission(ctx context.Context, req *api.UpdateRolePermissionReq) (*api.UpdateRolePermissionResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	if !req.Admin && req.Canwrite && !req.Canread {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateRolePermission] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	//permission check
	_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.RoleControlNodeId, true)
	if e != nil {
		log.Error(ctx, "[UpdateRolePermission] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !admin {
		return nil, ecode.ErrPermission
	}
	//logic
	if e = s.permissionDao.MongoUpdateRolePermission(ctx, operator, req.RoleName, req.NodeId, req.Admin, req.Canread, req.Canwrite); e != nil {
		log.Error(ctx, "[UpdateRolePermission] operator:", md["Token-Data"], "rolename:", req.RoleName, "nodeid:", req.NodeId, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.UpdateRolePermissionResp{}, nil
}
func (s *Service) AddNode(ctx context.Context, req *api.AddNodeReq) (*api.AddNodeResp, error) {
	if req.PnodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[AddNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.permissionDao.MongoAddNode(ctx, operator, req.PnodeId, req.NodeName, req.NodeData); e != nil {
		log.Error(ctx, "[AddNode] operator:", md["Token-Data"], "name:", req.NodeName, "data:", req.NodeData, "parent nodeid:", req.PnodeId, e)
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
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.permissionDao.MongoUpdateNode(ctx, operator, req.NodeId, req.NodeName, req.NodeData); e != nil {
		log.Error(ctx, "[UpdateNode] operator:", md["Token-Data"], "nodeid:", req.NodeId, e)
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
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[MoveNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	if e := s.permissionDao.MongoMoveNode(ctx, operator, req.NodeId, req.PnodeId); e != nil {
		log.Error(ctx, "[MoveNode] operator:", md["Token-Data"], "nodeid:", req.NodeId, "new parent nodeid:", req.PnodeId, e)
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
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	if e = s.permissionDao.MongoDeleteNode(ctx, operator, req.NodeId); e != nil {
		log.Error(ctx, "[DelNode] operator:", md["Token-Data"], "nodeid:", req.NodeId, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	return &api.DelNodeResp{}, nil
}
func (s *Service) ListUserNode(ctx context.Context, req *api.ListUserNodeReq) (*api.ListUserNodeResp, error) {
	if req.UserId == "" {
		//list self's
		md := metadata.GetMetadata(ctx)
		req.UserId = md["Token-Data"]
	} else {
		//list other user's
		//need to check the operator's permission
		md := metadata.GetMetadata(ctx)
		operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
		if e != nil {
			log.Error(ctx, "[ListUserNode] operator:", md["Token-Data"], "format wrong:", e)
			return nil, ecode.ErrAuth
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.UserControlNodeId, true)
		if e != nil {
			log.Error(ctx, "[ListUserNode] operator:", md["Token-Data"], "get permission failed:", e)
			return nil, ecode.ErrSystem
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[ListUserNode] target:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	undup := make(map[string]*api.NodeInfo)
	usernodes, e := s.permissionDao.MongoGetUserNodes(ctx, target, nil)
	if e != nil {
		log.Error(ctx, "[ListUserNode] target:", req.UserId, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	for _, usernode := range usernodes {
		nodeidstr, _ := json.Marshal(usernode.NodeId)
		if exist, ok := undup[common.Byte2str(nodeidstr)]; ok {
			if usernode.R {
				exist.Canread = usernode.R
			}
			if usernode.W {
				exist.Canwrite = usernode.W
			}
			if usernode.X {
				exist.Admin = usernode.X
			}
		} else {
			undup[common.Byte2str(nodeidstr)] = &api.NodeInfo{
				NodeId:   usernode.NodeId,
				NodeName: "",
				NodeData: "",
				Canread:  usernode.R,
				Canwrite: usernode.W,
				Admin:    usernode.X,
				Children: make([]*api.NodeInfo, 0, 10),
			}
		}
	}
	if req.NeedUserRoleNode {
		userrolenodes, e := s.permissionDao.MongoGetUserRoleNodes(ctx, target, nil)
		if e != nil {
			log.Error(ctx, "[ListUserNode] target:", req.UserId, e)
			if _, ok := e.(*cerror.Error); ok {
				return nil, e
			}
			return nil, ecode.ErrSystem
		}
		for _, v := range userrolenodes {
			for _, userrolenode := range v {
				nodeidstr, _ := json.Marshal(userrolenode.NodeId)
				if exist, ok := undup[common.Byte2str(nodeidstr)]; ok {
					if userrolenode.R {
						exist.Canread = userrolenode.R
					}
					if userrolenode.W {
						exist.Canwrite = userrolenode.W
					}
					if userrolenode.X {
						exist.Admin = userrolenode.X
					}
				} else {
					undup[common.Byte2str(nodeidstr)] = &api.NodeInfo{
						NodeId:   userrolenode.NodeId,
						NodeName: "",
						NodeData: "",
						Canread:  userrolenode.R,
						Canwrite: userrolenode.W,
						Admin:    userrolenode.X,
						Children: make([]*api.NodeInfo, 0, 10),
					}
				}
			}
		}
	}
	nodes := make([]*api.NodeInfo, 0, len(undup))
	for _, v := range undup {
		nodes = append(nodes, v)
	}
	sort.Slice(nodes, func(i, j int) bool {
		return len(nodes[i].NodeId) < len(nodes[j].NodeId)
	})
	adminnodeids := make([][]uint32, 0, len(nodes))
	nodeids := make([][]uint32, 0, len(nodes))
	for _, node := range nodes {
		//if this node is below another admin node,this node can be jumped
		jump := false
		for _, adminnodeid := range adminnodeids {
			isprefix := true
			for i := range adminnodeid {
				if adminnodeid[i] != node.NodeId[i] {
					isprefix = false
					break
				}
			}
			if isprefix {
				jump = true
				break
			}
		}
		if jump {
			continue
		}
		if node.Admin {
			adminnodeids = append(adminnodeids, node.NodeId)
		}
		nodeids = append(nodeids, node.NodeId)
	}
	nodeinfos, e := s.permissionDao.MongoGetNodes(ctx, nodeids)
	if e != nil {
		log.Error(ctx, "[ListUserNode] target:", req.UserId, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	for _, nodeinfo := range nodeinfos {
		nodeidstr, _ := json.Marshal(nodeinfo.NodeId)
		exist, _ := undup[common.Byte2str(nodeidstr)]
		exist.NodeName = nodeinfo.NodeName
		exist.NodeData = nodeinfo.NodeData
	}
	nodes = nodes[:0]
	for _, nodeid := range nodeids {
		nodeidstr, _ := json.Marshal(nodeid)
		nodes = append(nodes, undup[common.Byte2str(nodeidstr)])
	}
	g := egroup.GetGroup(ctx)
	for _, v := range nodes {
		if !v.Admin {
			continue
		}
		node := v
		g.Go(func(gctx context.Context) error {
			belowadminnodes, e := s.permissionDao.MongoListNode(ctx, node.NodeId)
			if e != nil {
				log.Error(gctx, "[ListUserNode] target:", req.UserId, "admin nodeid:", node.NodeId, e)
				return e
			}
			sort.Slice(belowadminnodes, func(i, j int) bool {
				return len(belowadminnodes[i].NodeId) < len(belowadminnodes[j].NodeId)
			})
			for _, belowadminnode := range belowadminnodes {
				addTreeNode(node, &api.NodeInfo{
					NodeId:   belowadminnode.NodeId,
					NodeName: belowadminnode.NodeName,
					NodeData: belowadminnode.NodeData,
					Canread:  true,
					Canwrite: true,
					Admin:    true,
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
	for _, node := range nodes {
		addTreeNode(root, node)
	}
	return &api.ListUserNodeResp{Nodes: root.Children}, nil
}
func (s *Service) ListRoleNode(ctx context.Context, req *api.ListRoleNodeReq) (*api.ListRoleNodeResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[ListRoleNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrAuth
	}
	canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, model.RoleControlNodeId, true)
	if e != nil {
		log.Error(ctx, "[ListRoleNode] operator:", md["Token-Data"], "get permission failed:", e)
		return nil, ecode.ErrSystem
	}
	if !canread && !admin {
		return nil, ecode.ErrPermission
	}
	rolenodes, e := s.permissionDao.MongoGetRoleNodes(ctx, req.RoleName, nil)
	if e != nil {
		log.Error(ctx, "[ListRoleNode] operator:", md["Token-Data"], "rolename:", req.RoleName, e)
		return nil, ecode.ErrSystem
	}
	nodes := make([]*api.NodeInfo, 0, len(rolenodes))
	undup := make(map[string]*api.NodeInfo, len(rolenodes))
	for _, rolenode := range rolenodes {
		nodeidstr, _ := json.Marshal(rolenode.NodeId)
		if exist, ok := undup[common.Byte2str(nodeidstr)]; ok {
			if rolenode.R {
				exist.Canread = rolenode.R
			}
			if rolenode.W {
				exist.Canwrite = rolenode.W
			}
			if rolenode.X {
				exist.Admin = rolenode.X
			}
		} else {
			tmp := &api.NodeInfo{
				NodeId:   rolenode.NodeId,
				NodeName: "",
				NodeData: "",
				Canread:  rolenode.R,
				Canwrite: rolenode.W,
				Admin:    rolenode.X,
			}
			undup[common.Byte2str(nodeidstr)] = tmp
			nodes = append(nodes, tmp)
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return len(nodes[i].NodeId) < len(nodes[j].NodeId)
	})
	adminnodeids := make([][]uint32, 0, len(nodes))
	nodeids := make([][]uint32, 0, len(nodes))
	for _, node := range nodes {
		//if this node is below another admin node,this node can be jumped
		jump := false
		for _, adminnodeid := range adminnodeids {
			isprefix := true
			for i := range adminnodeid {
				if adminnodeid[i] != node.NodeId[i] {
					isprefix = false
					break
				}
			}
			if isprefix {
				jump = true
				break
			}
		}
		if jump {
			continue
		}
		if node.Admin {
			adminnodeids = append(adminnodeids, node.NodeId)
		}
		nodeids = append(nodeids, node.NodeId)
	}
	nodeinfos, e := s.permissionDao.MongoGetNodes(ctx, nodeids)
	if e != nil {
		log.Error(ctx, "[ListRoleNode] operator:", md["Token-Data"], "rolename:", req.RoleName, e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	for _, nodeinfo := range nodeinfos {
		nodeidstr, _ := json.Marshal(nodeinfo.NodeId)
		exist, _ := undup[common.Byte2str(nodeidstr)]
		exist.NodeName = nodeinfo.NodeName
		exist.NodeData = nodeinfo.NodeData
	}
	nodes = nodes[:0]
	for _, nodeid := range nodeids {
		nodeidstr, _ := json.Marshal(nodeid)
		nodes = append(nodes, undup[common.Byte2str(nodeidstr)])
	}
	g := egroup.GetGroup(ctx)
	for _, v := range nodes {
		if !v.Admin {
			continue
		}
		node := v
		g.Go(func(gctx context.Context) error {
			belowadminnodes, e := s.permissionDao.MongoListNode(ctx, node.NodeId)
			if e != nil {
				log.Error(gctx, "[ListRoleNode] operator:", md["Token-Data"], "rolename:", req.RoleName, "admin nodeid:", node.NodeId, e)
				return e
			}
			sort.Slice(belowadminnodes, func(i, j int) bool {
				return len(belowadminnodes[i].NodeId) < len(belowadminnodes[j].NodeId)
			})
			for _, belowadminnode := range belowadminnodes {
				addTreeNode(node, &api.NodeInfo{
					NodeId:   belowadminnode.NodeId,
					NodeName: belowadminnode.NodeName,
					NodeData: belowadminnode.NodeData,
					Canread:  true,
					Canwrite: true,
					Admin:    true,
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
	for _, node := range nodes {
		addTreeNode(root, node)
	}
	return &api.ListRoleNodeResp{Nodes: root.Children}, nil
}
func (s *Service) ListAllNode(ctx context.Context, req *api.ListAllNodeReq) (*api.ListAllNodeResp, error) {
	md := metadata.GetMetadata(ctx)
	nodes, e := s.permissionDao.MongoListNode(ctx, nil)
	if e != nil {
		log.Error(ctx, "[ListAllNode] operator:", md["Token-Data"], e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	root := &api.NodeInfo{
		Children: make([]*api.NodeInfo, 0, 10),
	}
	//find the root
	for _, node := range nodes {
		if len(node.NodeId) == 1 {
			root.NodeId = node.NodeId
			root.NodeName = node.NodeName
			root.NodeData = node.NodeData
			break
		}
	}
	if len(root.NodeId) == 0 {
		return nil, ecode.ErrNotInited
	}
	sort.Slice(nodes, func(i, j int) bool {
		return len(nodes[i].NodeId) < len(nodes[j].NodeId)
	})
	for _, node := range nodes {
		if len(node.NodeId) <= 1 {
			continue
		}
		addTreeNode(root, &api.NodeInfo{
			NodeId:   node.NodeId,
			NodeName: node.NodeName,
			NodeData: node.NodeData,
			Children: make([]*api.NodeInfo, 0),
		})
	}
	return &api.ListAllNodeResp{Nodes: root.Children}, nil
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

// Stop -
func (s *Service) Stop() {

}
