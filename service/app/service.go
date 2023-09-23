package app

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	appdao "github.com/chenjie199234/admin/dao/app"
	initializedao "github.com/chenjie199234/admin/dao/initialize"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/pool"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/graceful"
	"github.com/chenjie199234/Corelib/util/name"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "github.com/chenjie199234/Corelib/web"
	// "github.com/chenjie199234/Corelib/cgrpc"
)

// Service subservice for config business
type Service struct {
	stop *graceful.Graceful

	appDao        *appdao.Dao
	permissionDao *permissiondao.Dao
	initializeDao *initializedao.Dao
}

// Start -
func Start() *Service {
	return &Service{
		stop: graceful.New(),

		appDao:        appdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		initializeDao: initializedao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}
}

func (s *Service) GetApp(ctx context.Context, req *api.GetAppReq) (*api.GetAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[GetApp] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[GetApp] get app's permission nodeid failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[GetApp] app's permission nodeid format wrong",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[GetApp] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("nodeid", nodeid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	app, e := s.appDao.MongoGetApp(ctx, projectid, req.GName, req.AName, req.Secret)
	if e != nil {
		log.Error(ctx, "[GetApp] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.GetAppResp{
		DiscoverMode:            app.DiscoverMode,
		KubernetesNamespace:     app.KubernetesNs,
		KubernetesLabelselector: app.KubernetesLS,
		DnsHost:                 app.DnsHost,
		DnsInterval:             app.DnsInterval,
		StaticAddrs:             app.StaticAddrs,
		Keys:                    make(map[string]*api.KeyConfigInfo),
		Paths:                   make(map[string]*api.ProxyPathInfo),
	}
	for k, v := range app.Keys {
		resp.Keys[k] = &api.KeyConfigInfo{
			CurIndex:     v.CurIndex,
			MaxIndex:     v.MaxIndex,
			CurVersion:   v.CurVersion,
			CurValue:     v.CurValue,
			CurValueType: v.CurValueType,
		}
	}
	for k, v := range app.Paths {
		nodeid, e := util.ParseNodeIDstr(v.PermissionNodeID)
		if e != nil {
			log.Error(ctx, "[GetApp] app's path's permission nodeid format wrong",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.String("path", k),
				log.String("nodeid", v.PermissionNodeID),
				log.CError(e))
			return nil, ecode.ErrDBDataBroken
		}
		resp.Paths[k] = &api.ProxyPathInfo{
			NodeId: nodeid,
			Read:   v.PermissionRead,
			Write:  v.PermissionWrite,
			Admin:  v.PermissionAdmin,
		}
	}
	return resp, nil
}

func (s *Service) SetApp(ctx context.Context, req *api.SetAppReq) (*api.SetAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[SetApp] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}
	if e := name.SingleCheck(req.GName, false); e != nil {
		log.Error(ctx, "[SetApp] group name format wrong", log.String("operator", md["Token-User"]), log.String("group", req.GName))
		return nil, ecode.ErrReq
	}
	if e := name.SingleCheck(req.AName, false); e != nil {
		log.Error(ctx, "[SetApp] app name format wrong", log.String("operator", md["Token-User"]), log.String("app", req.AName))
		return nil, ecode.ErrReq
	}
	switch req.DiscoverMode {
	case "kubernetes":
		if req.KubernetesNamespace == "" {
			log.Error(ctx, "[SetApp] kubernetes namesapce empty", log.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
		if req.KubernetesLabelselector == "" {
			log.Error(ctx, "[SetApp] kubernetes labelselector empty", log.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
	case "dns":
		if req.DnsHost == "" {
			log.Error(ctx, "[SetApp] dns host empty", log.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
		if req.DnsInterval == 0 {
			log.Error(ctx, "[SetApp] dns interval must be set", log.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
	case "static":
		if len(req.StaticAddrs) == 0 {
			log.Error(ctx, "[SetApp] static addrs empty", log.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	if !operator.IsZero() {
		//config control permission check
		if req.NewApp {
			//create new app need the AppControl's admin permission
			_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.AppControl, true)
			if e != nil {
				log.Error(ctx, "[SetApp] get operator's permission info failed",
					log.String("operator", md["Token-User"]),
					log.String("nodeid", projectid+model.AppControl),
					log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if !admin {
				return nil, ecode.ErrPermission
			}
		} else {
			//update app need the app's admin permission
			nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
			if e != nil {
				log.Error(ctx, "[SetApp] get app's permission nodeid failed",
					log.String("operator", md["Token-User"]),
					log.String("project_id", projectid),
					log.String("group", req.GName),
					log.String("app", req.AName),
					log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			nodeids := strings.Split(nodeid, ",")
			if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
				log.Error(ctx, "[SetApp] app's permission nodeid format wrong",
					log.String("operator", md["Token-User"]),
					log.String("project_id", projectid),
					log.String("group", req.GName),
					log.String("app", req.AName),
					log.String("nodeid", nodeid))
				return nil, ecode.ErrDBDataBroken
			}
			_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
			if e != nil {
				log.Error(ctx, "[SetApp] get operator's permission info failed",
					log.String("operator", md["Token-User"]),
					log.String("nodeid", nodeid),
					log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if !admin {
				return nil, ecode.ErrPermission
			}
		}
	}

	//logic
	var nodeidstr string
	if req.NewApp {
		nodeidstr, e = s.appDao.MongoCreateApp(ctx, projectid, req.GName, req.AName, req.Secret, req.DiscoverMode, req.KubernetesNamespace, req.KubernetesLabelselector, req.DnsHost, req.DnsInterval, req.StaticAddrs)
	} else {
		nodeidstr, e = s.appDao.MongoUpdateApp(ctx, projectid, req.GName, req.AName, req.Secret, req.DiscoverMode, req.KubernetesNamespace, req.KubernetesLabelselector, req.DnsHost, req.DnsInterval, req.StaticAddrs)
	}
	if e != nil {
		log.Error(ctx, "[SetApp] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeid, e := util.ParseNodeIDstr(nodeidstr)
	if e != nil {
		log.Error(ctx, "[SetApp] nodeid format wrong",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.CError(ecode.ErrDBDataBroken))
	}
	log.Info(ctx, "[SetApp] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("group", req.GName),
		log.String("app", req.AName))
	return &api.SetAppResp{NodeId: nodeid}, nil
}

func (s *Service) DelApp(ctx context.Context, req *api.DelAppReq) (*api.DelAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[DelApp] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
	if e != nil {
		log.Error(ctx, "[DelApp] get app's permission nodeid failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		log.Error(ctx, "[DelApp] app's permission nodeid format wrong",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("nodeid", nodeid))
		return nil, ecode.ErrDBDataBroken
	}
	//self can't be deleted
	if nodeids[1] == "1" && nodeids[3] == "1" {
		log.Error(ctx, "[DelApp] can't delete self",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName))
		return nil, ecode.ErrPermission
	}

	if !operator.IsZero() {
		//config control permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.AppControl, true)
		if e != nil {
			log.Error(ctx, "[DelApp] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("nodeid", projectid+model.AppControl),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoDelApp(ctx, projectid, req.GName, req.AName, req.Secret); e != nil {
		log.Error(ctx, "[DelApp] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelApp] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("group", req.GName),
		log.String("app", req.AName))
	return &api.DelAppResp{}, nil
}

func (s *Service) UpdateAppSecret(ctx context.Context, req *api.UpdateAppSecretReq) (*api.UpdateAppSecretResp, error) {
	if req.OldSecret == req.NewSecret {
		return &api.UpdateAppSecretResp{}, nil
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[UpdateAppSecret] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[UpdateAppSecret] get app's permission nodeid failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[UpdateAppSecret] app's permission nodeid format wrong",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[UpdateAppSecret] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("nodeid", nodeid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoUpdateAppSecret(ctx, projectid, req.GName, req.AName, req.OldSecret, req.NewSecret); e != nil {
		log.Error(ctx, "[UpdateAppSecret] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[UpdateAppSecret] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("group", req.GName),
		log.String("app", req.AName))
	return &api.UpdateAppSecretResp{}, nil
}

func (s *Service) DelKey(ctx context.Context, req *api.DelKeyReq) (*api.DelKeyResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[DelKey] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
	if e != nil {
		log.Error(ctx, "[DelKey] get app's permission nodeid failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		log.Error(ctx, "[DelKey] app's permission nodeid format wrong",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("nodeid", nodeid))
		return nil, ecode.ErrDBDataBroken
	}
	if nodeids[1] == "1" && nodeids[3] == "1" && (req.Key == "AppConfig" || req.Key == "SourceConfig") {
		log.Error(ctx, "[DelKey] can't delete self's 'AppConfig' or 'SourceConfig' key",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName))
		return nil, ecode.ErrPermission
	}

	if !operator.IsZero() {
		//config control permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[DelKey] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("nodeid", nodeid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	if e := s.appDao.MongoDelKey(ctx, projectid, req.GName, req.AName, req.Key, req.Secret); e != nil {
		log.Error(ctx, "[DelKey] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("key", req.Key),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelKey] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("group", req.GName),
		log.String("app", req.AName),
		log.String("key", req.Key))
	return &api.DelKeyResp{}, nil
}

func (s *Service) GetKeyConfig(ctx context.Context, req *api.GetKeyConfigReq) (*api.GetKeyConfigResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[GetKeyConfig] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[GetKeyConfig] get app's permission nodeid failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[GetKeyConfig] app's permission nodeid format wrong",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[GetKeyConfig] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("nodeid", nodeid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	keysummary, configlog, e := s.appDao.MongoGetKeyConfig(ctx, projectid, req.GName, req.AName, req.Key, req.Index, req.Secret)
	if e != nil {
		log.Error(ctx, "[GetKeyConfig] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("key", req.Key),
			log.Uint64("index", uint64(req.Index)),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.GetKeyConfigResp{
		CurIndex:   keysummary.CurIndex,
		MaxIndex:   keysummary.MaxIndex,
		CurVersion: keysummary.CurVersion,
		ThisIndex:  configlog.Index,
		Value:      configlog.Value,
		ValueType:  configlog.ValueType,
	}, nil
}

func (s *Service) SetKeyConfig(ctx context.Context, req *api.SetKeyConfigReq) (*api.SetKeyConfigResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[SetKeyConfig] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	req.Key = strings.TrimSpace(req.Key)
	if req.Key == "" {
		log.Error(ctx, "[SetKeyConfig] key empty",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName))
		return nil, ecode.ErrReq
	}
	req.Value = strings.TrimSpace(req.Value)
	if req.Value == "" {
		log.Error(ctx, "[SetKeyConfig] value empty",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("key", req.Key))
		return nil, ecode.ErrReq
	}
	switch req.ValueType {
	case "json":
		buf := bytes.NewBuffer(nil)
		if e := json.Compact(buf, common.Str2byte(req.Value)); e != nil {
			log.Error(ctx, "[SetKeyConfig] json value format wrong",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.String("key", req.Key),
				log.CError(e))
			return nil, ecode.ErrReq
		}
		req.Value = common.Byte2str(buf.Bytes())
	case "toml":
		//TODO
		fallthrough
	case "yaml":
		//TODO
		fallthrough
	case "raw":
		//TODO
		fallthrough
	default:
		log.Error(ctx, "[SetKeyConfig] unsupported value type",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("key", req.Key),
			log.String("valuetype", req.ValueType))
		return nil, ecode.ErrReq
	}

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[SetKeyConfig] get app's permission nodeid failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[SetKeyConfig] app's permission nodeid format wrong",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[SetKeyConfig] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("nodeid", nodeid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	index, version, e := s.appDao.MongoSetKeyConfig(ctx, projectid, req.GName, req.AName, req.Key, req.Secret, req.Value, req.ValueType, req.NewKey)
	if e != nil {
		log.Error(ctx, "[SetKeyConfig] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("key", req.Key),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[SetKeyConfig] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("group", req.GName),
		log.String("app", req.AName),
		log.String("key", req.Key),
		log.Uint64("new_index", uint64(index)),
		log.Uint64("new_version", uint64(version)))
	return &api.SetKeyConfigResp{}, nil
}

func (s *Service) Rollback(ctx context.Context, req *api.RollbackReq) (*api.RollbackResp, error) {
	if strings.Contains(req.Key, ".") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[Rollback] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[Rollback] get app's permission nodeid failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[Rollback] app's permission nodeid format wrong",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[Rollback] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("nodeid", nodeid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoRollbackKeyConfig(ctx, projectid, req.GName, req.AName, req.Key, req.Secret, req.Index); e != nil {
		log.Error(ctx, "[Rollback] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("key", req.Key),
			log.Uint64("index", uint64(req.Index)),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[Rollback] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("group", req.GName),
		log.String("app", req.AName),
		log.String("key", req.Key),
		log.Uint64("index", uint64(req.Index)))
	return &api.RollbackResp{}, nil
}

func (s *Service) Watch(ctx context.Context, req *api.WatchReq) (*api.WatchResp, error) {
	for k := range req.Keys {
		if strings.Contains(k, ".") {
			return nil, ecode.ErrReq
		}
	}
	if e := s.stop.Add(1); e != nil {
		if e == graceful.ErrClosing {
			return nil, ecode.ErrServerClosing
		}
		return nil, ecode.ErrBusy
	}
	defer s.stop.DoneOne()
	ch, cancel, e := config.Sdk.GetNoticeByProjectName(req.ProjectName, req.GName, req.AName)
	if e != nil {
		log.Error(ctx, "[Watch] get notice failed",
			log.String("project_name", req.ProjectName),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil, ecode.ReturnEcode(ctx.Err(), ecode.ErrSystem)
		case <-ch:
			app, e := config.Sdk.GetAppConfigByProjectName(req.ProjectName, req.GName, req.AName)
			if e != nil {
				log.Error(ctx, "[Watch] get config failed",
					log.String("project_name", req.ProjectName),
					log.String("group", req.GName),
					log.String("app", req.AName),
					log.CError(e))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			resp := &api.WatchResp{
				Datas: make(map[string]*api.WatchData, len(req.Keys)+3),
			}
			needreturn := false
			for key, clientversion := range req.Keys {
				k, ok := app.Keys[key]
				if !ok || k == nil || k.CurVersion == 0 {
					return nil, ecode.ErrKeyNotExist
				}
				if clientversion != k.CurVersion {
					needreturn = true
					resp.Datas[key] = &api.WatchData{
						Key:       key,
						Value:     k.CurValue,
						ValueType: k.CurValueType,
						Version:   k.CurVersion,
					}
				} else {
					resp.Datas[key] = &api.WatchData{
						Key:       key,
						Value:     "",
						ValueType: "",
						Version:   k.CurVersion,
					}
				}
			}
			if needreturn {
				return resp, nil
			}
		}
	}
}

func (s *Service) SetProxy(ctx context.Context, req *api.SetProxyReq) (*api.SetProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[SetProxy] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[SetProxy] get app's permission nodeid failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[SetProxy] app's permission nodeid format wrong",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[SetProxy] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("nodeid", nodeid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if req.Path[0] != '/' {
		req.Path = "/" + req.Path
	}
	newnodeidstr, e := s.appDao.MongoSetProxyPath(ctx, projectid, req.GName, req.AName, req.Secret, req.Path, req.Read, req.Write, req.Admin, req.NewPath)
	if e != nil {
		log.Error(ctx, "[SetProxy] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("path", req.Path),
			log.Any("read_write_admin", []bool{req.Read, req.Write, req.Admin}),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	newnodeid, _ := util.ParseNodeIDstr(newnodeidstr)
	log.Info(ctx, "[SetProxy] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("group", req.GName),
		log.String("app", req.AName),
		log.String("path", req.Path),
		log.Any("read_write_admin", []bool{req.Read, req.Write, req.Admin}))
	return &api.SetProxyResp{NodeId: newnodeid}, nil
}
func (s *Service) DelProxy(ctx context.Context, req *api.DelProxyReq) (*api.DelProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[DelProxy] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			log.Error(ctx, "[DelProxy] get app's permission nodeid failed",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			log.Error(ctx, "[DelProxy] app's permission nodeid format wrong",
				log.String("operator", md["Token-User"]),
				log.String("project_id", projectid),
				log.String("group", req.GName),
				log.String("app", req.AName),
				log.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[DelProxy] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("nodeid", nodeid),
				log.CError(e))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoDelProxyPath(ctx, projectid, req.GName, req.AName, req.Secret, req.Path); e != nil {
		log.Error(ctx, "[DelProxy] db op failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("path", req.Path),
			log.CError(e))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	log.Info(ctx, "[DelProxy] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("group", req.GName),
		log.String("app", req.AName),
		log.String("path", req.Path))
	return &api.DelProxyResp{}, nil
}
func (s *Service) Proxy(ctx context.Context, req *api.ProxyReq) (*api.ProxyResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := primitive.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		log.Error(ctx, "[Proxy] operator's token format wrong", log.String("operator", md["Token-User"]), log.CError(e))
		return nil, ecode.ErrToken
	}

	buf := pool.GetPool().Get(0)
	defer pool.GetPool().Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.Byte2str(buf)

	if req.Path[0] != '/' {
		req.Path = "/" + req.Path
	}
	pcheck := func(cctx context.Context, nodeid string, read, write, admin bool) error {
		if operator.IsZero() || (!read && !write && !admin) {
			return nil
		}
		//permission check
		canread, canwrite, canadmin, e := s.permissionDao.MongoGetUserPermission(cctx, operator, nodeid, true)
		if e != nil {
			log.Error(ctx, "[Proxy] get operator's permission info failed",
				log.String("operator", md["Token-User"]),
				log.String("nodeid", nodeid),
				log.CError(e))
			return ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if admin && !canadmin {
			return ecode.ErrPermission
		} else if write && !canwrite && !canadmin {
			return ecode.ErrPermission
		} else if read && !canread && !canadmin {
			return ecode.ErrPermission
		}
		return nil
	}
	out, e := config.Sdk.CallByPrjoectID(ctx, projectid, req.GName, req.AName, req.Path, common.Str2byte(req.Data), pcheck)
	if e != nil {
		log.Error(ctx, "[Proxy] call server failed",
			log.String("operator", md["Token-User"]),
			log.String("project_id", projectid),
			log.String("group", req.GName),
			log.String("app", req.AName),
			log.String("path", req.Path),
			log.String("reqdata", req.Data),
			log.CError(e))
		return nil, e
	}
	log.Info(ctx, "[Proxy] success",
		log.String("operator", md["Token-User"]),
		log.String("project_id", projectid),
		log.String("group", req.GName),
		log.String("app", req.AName),
		log.String("path", req.Path),
		log.String("reqdata", req.Data),
		log.String("respdata", common.Byte2str(out)))
	return &api.ProxyResp{Data: common.Byte2str(out)}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
