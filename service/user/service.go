package user

import (
	"context"
	"time"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	userdao "github.com/chenjie199234/admin/dao/user"
	"github.com/chenjie199234/admin/ecode"

	cerror "github.com/chenjie199234/Corelib/error"
	"github.com/chenjie199234/Corelib/log"
	publicmids "github.com/chenjie199234/Corelib/mids"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/chenjie199234/Corelib/web"
	//"github.com/chenjie199234/Corelib/cgrpc"
	//"github.com/chenjie199234/Corelib/crpc"
)

//Service subservice for user business
type Service struct {
	userDao *userdao.Dao
}

//Start -
func Start() *Service {
	return &Service{
		userDao: userdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}
}
func (s *Service) SuperAdminLogin(ctx context.Context, req *api.SuperAdminLoginReq) (*api.SuperAdminLoginResp, error) {
	users, e := s.userDao.MongoGetUsers(ctx, []primitive.ObjectID{primitive.NilObjectID})
	if e != nil {
		log.Error(ctx, "[SuperAdminLogin]", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	user, ok := users[primitive.NilObjectID]
	if !ok {
		return nil, ecode.ErrNotInited
	}
	if user.Password != req.Password {
		return nil, ecode.ErrAuth
	}
	start := time.Now()
	end := start.Add(config.AC.TokenExpire.StdDuration())
	tokenstr := publicmids.MakeToken(config.AC.TokenSecret, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, user.ID.Hex(), uint64(start.Unix()), uint64(end.Unix()))
	return &api.SuperAdminLoginResp{Token: tokenstr}, nil
}
func (s *Service) Login(ctx context.Context, req *api.LoginReq) (*api.LoginResp, error) {
	var userid string
	//TODO get userid
	start := time.Now()
	end := start.Add(config.AC.TokenExpire.StdDuration())
	tokenstr := publicmids.MakeToken(config.AC.TokenSecret, "corelib", *config.EC.DeployEnv, *config.EC.RunEnv, userid, uint64(start.Unix()), uint64(end.Unix()))
	return &api.LoginResp{Token: tokenstr}, nil
}
func (s *Service) GetUsers(ctx context.Context, req *api.GetUsersReq) (*api.GetUsersResp, error) {
	undup := make(map[primitive.ObjectID]*struct{})
	for _, userid := range req.UserIds {
		obj, e := primitive.ObjectIDFromHex(userid)
		if e != nil {
			log.Error(ctx, "[GetUsers] userid:", userid, "format error:", e)
			return nil, ecode.ErrReq
		}
		undup[obj] = nil
	}
	userids := make([]primitive.ObjectID, 0, len(undup))
	for userid := range undup {
		userids = append(userids, userid)
	}
	users, e := s.userDao.MongoGetUsers(ctx, userids)
	if e != nil {
		log.Error(ctx, "[GetUsers]", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	resp := &api.GetUsersResp{
		Users: make([]*api.UserInfo, 0, len(users)),
	}
	for _, user := range users {
		resp.Users = append(resp.Users, &api.UserInfo{
			UserId:     user.ID.Hex(),
			UserName:   user.Name,
			Department: user.Department,
			Ctime:      user.Ctime,
		})
	}
	return resp, nil
}
func (s *Service) SearchUsers(ctx context.Context, req *api.SearchUsersReq) (*api.SearchUsersResp, error) {
	users, e := s.userDao.MongoSearchUsers(ctx, req.UserName, 10)
	if e != nil {
		log.Error(ctx, "[SearchUser]", e)
		if _, ok := e.(*cerror.Error); ok {
			return nil, e
		}
		return nil, ecode.ErrSystem
	}
	resp := &api.SearchUsersResp{
		Users: make([]*api.UserInfo, 0, len(users)),
	}
	for _, user := range users {
		resp.Users = append(resp.Users, &api.UserInfo{
			UserId:     user.ID.Hex(),
			UserName:   user.Name,
			Department: user.Department,
			Ctime:      user.Ctime,
		})
	}
	return resp, nil
}

//Stop -
func (s *Service) Stop() {

}
