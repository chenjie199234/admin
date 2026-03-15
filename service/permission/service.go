package permission

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
	//"github.com/chenjie199234/Corelib/web"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/util/graceful"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// Service subservice for permission business
type Service struct {
	stop *graceful.Graceful

	userDao       *userdao.Dao
	permissionDao *permissiondao.Dao
}

// Start -
func Start() (*Service, error) {
	return &Service{
		stop: graceful.New(),

		userDao:       userdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}, nil
}
func (s *Service) GetUserPermission(ctx context.Context, req *api.GetUserPermissionReq) (*api.GetUserPermissionResp, error) {
	if req.GetNodeId()[0] != 0 {
		return nil, ecode.ErrReq
	}

	nodeid := util.FormID(req.GetNodeId())

	target, e := bson.ObjectIDFromHex(req.GetUserId())
	if e != nil {
		slog.ErrorContext(ctx, "[GetUserPermission] target's userid format wrong", slog.String("user_id", req.GetUserId()))
		return nil, ecode.ErrReq
	}
	canread, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, target, nodeid, true)
	if e != nil {
		slog.ErrorContext(ctx, "[GetUserPermission] db op failed",
			slog.String("user_id", req.GetUserId()), slog.String("node_id", nodeid), slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.GetUserPermissionResp{}
	resp.SetCanread(canread)
	resp.SetCanwrite(canwrite)
	resp.SetAdmin(admin)
	return resp, nil
}
func (s *Service) UpdateUserPermission(ctx context.Context, req *api.UpdateUserPermissionReq) (*api.UpdateUserPermissionResp, error) {
	if req.GetNodeId()[0] != 0 {
		return nil, ecode.ErrReq
	}
	if !req.GetAdmin() && req.GetCanwrite() && !req.GetCanread() {
		return nil, ecode.ErrReq
	}

	nodeid := util.FormID(req.GetNodeId())
	projectid := util.FormID(req.GetNodeId()[:2])

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateUserPermission] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	target, e := bson.ObjectIDFromHex(req.GetUserId())
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateUserPermission] target's userid format wrong", slog.String("user_id", req.GetUserId()))
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
	if e = s.permissionDao.MongoUpdateUserPermission(ctx, operator, target, nodeid, req.GetAdmin(), req.GetCanread(), req.GetCanwrite()); e != nil {
		slog.ErrorContext(ctx, "[UpdateUserPermission] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("user_id", req.GetUserId()),
			slog.String("node_id", nodeid),
			slog.Any("new_permission", []bool{req.GetCanread(), req.GetCanwrite(), req.GetAdmin()}),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[UpdateUserPermission] success",
		slog.String("project_id", projectid),
		slog.String("user_id", req.GetUserId()),
		slog.String("node_id", nodeid),
		slog.Any("new_permission", []bool{req.GetCanread(), req.GetCanwrite(), req.GetAdmin()}))
	return &api.UpdateUserPermissionResp{}, nil
}
func (s *Service) UpdateRolePermission(ctx context.Context, req *api.UpdateRolePermissionReq) (*api.UpdateRolePermissionResp, error) {
	if req.GetNodeId()[0] != 0 || req.GetProjectId()[0] != 0 || (req.GetProjectId()[1] != req.GetNodeId()[1]) {
		return nil, ecode.ErrReq
	}
	if !req.GetAdmin() && req.GetCanwrite() && !req.GetCanread() {
		return nil, ecode.ErrReq
	}

	nodeid := util.FormID(req.GetNodeId())
	projectid := util.FormID(req.GetProjectId())

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
	if e = s.permissionDao.MongoUpdateRolePermission(ctx, operator, projectid, req.GetRoleName(), nodeid, req.GetAdmin(), req.GetCanread(), req.GetCanwrite()); e != nil {
		slog.ErrorContext(ctx, "[UpdateRolePermission] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("role_name", req.GetRoleName()),
			slog.String("node_id", nodeid),
			slog.Any("new_permission", []bool{req.GetCanread(), req.GetCanwrite(), req.GetAdmin()}),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[UpdateRolePermission] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("role_name", req.GetRoleName()),
		slog.String("node_id", nodeid),
		slog.Any("new_permission", []bool{req.GetCanread(), req.GetCanwrite(), req.GetAdmin()}))
	return &api.UpdateRolePermissionResp{}, nil
}
func (s *Service) AddNode(ctx context.Context, req *api.AddNodeReq) (*api.AddNodeResp, error) {
	if req.GetPnodeId()[0] != 0 {
		return nil, ecode.ErrReq
	}
	if req.GetPnodeId()[1] == 1 {
		//0,1 -> admin project
		//can't add
		return nil, ecode.ErrPermission
	}
	if len(req.GetPnodeId()) >= 3 && (req.GetPnodeId()[2] == 1 || req.GetPnodeId()[2] == 2) {
		//0,x,1 -> UserAndRoleControl
		//0,x,2 -> AppControl
		//these are default,already exist
		return nil, ecode.ErrPermission
	}

	pnodeid := util.FormID(req.GetPnodeId())

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[AddNode] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	var nodeidstr string
	if nodeidstr, e = s.permissionDao.MongoAddNode(ctx, operator, pnodeid, req.GetNodeName(), req.GetNodeData()); e != nil {
		slog.ErrorContext(ctx, "[AddNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("node_name", req.GetNodeName()),
			slog.String("node_data", req.GetNodeData()),
			slog.String("parent_node_id", pnodeid),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeid, _ := util.ParseID(nodeidstr)
	slog.InfoContext(ctx, "[AddNode] success",
		slog.String("operator", md["Token-User"]),
		slog.String("node_id", nodeidstr),
		slog.String("node_name", req.GetNodeName()),
		slog.String("node_data", req.GetNodeData()))
	resp := &api.AddNodeResp{}
	resp.SetNodeId(nodeid)
	return resp, nil
}
func (s *Service) UpdateNode(ctx context.Context, req *api.UpdateNodeReq) (*api.UpdateNodeResp, error) {
	if req.GetNodeId()[0] != 0 {
		return nil, ecode.ErrReq
	}
	if req.GetNodeId()[1] == 1 {
		//0,1 -> admin project
		//can't update
		return nil, ecode.ErrPermission
	}
	if len(req.GetNodeId()) >= 3 && (req.GetNodeId()[2] == 1 || req.GetNodeId()[2] == 2) {
		//0,x,1 -> UserAndRoleControl
		//0,x,2 -> AppControl
		//these are default,can't update
		return nil, ecode.ErrPermission
	}

	nodeid := util.FormID(req.GetNodeId())

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateNode] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	oldnode, e := s.permissionDao.MongoUpdateNode(ctx, operator, nodeid, req.GetNewNodeName(), req.GetNewNodeData())
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("new_node_name", req.GetNewNodeName()),
			slog.String("new_node_data", req.GetNewNodeData()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	if oldnode.NodeName != req.GetNewNodeName() && oldnode.NodeData != req.GetNewNodeData() {
		slog.InfoContext(ctx, "[UpdateNode] success",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("old_node_name", oldnode.NodeName),
			slog.String("new_node_name", req.GetNewNodeName()),
			slog.String("old_node_data", oldnode.NodeData),
			slog.String("new_node_data", req.GetNewNodeData()))
	} else if oldnode.NodeName != req.GetNewNodeName() {
		slog.InfoContext(ctx, "[UpdateNode] success",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("old_node_name", oldnode.NodeName),
			slog.String("new_node_name", req.GetNewNodeName()))
	} else if oldnode.NodeData != req.GetNewNodeData() {
		slog.InfoContext(ctx, "[UpdateNode] success",
			slog.String("operator", md["Token-User"]),
			slog.String("node_id", nodeid),
			slog.String("old_node_data", oldnode.NodeData),
			slog.String("new_node_data", req.GetNewNodeData()))
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
	if req.GetNodeId()[0] != 0 || req.GetPnodeId()[0] != 0 || (req.GetNodeId()[1] != req.GetPnodeId()[1]) {
		//can't cross project
		return nil, ecode.ErrReq
	}
	if req.GetNodeId()[1] == 1 || req.GetPnodeId()[1] == 1 {
		//0,1 -> admin project
		//can't modify
		return nil, ecode.ErrPermission
	}
	if req.GetNodeId()[2] == 1 || req.GetNodeId()[2] == 2 {
		//0,x,1 -> UserAndRoleControl
		//0,x,2 -> AppControl
		//these are default,can't modify
		return nil, ecode.ErrPermission
	}
	if len(req.GetPnodeId()) >= 3 && (req.GetPnodeId()[2] == 1 || req.GetPnodeId()[2] == 2) {
		//0,x,1 -> UserAndRoleControl
		//0,x,2 -> AppControl
		//these are default,can't modify
		return nil, ecode.ErrPermission
	}
	if len(req.GetPnodeId())+1 == len(req.GetNodeId()) {
		//0,x,y,z move to 0,x,y is equal to not move
		child := true
		for i := range req.GetPnodeId() {
			if req.GetPnodeId()[i] != req.GetNodeId()[i] {
				child = false
				break
			}
		}
		if child {
			return &api.MoveNodeResp{}, nil
		}
	}

	nodeid := util.FormID(req.GetNodeId())
	pnodeid := util.FormID(req.GetPnodeId())

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
	if req.GetNodeId()[0] != 0 {
		return nil, ecode.ErrReq
	}
	if req.GetNodeId()[1] == 1 {
		//0,1 -> admin project
		//can't delete
		return nil, ecode.ErrPermission
	}
	if req.GetNodeId()[2] == 1 || req.GetNodeId()[2] == 2 {
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

	nodeid := util.FormID(req.GetNodeId())

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
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}

	projectid := util.FormID(req.GetProjectId())

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[ListUserNode] operator's token format wrong", slog.String("operator", md["Token-User"]))
		return nil, ecode.ErrToken
	}
	var target bson.ObjectID
	if req.GetUserId() == "" || req.GetUserId() == md["Token-User"] {
		//list self's
		req.SetUserId(md["Token-User"])
		target = operator
	} else {
		//list other user's
		target, e = bson.ObjectIDFromHex(req.GetUserId())
		if e != nil {
			slog.ErrorContext(ctx, "[ListUserNode] target's userid format wrong", slog.String("user_id", req.GetUserId()))
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
			slog.String("user_id", req.GetUserId()),
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
			slog.String("user_id", req.GetUserId()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	var usernodes model.UserNodes
	if !target.IsZero() {
		if usernodes, e = s.permissionDao.MongoGetUserNodes(ctx, target, projectid, nil); e != nil {
			slog.ErrorContext(ctx, "[ListUserNode] db op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("user_id", req.GetUserId()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	var userrolenodes map[string]model.RoleNodes
	if !target.IsZero() && req.GetNeedUserRoleNode() {
		if userrolenodes, e = s.permissionDao.MongoGetUserRoleNodes(ctx, target, projectid, nil); e != nil {
			slog.ErrorContext(ctx, "[ListUserNode] db op failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("user_id", req.GetUserId()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
	}
	projectnode := &api.NodeInfo{}
	projectnode.SetNodeId(req.GetProjectId())
	projectnode.SetNodeName(root.NodeName)
	projectnode.SetNodeData(root.NodeData)
	projectnode.SetChildren(make([]*api.NodeInfo, 0, 10))
	if target.IsZero() {
		projectnode.SetCanread(true)
		projectnode.SetCanwrite(true)
		projectnode.SetAdmin(true)
	} else {
		canread, canwrite, admin := usernodes.CheckNode(projectid)
		projectnode.SetCanread(canread)
		projectnode.SetCanwrite(canwrite)
		projectnode.SetAdmin(admin)
		for _, rolenodes := range userrolenodes {
			if projectnode.GetAdmin() {
				break
			}
			tmpr, tmpw, tmpa := rolenodes.CheckNode(projectid)
			if tmpr {
				projectnode.SetCanread(tmpr)
			}
			if tmpw {
				projectnode.SetCanwrite(tmpw)
			}
			if tmpa {
				projectnode.SetAdmin(tmpa)
			}
		}
	}
	sort.Slice(children, func(i, j int) bool {
		return strings.Count(children[i].NodeId, ",") < strings.Count(children[j].NodeId, ",")
	})
	for _, node := range children {
		nodeid, e := util.ParseID(node.NodeId)
		if e != nil {
			slog.ErrorContext(ctx, "[ListUserNode] target's node's nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("user_id", req.GetUserId()),
				slog.String("node_id", node.NodeId))
			return nil, ecode.ErrSystem
		}
		tmp := &api.NodeInfo{}
		tmp.SetNodeId(nodeid)
		tmp.SetNodeName(node.NodeName)
		tmp.SetNodeData(node.NodeData)
		tmp.SetChildren(make([]*api.NodeInfo, 0, 10))
		if target.IsZero() {
			tmp.SetCanread(true)
			tmp.SetCanwrite(true)
			tmp.SetAdmin(true)
		} else {
			canread, canwrite, admin := usernodes.CheckNode(node.NodeId)
			tmp.SetCanread(canread)
			tmp.SetCanwrite(canwrite)
			tmp.SetAdmin(admin)
			for _, rolenodes := range userrolenodes {
				if tmp.GetAdmin() {
					break
				}
				tmpr, tmpw, tmpa := rolenodes.CheckNode(node.NodeId)
				if tmpr {
					tmp.SetCanread(tmpr)
				}
				if tmpw {
					tmp.SetCanwrite(tmpw)
				}
				if tmpa {
					tmp.SetAdmin(tmpa)
				}
			}
		}
		addTreeNode(projectnode, tmp)
	}
	sortTreeNodes(projectnode.GetChildren())
	resp := &api.ListUserNodeResp{}
	resp.SetNode(projectnode)
	return resp, nil
}
func (s *Service) ListRoleNode(ctx context.Context, req *api.ListRoleNodeReq) (*api.ListRoleNodeResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}

	projectid := util.FormID(req.GetProjectId())

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
			slog.String("role_name", req.GetRoleName()),
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
			slog.String("role_name", req.GetRoleName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	rolenodes, e := s.permissionDao.MongoGetRoleNodes(ctx, projectid, req.GetRoleName(), nil)
	if e != nil {
		slog.ErrorContext(ctx, "[ListRoleNode] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("role_name", req.GetRoleName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	projectnode := &api.NodeInfo{}
	projectnode.SetNodeId(req.GetProjectId())
	projectnode.SetNodeName(root.NodeName)
	projectnode.SetNodeData(root.NodeData)
	projectnode.SetChildren(make([]*api.NodeInfo, 0, 10))
	canread, canwrite, admin := rolenodes.CheckNode(projectid)
	projectnode.SetCanread(canread)
	projectnode.SetCanwrite(canwrite)
	projectnode.SetAdmin(admin)
	sort.Slice(children, func(i, j int) bool {
		return strings.Count(children[i].NodeId, ",") < strings.Count(children[j].NodeId, ",")
	})
	for _, node := range children {
		nodeid, e := util.ParseID(node.NodeId)
		if e != nil {
			slog.ErrorContext(ctx, "[ListRoleNode] role's node's nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("role_name", req.GetRoleName()),
				slog.String("node_id", node.NodeId))
			return nil, ecode.ErrSystem
		}
		tmp := &api.NodeInfo{}
		tmp.SetNodeId(nodeid)
		tmp.SetNodeName(node.NodeName)
		tmp.SetNodeData(node.NodeData)
		tmp.SetChildren(make([]*api.NodeInfo, 0, 10))
		canread, canwrite, admin := rolenodes.CheckNode(node.NodeId)
		tmp.SetCanread(canread)
		tmp.SetCanwrite(canwrite)
		tmp.SetAdmin(admin)
		addTreeNode(projectnode, tmp)
	}
	sortTreeNodes(projectnode.GetChildren())
	resp := &api.ListRoleNodeResp{}
	resp.SetNode(projectnode)
	return resp, nil
}
func (s *Service) ListProjectNode(ctx context.Context, req *api.ListProjectNodeReq) (*api.ListProjectNodeResp, error) {
	if req.GetProjectId()[0] != 0 {
		return nil, ecode.ErrReq
	}

	projectid := util.FormID(req.GetProjectId())

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
	projectnode := &api.NodeInfo{}
	projectnode.SetNodeId(req.GetProjectId())
	projectnode.SetNodeName(root.NodeName)
	projectnode.SetNodeData(root.NodeData)
	projectnode.SetChildren(make([]*api.NodeInfo, 0, 10))
	sort.Slice(children, func(i, j int) bool {
		return strings.Count(children[i].NodeId, ",") < strings.Count(children[j].NodeId, ",")
	})
	for _, node := range children {
		nodeid, e := util.ParseID(node.NodeId)
		if e != nil {
			slog.ErrorContext(ctx, "[ListProjectNode] project's node's nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("node_id", node.NodeId))
			return nil, ecode.ErrSystem
		}
		tmp := &api.NodeInfo{}
		tmp.SetNodeId(nodeid)
		tmp.SetNodeName(node.NodeName)
		tmp.SetNodeData(node.NodeData)
		tmp.SetChildren(make([]*api.NodeInfo, 0, 10))
		addTreeNode(projectnode, tmp)
	}
	sortTreeNodes(projectnode.GetChildren())
	resp := &api.ListProjectNodeResp{}
	resp.SetNode(projectnode)
	return resp, nil
}

func addTreeNode(root, node *api.NodeInfo) bool {
	if len(root.GetNodeId()) > len(node.GetNodeId()) {
		return false
	}
	isprefix := true
	for i := range root.GetNodeId() {
		if root.GetNodeId()[i] != node.GetNodeId()[i] {
			isprefix = false
			break
		}
	}
	if !isprefix {
		return false
	}
	if len(root.GetNodeId()) == len(node.GetNodeId()) {
		return true
	}
	for _, child := range root.GetChildren() {
		if addTreeNode(child, node) {
			return true
		}
	}
	root.SetChildren(append(root.GetChildren(), node))
	return true
}
func sortTreeNodes(nodes []*api.NodeInfo) {
	sort.Slice(nodes, func(i, j int) bool {
		if len(nodes[i].GetNodeId()) < len(nodes[j].GetNodeId()) {
			return true
		} else if len(nodes[i].GetNodeId()) > len(nodes[j].GetNodeId()) {
			return false
		}
		for k := range len(nodes[i].GetNodeId()) {
			if nodes[i].GetNodeId()[k] < nodes[j].GetNodeId()[k] {
				return true
			}
		}
		return false
	})
	for _, node := range nodes {
		if len(node.GetChildren()) > 1 {
			sortTreeNodes(node.GetChildren())
		}
	}
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
