package permission

import (
	"context"
	"sort"
	"strings"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/pool"
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
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.NodeId {
		buf.AppendUint32(v)
		if i != len(req.NodeId)-1 {
			buf.AppendByte(',')
		}
	}
	nodeid := buf.String()
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[GetUserPermission] target:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	canread, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, target, nodeid, true)
	if e != nil {
		log.Error(ctx, "[GetUserPermission] target:", req.UserId, "nodeid:", nodeid, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
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
	buf1 := pool.GetBuffer()
	defer pool.PutBuffer(buf1)
	for i, v := range req.NodeId {
		buf1.AppendUint32(v)
		if i != len(req.NodeId)-1 {
			buf1.AppendByte(',')
		}
	}
	nodeid := buf1.String()
	buf2 := pool.GetBuffer()
	defer pool.PutBuffer(buf2)
	buf2.AppendUint32(req.NodeId[0])
	buf2.AppendByte(',')
	buf2.AppendUint32(req.NodeId[1])
	project := buf2.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateUserPermission] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	target, e := primitive.ObjectIDFromHex(req.UserId)
	if e != nil {
		log.Error(ctx, "[UpdateUserPermission] target:", req.UserId, "format wrong:", e)
		return nil, ecode.ErrReq
	}
	if !operator.IsZero() {
		//user control permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.UserControl, true)
		if e != nil {
			log.Error(ctx, "[UpdateUserPermission] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}
	//logic
	if e = s.permissionDao.MongoUpdateUserPermission(ctx, operator, target, nodeid, req.Admin, req.Canread, req.Canwrite); e != nil {
		log.Error(ctx, "[UpdateUserPermission] operator:", md["Token-Data"], "target:", req.UserId, "nodeid:", nodeid, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.UpdateUserPermissionResp{}, nil
}
func (s *Service) UpdateRolePermission(ctx context.Context, req *api.UpdateRolePermissionReq) (*api.UpdateRolePermissionResp, error) {
	if req.NodeId[0] != 0 || req.Project[0] != 0 || (req.Project[1] != req.NodeId[1]) {
		return nil, ecode.ErrReq
	}
	if !req.Admin && req.Canwrite && !req.Canread {
		return nil, ecode.ErrReq
	}
	buf1 := pool.GetBuffer()
	defer pool.PutBuffer(buf1)
	for i, v := range req.NodeId {
		buf1.AppendUint32(v)
		if i != len(req.NodeId)-1 {
			buf1.AppendByte(',')
		}
	}
	nodeid := buf1.String()
	buf2 := pool.GetBuffer()
	defer pool.PutBuffer(buf2)
	for i, v := range req.Project {
		buf2.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf2.AppendByte(',')
		}
	}
	project := buf2.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateRolePermission] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//role control permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.RoleControl, true)
		if e != nil {
			log.Error(ctx, "[UpdateRolePermission] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}
	//logic
	if e = s.permissionDao.MongoUpdateRolePermission(ctx, operator, project, req.RoleName, nodeid, req.Admin, req.Canread, req.Canwrite); e != nil {
		log.Error(ctx, "[UpdateRolePermission] operator:", md["Token-Data"], "project:", project, "rolename:", req.RoleName, "nodeid:", nodeid, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.UpdateRolePermissionResp{}, nil
}
func (s *Service) AddNode(ctx context.Context, req *api.AddNodeReq) (*api.AddNodeResp, error) {
	if req.PnodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	if req.PnodeId[1] == 1 || req.PnodeId[1] == 2 || req.PnodeId[1] == 3 {
		//0,1 -> UserControl
		//0,2 -> RoleControl
		//0,3 -> ConfigControl
		//these are default,can't modify
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.PnodeId {
		buf.AppendUint32(v)
		if i != len(req.PnodeId)-1 {
			buf.AppendByte(',')
		}
	}
	pnodeid := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[AddNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if e = s.permissionDao.MongoAddNode(ctx, operator, pnodeid, req.NodeName, req.NodeData); e != nil {
		log.Error(ctx, "[AddNode] operator:", md["Token-Data"], "nodename:", req.NodeName, "nodedata:", req.NodeData, "parent nodeid:", pnodeid, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.AddNodeResp{}, nil
}
func (s *Service) UpdateNode(ctx context.Context, req *api.UpdateNodeReq) (*api.UpdateNodeResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	if req.NodeId[1] == 1 || req.NodeId[1] == 2 || req.NodeId[1] == 3 {
		//0,1 -> UserControl
		//0,2 -> RoleControl
		//0,3 -> ConfigControl
		//these are default,can't modify
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.NodeId {
		buf.AppendUint32(v)
		if i != len(req.NodeId)-1 {
			buf.AppendByte(',')
		}
	}
	nodeid := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[UpdateNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if e = s.permissionDao.MongoUpdateNode(ctx, operator, nodeid, req.NodeName, req.NodeData); e != nil {
		log.Error(ctx, "[UpdateNode] operator:", md["Token-Data"], "nodeid:", nodeid, "nodename:", req.NodeName, "nodedata:", req.NodeData, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.UpdateNodeResp{}, nil
}
func (s *Service) MoveNode(ctx context.Context, req *api.MoveNodeReq) (*api.MoveNodeResp, error) {
	if req.NodeId[0] != 0 || req.PnodeId[0] != 0 || (req.NodeId[1] != req.PnodeId[1]) {
		return nil, ecode.ErrReq
	}
	if req.NodeId[1] == 1 || req.NodeId[1] == 2 || req.NodeId[1] == 3 {
		//0,1 -> UserControl
		//0,2 -> RoleControl
		//0,3 -> ConfigControl
		//these are default,can't modify
		return nil, ecode.ErrReq
	}
	buf1 := pool.GetBuffer()
	defer pool.PutBuffer(buf1)
	for i, v := range req.NodeId {
		buf1.AppendUint32(v)
		if i != len(req.NodeId)-1 {
			buf1.AppendByte(',')
		}
	}
	nodeid := buf1.String()
	buf2 := pool.GetBuffer()
	defer pool.PutBuffer(buf2)
	for i, v := range req.PnodeId {
		buf1.AppendUint32(v)
		if i != len(req.PnodeId)-1 {
			buf2.AppendByte(',')
		}
	}
	pnodeid := buf2.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[MoveNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if e := s.permissionDao.MongoMoveNode(ctx, operator, nodeid, pnodeid); e != nil {
		log.Error(ctx, "[MoveNode] operator:", md["Token-Data"], "nodeid:", nodeid, "new parent nodeid:", pnodeid, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.MoveNodeResp{}, nil
}
func (s *Service) DelNode(ctx context.Context, req *api.DelNodeReq) (*api.DelNodeResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	if req.NodeId[1] == 1 || req.NodeId[1] == 2 || req.NodeId[1] == 3 {
		//0,1 -> UserControl
		//0,2 -> RoleControl
		//0,3 -> ConfigControl
		//these are default,can't modify
		return nil, ecode.ErrReq
	}
	buf1 := pool.GetBuffer()
	defer pool.PutBuffer(buf1)
	for i, v := range req.NodeId {
		buf1.AppendUint32(v)
		if i != len(req.NodeId)-1 {
			buf1.AppendByte(',')
		}
	}
	nodeid := buf1.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[DelNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	if e = s.permissionDao.MongoDeleteNode(ctx, operator, nodeid); e != nil {
		log.Error(ctx, "[DelNode] operator:", md["Token-Data"], "nodeid:", nodeid, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.DelNodeResp{}, nil
}
func (s *Service) ListUserNode(ctx context.Context, req *api.ListUserNodeReq) (*api.ListUserNodeResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[ListUserNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	var target primitive.ObjectID
	if req.UserId == "" {
		//list self's
		req.UserId = md["Token-Data"]
		target = operator
	} else {
		//list other user's
		target, e = primitive.ObjectIDFromHex(req.UserId)
		if e != nil {
			log.Error(ctx, "[ListUserNode] target:", req.UserId, "format wrong:", e)
			return nil, ecode.ErrReq
		}
		if !operator.IsZero() {
			//user control permission check
			canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.UserControl, true)
			if e != nil {
				log.Error(ctx, "[ListUserNode] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if !canread && !admin {
				return nil, ecode.ErrPermission
			}
		}
	}
	//logic
	undup := make(map[string]*api.NodeInfo)
	usernodes, e := s.permissionDao.MongoGetUserNodes(ctx, target, project, nil)
	if e != nil {
		log.Error(ctx, "[ListUserNode] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	for _, usernode := range usernodes {
		if exist, ok := undup[usernode.NodeId]; ok {
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
			nodeid, e := util.ParseNodeIDstr(usernode.NodeId)
			if e != nil {
				log.Error(ctx, "[ListUserNode] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, "nodeid:", usernode.NodeId, "format wrong:", e)
				return nil, ecode.ErrSystem
			}
			undup[usernode.NodeId] = &api.NodeInfo{
				NodeId:   nodeid,
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
		userrolenodes, e := s.permissionDao.MongoGetUserRoleNodes(ctx, target, project, nil)
		if e != nil {
			log.Error(ctx, "[ListUserNode] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, e)
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		for _, v := range userrolenodes {
			for _, userrolenode := range v {
				if exist, ok := undup[userrolenode.NodeId]; ok {
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
					nodeid, e := util.ParseNodeIDstr(userrolenode.NodeId)
					if e != nil {
						log.Error(ctx, "[ListUserNode] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, "nodeid:", userrolenode.NodeId, "format wrong:", e)
						return nil, ecode.ErrSystem
					}
					undup[userrolenode.NodeId] = &api.NodeInfo{
						NodeId:   nodeid,
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
	index := make(map[*api.NodeInfo]string, len(undup))
	for nodeidstr, v := range undup {
		nodes = append(nodes, v)
		index[v] = nodeidstr
	}
	sort.Slice(nodes, func(i, j int) bool {
		return len(nodes[i].NodeId) < len(nodes[j].NodeId)
	})
	adminnodeids := make([]string, 0, len(nodes))
	nodeids := make([]string, 0, len(nodes))
	for _, node := range nodes {
		//if this node is below another admin node,this node can be jumped
		nodeid := index[node]
		jump := false
		for _, adminnodeid := range adminnodeids {
			if strings.HasPrefix(adminnodeid, nodeid) {
				jump = true
				break
			}
		}
		if jump {
			continue
		}
		if node.Admin {
			adminnodeids = append(adminnodeids, nodeid)
		}
		nodeids = append(nodeids, nodeid)
	}
	nodeinfos, e := s.permissionDao.MongoGetNodes(ctx, nodeids)
	if e != nil {
		log.Error(ctx, "[ListUserNode] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	for _, nodeinfo := range nodeinfos {
		exist, _ := undup[nodeinfo.NodeId]
		exist.NodeName = nodeinfo.NodeName
		exist.NodeData = nodeinfo.NodeData
	}
	nodes = nodes[:0]
	for _, nodeid := range nodeids {
		nodes = append(nodes, undup[nodeid])
	}
	g := egroup.GetGroup(ctx)
	for _, v := range nodes {
		if !v.Admin {
			continue
		}
		node := v
		nodeid := index[v]
		g.Go(func(gctx context.Context) error {
			tmpnodes, e := s.permissionDao.MongoListNodes(ctx, nodeid, true)
			if e != nil {
				log.Error(gctx, "[ListUserNode] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, "admin nodeid:", nodeid, e)
				return e
			}
			sort.Slice(tmpnodes, func(i, j int) bool {
				return strings.Count(tmpnodes[i].NodeId, ",") < strings.Count(tmpnodes[j].NodeId, ",")
			})
			for _, tmpnode := range tmpnodes {
				tmpnodeid, e := util.ParseNodeIDstr(tmpnode.NodeId)
				if e != nil {
					log.Error(ctx, "[ListUserNode] operator:", md["Token-Data"], "project:", project, "target:", req.UserId, "nodeid:", tmpnode.NodeId, "format wrong:", e)
					return ecode.ErrSystem
				}
				addTreeNode(node, &api.NodeInfo{
					NodeId:   tmpnodeid,
					NodeName: tmpnode.NodeName,
					NodeData: tmpnode.NodeData,
					Canread:  true,
					Canwrite: true,
					Admin:    true,
				})
			}
			return nil
		})
	}
	if e := egroup.PutGroup(g); e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	root := &api.NodeInfo{
		NodeId: []uint32{0},
	}
	for _, node := range nodes {
		addTreeNode(root, node)
	}
	return &api.ListUserNodeResp{Nodes: root.Children}, nil
}
func (s *Service) ListRoleNode(ctx context.Context, req *api.ListRoleNodeReq) (*api.ListRoleNodeResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-Data"])
	if e != nil {
		log.Error(ctx, "[ListRoleNode] operator:", md["Token-Data"], "format wrong:", e)
		return nil, ecode.ErrToken
	}
	//role control permission check
	canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, project+model.RoleControl, true)
	if e != nil {
		log.Error(ctx, "[ListRoleNode] operator:", md["Token-Data"], "project:", project, "get permission failed:", e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if !canread && !admin {
		return nil, ecode.ErrPermission
	}
	//logic
	rolenodes, e := s.permissionDao.MongoGetRoleNodes(ctx, project, req.RoleName, nil)
	if e != nil {
		log.Error(ctx, "[ListRoleNode] operator:", md["Token-Data"], "project:", project, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodes := make([]*api.NodeInfo, 0, len(rolenodes))
	undup := make(map[string]*api.NodeInfo, len(rolenodes))
	index := make(map[*api.NodeInfo]string)
	for _, rolenode := range rolenodes {
		if exist, ok := undup[rolenode.NodeId]; ok {
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
			nodeid, e := util.ParseNodeIDstr(rolenode.NodeId)
			if e != nil {
				log.Error(ctx, "[ListRoleNode] operator:", md["Token-Data"], "project:", project, "rolename:", req.RoleName, "nodeid:", rolenode.NodeId, "format wrong:", e)
				return nil, ecode.ErrSystem
			}
			tmp := &api.NodeInfo{
				NodeId:   nodeid,
				NodeName: "",
				NodeData: "",
				Canread:  rolenode.R,
				Canwrite: rolenode.W,
				Admin:    rolenode.X,
			}
			nodes = append(nodes, tmp)
			undup[rolenode.NodeId] = tmp
			index[tmp] = rolenode.NodeId
		}
	}
	sort.Slice(nodes, func(i, j int) bool {
		return len(nodes[i].NodeId) < len(nodes[j].NodeId)
	})
	adminnodeids := make([]string, 0, len(nodes))
	nodeids := make([]string, 0, len(nodes))
	for _, node := range nodes {
		//if this node is below another admin node,this node can be jumped
		nodeid := index[node]
		jump := false
		for _, adminnodeid := range adminnodeids {
			if strings.HasPrefix(adminnodeid, nodeid) {
				jump = true
				break
			}
		}
		if jump {
			continue
		}
		if node.Admin {
			adminnodeids = append(adminnodeids, nodeid)
		}
		nodeids = append(nodeids, nodeid)
	}
	nodeinfos, e := s.permissionDao.MongoGetNodes(ctx, nodeids)
	if e != nil {
		log.Error(ctx, "[ListRoleNode] operator:", md["Token-Data"], "project:", project, "rolename:", req.RoleName, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	for _, nodeinfo := range nodeinfos {
		exist, _ := undup[nodeinfo.NodeId]
		exist.NodeName = nodeinfo.NodeName
		exist.NodeData = nodeinfo.NodeData
	}
	nodes = nodes[:0]
	for _, nodeid := range nodeids {
		nodes = append(nodes, undup[nodeid])
	}
	g := egroup.GetGroup(ctx)
	for _, v := range nodes {
		if !v.Admin {
			continue
		}
		node := v
		nodeid := index[node]
		g.Go(func(gctx context.Context) error {
			tmpnodes, e := s.permissionDao.MongoListNodes(ctx, nodeid, true)
			if e != nil {
				log.Error(gctx, "[ListRoleNode] operator:", md["Token-Data"], "project:", project, "rolename:", req.RoleName, "admin nodeid:", nodeid, e)
				return e
			}
			sort.Slice(tmpnodes, func(i, j int) bool {
				return strings.Count(tmpnodes[i].NodeId, ",") < strings.Count(tmpnodes[j].NodeId, ",")
			})
			for _, tmpnode := range tmpnodes {
				tmpnodeid, e := util.ParseNodeIDstr(tmpnode.NodeId)
				if e != nil {
					log.Error(ctx, "[ListRoleNode] operator:", md["Token-Data"], "project:", project, "rolename:", req.RoleName, "nodeid:", tmpnode.NodeId, "format wrong:", e)
					return ecode.ErrSystem
				}
				addTreeNode(node, &api.NodeInfo{
					NodeId:   tmpnodeid,
					NodeName: tmpnode.NodeName,
					NodeData: tmpnode.NodeData,
					Canread:  true,
					Canwrite: true,
					Admin:    true,
				})
			}
			return nil
		})
	}
	if e := egroup.PutGroup(g); e != nil {
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	root := &api.NodeInfo{
		NodeId: []uint32{0},
	}
	for _, node := range nodes {
		addTreeNode(root, node)
	}
	return &api.ListRoleNodeResp{Nodes: root.Children}, nil
}
func (s *Service) ListProjectNode(ctx context.Context, req *api.ListProjectNodeReq) (*api.ListProjectNodeResp, error) {
	if req.Project[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	for i, v := range req.Project {
		buf.AppendUint32(v)
		if i != len(req.Project)-1 {
			buf.AppendByte(',')
		}
	}
	project := buf.String()
	md := metadata.GetMetadata(ctx)
	nodes, e := s.permissionDao.MongoListNodes(ctx, project, true)
	if e != nil {
		log.Error(ctx, "[ListAllNode] operator:", md["Token-Data"], "project:", req.Project, e)
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	projectnode := &api.NodeInfo{
		NodeId:   req.Project,
		Children: make([]*api.NodeInfo, 0, 10),
	}
	sort.Slice(nodes, func(i, j int) bool {
		return strings.Count(nodes[i].NodeId, ",") < strings.Count(nodes[j].NodeId, ",")
	})
	for _, node := range nodes {
		nodeid, e := util.ParseNodeIDstr(node.NodeId)
		if e != nil {
			log.Error(ctx, "[ListAllNode] operator:", md["Token-Data"], "project:", project, "nodeid:", node.NodeId, "format wrong:", e)
			return nil, ecode.ErrSystem
		}
		addTreeNode(projectnode, &api.NodeInfo{
			NodeId:   nodeid,
			NodeName: node.NodeName,
			NodeData: node.NodeData,
			Children: make([]*api.NodeInfo, 0),
		})
	}
	return &api.ListProjectNodeResp{Nodes: projectnode.Children}, nil
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
