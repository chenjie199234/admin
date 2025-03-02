package permission

import (
	"context"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/pool/bpool"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/v2/bson"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
)

// Service subservice for permission business
type Service struct {
	stop *graceful.Graceful

	userDao       *userdao.Dao
	permissionDao *permissiondao.Dao
}

// Start -
func Start() *Service {
	return &Service{
		stop: graceful.New(),

		userDao:       userdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}
}
func (s *Service) GetUserPermission(ctx context.Context, req *api.GetUserPermissionReq) (*api.GetUserPermissionResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.NodeId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	nodeid := common.BTS(buf)

	target, e := bson.ObjectIDFromHex(req.UserId)
	if e != nil {
		slog.ErrorContext(ctx, "[GetUserPermission] target's userid format wrong", slog.String("user_id", req.UserId))
		return nil, ecode.ErrReq
	}
	canread, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, target, nodeid, true)
	if e != nil {
		slog.ErrorContext(ctx, "[GetUserPermission] db op failed", slog.String("user_id", req.UserId), slog.String("node_id", nodeid), slog.String("error", e.Error()))
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
	buf1 := bpool.Get(0)
	defer bpool.Put(&buf1)
	for i, v := range req.NodeId {
		if i != 0 {
			buf1 = append(buf1, ',')
		}
		buf1 = strconv.AppendUint(buf1, uint64(v), 10)
	}
	nodeid := common.BTS(buf1)

	buf2 := bpool.Get(0)
	defer bpool.Put(&buf2)
	buf2 = strconv.AppendUint(buf2, uint64(req.NodeId[0]), 10)
	buf2 = append(buf2, ',')
	buf2 = strconv.AppendUint(buf2, uint64(req.NodeId[1]), 10)
	projectid := common.BTS(buf2)

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateUserPermission] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	target, e := bson.ObjectIDFromHex(req.UserId)
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateUserPermission] target's userid format wrong", slog.String("user_id", req.UserId))
		return nil, ecode.ErrReq
	}
	if !operator.IsZero() {
		//permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateUserPermission] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}
	//logic
	if e = s.permissionDao.MongoUpdateUserPermission(ctx, operator, target, nodeid, req.Admin, req.Canread, req.Canwrite); e != nil {
		slog.ErrorContext(ctx, "[UpdateUserPermission] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("user_id", req.UserId),
			slog.String("node_id", nodeid),
			slog.Any("new_permission", []bool{req.Canread, req.Canwrite, req.Admin}),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[UpdateUserPermission] success",
		slog.String("project_id", projectid),
		slog.String("user_id", req.UserId),
		slog.String("node_id", nodeid),
		slog.Any("new_permission", []bool{req.Canread, req.Canwrite, req.Admin}))
	return &api.UpdateUserPermissionResp{}, nil
}
func (s *Service) UpdateRolePermission(ctx context.Context, req *api.UpdateRolePermissionReq) (*api.UpdateRolePermissionResp, error) {
	if req.NodeId[0] != 0 || req.ProjectId[0] != 0 || (req.ProjectId[1] != req.NodeId[1]) {
		return nil, ecode.ErrReq
	}
	if !req.Admin && req.Canwrite && !req.Canread {
		return nil, ecode.ErrReq
	}
	buf1 := bpool.Get(0)
	defer bpool.Put(&buf1)
	for i, v := range req.NodeId {
		if i != 0 {
			buf1 = append(buf1, ',')
		}
		buf1 = strconv.AppendUint(buf1, uint64(v), 10)
	}
	nodeid := common.BTS(buf1)

	buf2 := bpool.Get(0)
	defer bpool.Put(&buf2)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf2 = append(buf2, ',')
		}
		buf2 = strconv.AppendUint(buf2, uint64(v), 10)
	}
	projectid := common.BTS(buf2)

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateRolePermission] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//role control permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateRolePermission] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}
	//logic
	if e = s.permissionDao.MongoUpdateRolePermission(ctx, operator, projectid, req.RoleName, nodeid, req.Admin, req.Canread, req.Canwrite); e != nil {
		slog.ErrorContext(ctx, "[UpdateRolePermission] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("role_name", req.RoleName),
			slog.String("node_id", nodeid),
			slog.Any("new_permission", []bool{req.Canread, req.Canwrite, req.Admin}),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[UpdateRolePermission] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("role_name", req.RoleName),
		slog.String("node_id", nodeid),
		slog.Any("new_permission", []bool{req.Canread, req.Canwrite, req.Admin}))
	return &api.UpdateRolePermissionResp{}, nil
}
func (s *Service) AddNode(ctx context.Context, req *api.AddNodeReq) (*api.AddNodeResp, error) {
	if req.PnodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	if req.PnodeId[1] == 1 {
		//0,1 -> admin project
		//can't add
		return nil, ecode.ErrPermission
	}
	if len(req.PnodeId) >= 3 && (req.PnodeId[2] == 1 || req.PnodeId[2] == 2) {
		//0,x,1 -> UserAndRoleControl
		//0,x,2 -> AppControl
		//these are default,already exist
		return nil, ecode.ErrPermission
	}
	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.PnodeId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	pnodeid := common.BTS(buf)

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[AddNode] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	var nodeidstr string
	if nodeidstr, e = s.permissionDao.MongoAddNode(ctx, operator, pnodeid, req.NodeName, req.NodeData); e != nil {
		slog.ErrorContext(ctx, "[AddNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("node_name", req.NodeName),
			slog.String("node_data", req.NodeData),
			slog.String("parent_node_id", pnodeid),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeid, _ := util.ParseNodeIDstr(nodeidstr)
	slog.InfoContext(ctx, "[AddNode] success",
		slog.String("operator", md["Token-User"]),
		slog.String("node_id", nodeidstr),
		slog.String("node_name", req.NodeName),
		slog.String("node_data", req.NodeData))
	return &api.AddNodeResp{NodeId: nodeid}, nil
}
func (s *Service) UpdateNode(ctx context.Context, req *api.UpdateNodeReq) (*api.UpdateNodeResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	if req.NodeId[1] == 1 {
		//0,1 -> admin project
		//can't update
		return nil, ecode.ErrPermission
	}
	if len(req.NodeId) >= 3 && (req.NodeId[2] == 1 || req.NodeId[2] == 2) {
		//0,x,1 -> UserAndRoleControl
		//0,x,2 -> AppControl
		//these are default,can't update
		return nil, ecode.ErrPermission
	}
	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.NodeId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	nodeid := common.BTS(buf)

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateNode] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	oldnode, e := s.permissionDao.MongoUpdateNode(ctx, operator, nodeid, req.NewNodeName, req.NewNodeData)
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("new_node_name", req.NewNodeName),
			slog.String("new_node_data", req.NewNodeData),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if oldnode.NodeName != req.NewNodeName && oldnode.NodeData != req.NewNodeData {
		slog.InfoContext(ctx, "[UpdateNode] success",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("old_node_name", oldnode.NodeName),
			slog.String("new_node_name", req.NewNodeName),
			slog.String("old_node_data", oldnode.NodeData),
			slog.String("new_node_data", req.NewNodeData))
	} else if oldnode.NodeName != req.NewNodeName {
		slog.InfoContext(ctx, "[UpdateNode] success",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("old_node_name", oldnode.NodeName),
			slog.String("new_node_name", req.NewNodeName))
	} else if oldnode.NodeData != req.NewNodeData {
		slog.InfoContext(ctx, "[UpdateNode] success",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("old_node_data", oldnode.NodeData),
			slog.String("new_node_data", req.NewNodeData))
	} else {
		slog.InfoContext(ctx, "[UpdateNode] success,nothing changed",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("node_name", oldnode.NodeName),
			slog.String("node_data", oldnode.NodeData))
	}
	return &api.UpdateNodeResp{}, nil
}
func (s *Service) MoveNode(ctx context.Context, req *api.MoveNodeReq) (*api.MoveNodeResp, error) {
	if req.NodeId[0] != 0 || req.PnodeId[0] != 0 || (req.NodeId[1] != req.PnodeId[1]) {
		//can't cross project
		return nil, ecode.ErrReq
	}
	if req.NodeId[1] == 1 || req.PnodeId[1] == 1 {
		//0,1 -> admin project
		//can't modify
		return nil, ecode.ErrPermission
	}
	if req.NodeId[2] == 1 || req.NodeId[2] == 2 {
		//0,x,1 -> UserAndRoleControl
		//0,x,2 -> AppControl
		//these are default,can't modify
		return nil, ecode.ErrPermission
	}
	if len(req.PnodeId) >= 3 && (req.PnodeId[2] == 1 || req.PnodeId[2] == 2) {
		//0,x,1 -> UserAndRoleControl
		//0,x,2 -> AppControl
		//these are default,can't modify
		return nil, ecode.ErrPermission
	}
	if len(req.PnodeId)+1 == len(req.NodeId) {
		//0,x,y,z move to 0,x,y is equal to not move
		child := true
		for i := range req.PnodeId {
			if req.PnodeId[i] != req.NodeId[i] {
				child = false
				break
			}
		}
		if child {
			return &api.MoveNodeResp{}, nil
		}
	}
	buf1 := bpool.Get(0)
	defer bpool.Put(&buf1)
	for i, v := range req.NodeId {
		if i != 0 {
			buf1 = append(buf1, ',')
		}
		buf1 = strconv.AppendUint(buf1, uint64(v), 10)
	}
	nodeid := common.BTS(buf1)

	buf2 := bpool.Get(0)
	defer bpool.Put(&buf2)
	for i, v := range req.PnodeId {
		if i != 0 {
			buf2 = append(buf2, ',')
		}
		buf2 = strconv.AppendUint(buf2, uint64(v), 10)
	}
	pnodeid := common.BTS(buf2)

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[MoveNode] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	newnodeid, e := s.permissionDao.MongoMoveNode(ctx, operator, nodeid, pnodeid)
	if e != nil {
		slog.ErrorContext(ctx, "[MoveNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("new_parent_node_id", pnodeid),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[MoveNode] success",
		slog.String("operator", md["Token-User"]),
		slog.String("node_id", nodeid),
		slog.String("new_node_id", newnodeid))
	return &api.MoveNodeResp{}, nil
}
func (s *Service) DelNode(ctx context.Context, req *api.DelNodeReq) (*api.DelNodeResp, error) {
	if req.NodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	if req.NodeId[1] == 1 {
		//0,1 -> admin project
		//can't delete
		return nil, ecode.ErrPermission
	}
	if req.NodeId[2] == 1 || req.NodeId[2] == 2 {
		//0,x,1 -> UserAndRoleControl node
		//0,x,2 -> AppControl node
		//these are default,can't modify
		return nil, ecode.ErrPermission
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[DelNode] operator's token format failed", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.NodeId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	nodeid := common.BTS(buf)

	node, e := s.permissionDao.MongoDeleteNode(ctx, operator, nodeid)
	if e != nil {
		slog.ErrorContext(ctx, "[DelNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[DelNode] success",
		slog.String("operator", md["Token-User"]),
		slog.String("node_id", nodeid),
		slog.String("node_name", node.NodeName),
		slog.String("node_data", node.NodeData))
	return &api.DelNodeResp{}, nil
}
func (s *Service) ListUserNode(ctx context.Context, req *api.ListUserNodeReq) (*api.ListUserNodeResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[ListUserNode] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	var target bson.ObjectID
	if req.UserId == "" || req.UserId == md["Token-User"] {
		//list self's
		req.UserId = md["Token-User"]
		target = operator
	} else {
		//list other user's
		target, e = bson.ObjectIDFromHex(req.UserId)
		if e != nil {
			slog.ErrorContext(ctx, "[ListUserNode] target's userid format wrong", slog.String("user_id", req.UserId))
			return nil, ecode.ErrReq
		}
	}
	//get self's don't need check permission
	//get other's need check permission
	if operator.Hex() != target.Hex() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[ListUserNode] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}
	//logic
	root, e := s.permissionDao.MongoGetNode(ctx, projectid)
	if e != nil {
		slog.ErrorContext(ctx, "[ListUserNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("user_id", req.UserId),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if root == nil {
		return nil, ecode.ErrProjectNotExist
	}
	children, e := s.permissionDao.MongoListChildrenNodes(ctx, projectid, true)
	if e != nil {
		slog.ErrorContext(ctx, "[ListUserNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("user_id", req.UserId),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var usernodes model.UserNodes
	if !target.IsZero() {
		if usernodes, e = s.permissionDao.MongoGetUserNodes(ctx, target, projectid, nil); e != nil {
			slog.ErrorContext(ctx, "[ListUserNode] db op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("user_id", req.UserId),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	var userrolenodes map[string]model.RoleNodes
	if !target.IsZero() && req.NeedUserRoleNode {
		if userrolenodes, e = s.permissionDao.MongoGetUserRoleNodes(ctx, target, projectid, nil); e != nil {
			slog.ErrorContext(ctx, "[ListUserNode] db op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("user_id", req.UserId),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	projectnode := &api.NodeInfo{
		NodeId:   req.ProjectId,
		NodeName: root.NodeName,
		NodeData: root.NodeData,
		Children: make([]*api.NodeInfo, 0, 10),
	}
	if target.IsZero() {
		projectnode.Canread = true
		projectnode.Canwrite = true
		projectnode.Admin = true
	} else {
		projectnode.Canread, projectnode.Canwrite, projectnode.Admin = usernodes.CheckNode(projectid)
		for _, rolenodes := range userrolenodes {
			if projectnode.Admin {
				break
			}
			tmpr, tmpw, tmpa := rolenodes.CheckNode(projectid)
			if tmpr {
				projectnode.Canread = true
			}
			if tmpw {
				projectnode.Canwrite = true
			}
			if tmpa {
				projectnode.Admin = true
			}
		}
	}
	sort.Slice(children, func(i, j int) bool {
		return strings.Count(children[i].NodeId, ",") < strings.Count(children[j].NodeId, ",")
	})
	for _, node := range children {
		nodeid, e := util.ParseNodeIDstr(node.NodeId)
		if e != nil {
			slog.ErrorContext(ctx, "[ListUserNode] target's node's nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("user_id", req.UserId),
				slog.String("node_id", node.NodeId))
			return nil, ecode.ErrSystem
		}
		tmp := &api.NodeInfo{
			NodeId:   nodeid,
			NodeName: node.NodeName,
			NodeData: node.NodeData,
			Children: make([]*api.NodeInfo, 0, 10),
		}
		if target.IsZero() {
			tmp.Canread = true
			tmp.Canwrite = true
			tmp.Admin = true
		} else {
			tmp.Canread, tmp.Canwrite, tmp.Admin = usernodes.CheckNode(node.NodeId)
			for _, rolenodes := range userrolenodes {
				if tmp.Admin {
					break
				}
				tmpr, tmpw, tmpa := rolenodes.CheckNode(node.NodeId)
				if tmpr {
					tmp.Canread = tmpr
				}
				if tmpw {
					tmp.Canwrite = tmpw
				}
				if tmpa {
					tmp.Admin = tmpa
				}
			}
		}
		addTreeNode(projectnode, tmp)
	}
	sortTreeNodes(projectnode.Children)
	return &api.ListUserNodeResp{Node: projectnode}, nil
}
func (s *Service) ListRoleNode(ctx context.Context, req *api.ListRoleNodeReq) (*api.ListRoleNodeResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[ListRoleNode] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	if !operator.IsZero() {
		//permission check
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.UserAndRoleControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[ListRoleNode] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	root, e := s.permissionDao.MongoGetNode(ctx, projectid)
	if e != nil {
		slog.ErrorContext(ctx, "[ListRoleNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("role_name", req.RoleName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if root == nil {
		return nil, ecode.ErrProjectNotExist
	}
	children, e := s.permissionDao.MongoListChildrenNodes(ctx, projectid, true)
	if e != nil {
		slog.ErrorContext(ctx, "[ListRoleNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("role_name", req.RoleName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	rolenodes, e := s.permissionDao.MongoGetRoleNodes(ctx, projectid, req.RoleName, nil)
	if e != nil {
		slog.ErrorContext(ctx, "[ListRoleNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("role_name", req.RoleName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	projectnode := &api.NodeInfo{
		NodeId:   req.ProjectId,
		NodeName: root.NodeName,
		NodeData: root.NodeData,
		Children: make([]*api.NodeInfo, 0, 10),
	}
	projectnode.Canread, projectnode.Canwrite, projectnode.Admin = rolenodes.CheckNode(projectid)
	sort.Slice(children, func(i, j int) bool {
		return strings.Count(children[i].NodeId, ",") < strings.Count(children[j].NodeId, ",")
	})
	for _, node := range children {
		nodeid, e := util.ParseNodeIDstr(node.NodeId)
		if e != nil {
			slog.ErrorContext(ctx, "[ListRoleNode] role's node's nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("role_name", req.RoleName),
				slog.String("node_id", node.NodeId))
			return nil, ecode.ErrSystem
		}
		tmp := &api.NodeInfo{
			NodeId:   nodeid,
			NodeName: node.NodeName,
			NodeData: node.NodeData,
			Children: make([]*api.NodeInfo, 0, 10),
		}
		tmp.Canread, tmp.Canwrite, tmp.Admin = rolenodes.CheckNode(node.NodeId)
		addTreeNode(projectnode, tmp)
	}
	sortTreeNodes(projectnode.Children)
	return &api.ListRoleNodeResp{Node: projectnode}, nil
}
func (s *Service) ListProjectNode(ctx context.Context, req *api.ListProjectNodeReq) (*api.ListProjectNodeResp, error) {
	if req.ProjectId[0] != 0 {
		return nil, ecode.ErrReq
	}
	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	md := metadata.GetMetadata(ctx)

	root, e := s.permissionDao.MongoGetNode(ctx, projectid)
	if e != nil {
		slog.ErrorContext(ctx, "[ListProjectNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("error", e.Error()))
		return nil, e
	}
	if root == nil {
		return nil, ecode.ErrProjectNotExist
	}
	children, e := s.permissionDao.MongoListChildrenNodes(ctx, projectid, true)
	if e != nil {
		slog.ErrorContext(ctx, "[ListProjectNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	projectnode := &api.NodeInfo{
		NodeId:   req.ProjectId,
		NodeName: root.NodeName,
		NodeData: root.NodeData,
		Children: make([]*api.NodeInfo, 0, 10),
	}
	sort.Slice(children, func(i, j int) bool {
		return strings.Count(children[i].NodeId, ",") < strings.Count(children[j].NodeId, ",")
	})
	for _, node := range children {
		nodeid, e := util.ParseNodeIDstr(node.NodeId)
		if e != nil {
			slog.ErrorContext(ctx, "[ListProjectNode] project's node's nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("node_id", node.NodeId))
			return nil, ecode.ErrSystem
		}
		addTreeNode(projectnode, &api.NodeInfo{
			NodeId:   nodeid,
			NodeName: node.NodeName,
			NodeData: node.NodeData,
			Children: make([]*api.NodeInfo, 0),
		})
	}
	sortTreeNodes(projectnode.Children)
	return &api.ListProjectNodeResp{Node: projectnode}, nil
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
func sortTreeNodes(nodes []*api.NodeInfo) {
	sort.Slice(nodes, func(i, j int) bool {
		if len(nodes[i].NodeId) < len(nodes[j].NodeId) {
			return true
		} else if len(nodes[i].NodeId) > len(nodes[j].NodeId) {
			return false
		}
		for k := range len(nodes[i].NodeId) {
			if nodes[i].NodeId[k] < nodes[j].NodeId[k] {
				return true
			}
		}
		return false
	})
	for _, node := range nodes {
		if len(node.Children) > 1 {
			sortTreeNodes(node.Children)
		}
	}
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
