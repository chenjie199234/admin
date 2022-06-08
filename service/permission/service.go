package permission

import (
	"context"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"

	cerror "github.com/chenjie199234/Corelib/error"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
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
	if req.PnodeId[0] != 0 {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	userid := md["Token-Data"]
	obj, e := primitive.ObjectIDFromHex(userid)
	if e != nil {
		log.Error(ctx, "[ListUserNode] userid:", userid, "format error:", e)
		return nil, ecode.ErrAuth
	}
	r, _, x, e := s.permissionDao.MongoGetUserPermission(ctx, obj, req.PnodeId)
	if e != nil {
		log.Error(ctx, "[ListUserNode] userid:", userid, "pnodeid:", req.PnodeId, "get user permission error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	if !r && !x {
		log.Error(ctx, "[ListUserNode] userid:", userid, "pnodeid:", req.PnodeId, "missing permission")
		return nil, ecode.ErrPermission
	}
	nodes, e := s.permissionDao.MongoListNode(ctx, obj, req.PnodeId)
	if e != nil {
		log.Error(ctx, "[ListUserNode] userid:", userid, "pnodeid:", req.PnodeId, "error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	resp := &api.ListUserNodeResp{
		Nodes: make([]*api.NodeInfo, 0, len(nodes)),
	}
	for _, node := range nodes {
		resp.Nodes = append(resp.Nodes, &api.NodeInfo{
			NodeId:   node.NodeId,
			NodeName: node.NodeName,
			NodeData: node.NodeData,
		})
	}
	return resp, nil
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
	r, _, x, e := s.permissionDao.MongoGetUserPermission(ctx, obj, req.NodeId)
	if e != nil {
		log.Error(ctx, "[ListNodeUser] userid:", userid, "nodeid:", req.NodeId, "get user permission error:", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	if !x && !r {
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

//Stop -
func (s *Service) Stop() {

}
