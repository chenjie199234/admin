package app

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"slices"
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

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/pool/bpool"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/egroup"
	"github.com/chenjie199234/Corelib/util/graceful"
	"github.com/chenjie199234/Corelib/util/name"
	"go.mongodb.org/mongo-driver/v2/bson"
	// "github.com/chenjie199234/Corelib/web"
	// "github.com/chenjie199234/Corelib/crpc"
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
func Start() (*Service, error) {
	return &Service{
		stop: graceful.New(),

		appDao:        appdao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		permissionDao: permissiondao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
		initializeDao: initializedao.NewDao(nil, nil, config.GetMongo("admin_mongo")),
	}, nil
}

func (s *Service) GetApp(ctx context.Context, req *api.GetAppReq) (*api.GetAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[GetApp] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
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

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			slog.ErrorContext(ctx, "[GetApp] get app's permission nodeid failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			slog.ErrorContext(ctx, "[GetApp] app's permission nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			slog.ErrorContext(ctx, "[GetApp] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("nodeid", nodeid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	app, e := s.appDao.MongoGetApp(ctx, projectid, req.GName, req.AName, req.Secret)
	if e != nil {
		slog.ErrorContext(ctx, "[GetApp] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.GetAppResp{
		DiscoverMode:            app.DiscoverMode,
		KubernetesNamespace:     app.KubernetesNs,
		KubernetesLabelselector: app.KubernetesLS,
		KubernetesFieldselector: app.KubernetesFS,
		DnsHost:                 app.DnsHost,
		DnsInterval:             app.DnsInterval,
		StaticAddrs:             app.StaticAddrs,
		CrpcPort:                app.CrpcPort,
		CgrpcPort:               app.CGrpcPort,
		WebPort:                 app.WebPort,
		Keys:                    make(map[string]*api.KeyConfigInfo),
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
	return resp, nil
}

func (s *Service) SetApp(ctx context.Context, req *api.SetAppReq) (*api.SetAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[SetApp] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	if e := name.SingleCheck(req.GName, false); e != nil {
		slog.ErrorContext(ctx, "[SetApp] group name format wrong", slog.String("operator", md["Token-User"]), slog.String("group", req.GName))
		return nil, ecode.ErrReq
	}
	if e := name.SingleCheck(req.AName, false); e != nil {
		slog.ErrorContext(ctx, "[SetApp] app name format wrong", slog.String("operator", md["Token-User"]), slog.String("app", req.AName))
		return nil, ecode.ErrReq
	}
	switch req.DiscoverMode {
	case "kubernetes":
		if req.KubernetesNamespace == "" {
			slog.ErrorContext(ctx, "[SetApp] kubernetes namesapce empty", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
		if req.KubernetesLabelselector == "" && req.KubernetesFieldselector == "" {
			slog.ErrorContext(ctx, "[SetApp] kubernetes labelselector and fieldselector empty", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
	case "dns":
		if req.DnsHost == "" {
			slog.ErrorContext(ctx, "[SetApp] dns host empty", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
		if req.DnsInterval == 0 {
			slog.ErrorContext(ctx, "[SetApp] dns interval must be set", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
	case "static":
		if len(req.StaticAddrs) == 0 {
			slog.ErrorContext(ctx, "[SetApp] static addrs empty", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
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

	if !operator.IsZero() {
		//config control permission check
		if req.NewApp {
			//create new app need the AppControl's admin permission
			_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.AppControl, true)
			if e != nil {
				slog.ErrorContext(ctx, "[SetApp] get operator's permission info failed",
					slog.String("operator", md["Token-User"]),
					slog.String("nodeid", projectid+model.AppControl),
					slog.String("error", e.Error()))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			if !admin {
				return nil, ecode.ErrPermission
			}
		} else {
			//update app need the app's admin permission
			nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
			if e != nil {
				slog.ErrorContext(ctx, "[SetApp] get app's permission nodeid failed",
					slog.String("operator", md["Token-User"]),
					slog.String("project_id", projectid),
					slog.String("group", req.GName),
					slog.String("app", req.AName),
					slog.String("error", e.Error()))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			nodeids := strings.Split(nodeid, ",")
			if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
				slog.ErrorContext(ctx, "[SetApp] app's permission nodeid format wrong",
					slog.String("operator", md["Token-User"]),
					slog.String("project_id", projectid),
					slog.String("group", req.GName),
					slog.String("app", req.AName),
					slog.String("nodeid", nodeid))
				return nil, ecode.ErrDBDataBroken
			}
			_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
			if e != nil {
				slog.ErrorContext(ctx, "[SetApp] get operator's permission info failed",
					slog.String("operator", md["Token-User"]),
					slog.String("nodeid", nodeid),
					slog.String("error", e.Error()))
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
		nodeidstr, e = s.appDao.MongoCreateApp(
			ctx,
			projectid,
			req.GName,
			req.AName,
			req.Secret,
			req.DiscoverMode,
			req.KubernetesNamespace,
			req.KubernetesLabelselector,
			req.KubernetesFieldselector,
			req.DnsHost,
			req.DnsInterval,
			req.StaticAddrs,
			req.CrpcPort,
			req.CgrpcPort,
			req.WebPort)
	} else {
		nodeidstr, e = s.appDao.MongoUpdateApp(
			ctx,
			projectid,
			req.GName,
			req.AName,
			req.Secret,
			req.DiscoverMode,
			req.KubernetesNamespace,
			req.KubernetesLabelselector,
			req.KubernetesFieldselector,
			req.DnsHost,
			req.DnsInterval,
			req.StaticAddrs,
			req.CrpcPort,
			req.CgrpcPort,
			req.WebPort)
	}
	if e != nil {
		slog.ErrorContext(ctx, "[SetApp] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeid, e := util.ParseNodeIDstr(nodeidstr)
	if e != nil {
		slog.ErrorContext(ctx, "[SetApp] nodeid format wrong",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", ecode.ErrDBDataBroken.Error()))
	}
	slog.InfoContext(ctx, "[SetApp] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GName),
		slog.String("app", req.AName))
	return &api.SetAppResp{NodeId: nodeid}, nil
}

func (s *Service) DelApp(ctx context.Context, req *api.DelAppReq) (*api.DelAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[DelApp] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
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

	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
	if e != nil {
		slog.ErrorContext(ctx, "[DelApp] get app's permission nodeid failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		slog.ErrorContext(ctx, "[DelApp] app's permission nodeid format wrong",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("nodeid", nodeid))
		return nil, ecode.ErrDBDataBroken
	}
	//self can't be deleted
	if nodeids[1] == "1" && nodeids[3] == "1" {
		slog.ErrorContext(ctx, "[DelApp] can't delete self",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName))
		return nil, ecode.ErrPermission
	}

	if !operator.IsZero() {
		//config control permission check
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, projectid+model.AppControl, true)
		if e != nil {
			slog.ErrorContext(ctx, "[DelApp] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("nodeid", projectid+model.AppControl),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoDelApp(ctx, projectid, req.GName, req.AName, req.Secret); e != nil {
		slog.ErrorContext(ctx, "[DelApp] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[DelApp] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GName),
		slog.String("app", req.AName))
	return &api.DelAppResp{}, nil
}

func (s *Service) UpdateAppSecret(ctx context.Context, req *api.UpdateAppSecretReq) (*api.UpdateAppSecretResp, error) {
	if req.OldSecret == req.NewSecret {
		return &api.UpdateAppSecretResp{}, nil
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateAppSecret] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
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

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateAppSecret] get app's permission nodeid failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			slog.ErrorContext(ctx, "[UpdateAppSecret] app's permission nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		_, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateAppSecret] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("nodeid", nodeid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoUpdateAppSecret(ctx, projectid, req.GName, req.AName, req.OldSecret, req.NewSecret); e != nil {
		slog.ErrorContext(ctx, "[UpdateAppSecret] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[UpdateAppSecret] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GName),
		slog.String("app", req.AName))
	return &api.UpdateAppSecretResp{}, nil
}

func (s *Service) DelKey(ctx context.Context, req *api.DelKeyReq) (*api.DelKeyResp, error) {
	if strings.Contains(req.Key, ".") || strings.Contains(req.Key, "$") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[DelKey] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
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

	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
	if e != nil {
		slog.ErrorContext(ctx, "[DelKey] get app's permission nodeid failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		slog.ErrorContext(ctx, "[DelKey] app's permission nodeid format wrong",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("nodeid", nodeid))
		return nil, ecode.ErrDBDataBroken
	}
	if nodeids[1] == "1" && nodeids[3] == "1" && (req.Key == "AppConfig" || req.Key == "SourceConfig") {
		slog.ErrorContext(ctx, "[DelKey] can't delete self's 'AppConfig' or 'SourceConfig' key",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName))
		return nil, ecode.ErrPermission
	}

	if !operator.IsZero() {
		//config control permission check
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			slog.ErrorContext(ctx, "[DelKey] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("nodeid", nodeid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	if e := s.appDao.MongoDelKey(ctx, projectid, req.GName, req.AName, req.Key, req.Secret); e != nil {
		slog.ErrorContext(ctx, "[DelKey] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("key", req.Key),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[DelKey] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GName),
		slog.String("app", req.AName),
		slog.String("key", req.Key))
	return &api.DelKeyResp{}, nil
}

func (s *Service) GetKeyConfig(ctx context.Context, req *api.GetKeyConfigReq) (*api.GetKeyConfigResp, error) {
	if strings.Contains(req.Key, ".") || strings.Contains(req.Key, "$") {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[GetKeyConfig] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
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

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			slog.ErrorContext(ctx, "[GetKeyConfig] get app's permission nodeid failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			slog.ErrorContext(ctx, "[GetKeyConfig] app's permission nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		canread, _, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			slog.ErrorContext(ctx, "[GetKeyConfig] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("nodeid", nodeid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canread && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	keysummary, configlog, e := s.appDao.MongoGetKeyConfig(ctx, projectid, req.GName, req.AName, req.Key, req.Index, req.Secret)
	if e != nil {
		slog.ErrorContext(ctx, "[GetKeyConfig] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("key", req.Key),
			slog.Uint64("index", uint64(req.Index)),
			slog.String("error", e.Error()))
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
	if strings.Contains(req.Key, ".") || strings.Contains(req.Key, "$") {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[SetKeyConfig] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
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

	req.Key = strings.TrimSpace(req.Key)
	if req.Key == "" {
		slog.ErrorContext(ctx, "[SetKeyConfig] key empty",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName))
		return nil, ecode.ErrReq
	}
	req.Value = strings.TrimSpace(req.Value)
	if req.Value == "" {
		slog.ErrorContext(ctx, "[SetKeyConfig] value empty",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("key", req.Key))
		return nil, ecode.ErrReq
	}
	switch req.ValueType {
	case "json":
		buf := bytes.NewBuffer(nil)
		if e := json.Compact(buf, common.STB(req.Value)); e != nil {
			slog.ErrorContext(ctx, "[SetKeyConfig] json value format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("key", req.Key),
				slog.String("error", e.Error()))
			return nil, ecode.ErrReq
		}
		req.Value = common.BTS(buf.Bytes())
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
		slog.ErrorContext(ctx, "[SetKeyConfig] unsupported value type",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("key", req.Key),
			slog.String("valuetype", req.ValueType))
		return nil, ecode.ErrReq
	}

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			slog.ErrorContext(ctx, "[SetKeyConfig] get app's permission nodeid failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			slog.ErrorContext(ctx, "[SetKeyConfig] app's permission nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			slog.ErrorContext(ctx, "[SetKeyConfig] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("nodeid", nodeid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	index, version, e := s.appDao.MongoSetKeyConfig(ctx, projectid, req.GName, req.AName, req.Key, req.Secret, req.Value, req.ValueType, req.NewKey)
	if e != nil {
		slog.ErrorContext(ctx, "[SetKeyConfig] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("key", req.Key),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[SetKeyConfig] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GName),
		slog.String("app", req.AName),
		slog.String("key", req.Key),
		slog.Uint64("new_index", uint64(index)),
		slog.Uint64("new_version", uint64(version)))
	return &api.SetKeyConfigResp{}, nil
}

func (s *Service) Rollback(ctx context.Context, req *api.RollbackReq) (*api.RollbackResp, error) {
	if strings.Contains(req.Key, ".") || strings.Contains(req.Key, "$") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[Rollback] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
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

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GName, req.AName)
		if e != nil {
			slog.ErrorContext(ctx, "[Rollback] get app's permission nodeid failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			slog.ErrorContext(ctx, "[Rollback] app's permission nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GName),
				slog.String("app", req.AName),
				slog.String("nodeid", nodeid))
			return nil, ecode.ErrDBDataBroken
		}
		_, canwrite, admin, e := s.permissionDao.MongoGetUserPermission(ctx, operator, nodeid, true)
		if e != nil {
			slog.ErrorContext(ctx, "[Rollback] get operator's permission info failed",
				slog.String("operator", md["Token-User"]),
				slog.String("nodeid", nodeid),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		if !canwrite && !admin {
			return nil, ecode.ErrPermission
		}
	}

	//logic
	if e := s.appDao.MongoRollbackKeyConfig(ctx, projectid, req.GName, req.AName, req.Key, req.Secret, req.Index); e != nil {
		slog.ErrorContext(ctx, "[Rollback] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("key", req.Key),
			slog.Uint64("index", uint64(req.Index)),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[Rollback] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GName),
		slog.String("app", req.AName),
		slog.String("key", req.Key),
		slog.Uint64("index", uint64(req.Index)))
	return &api.RollbackResp{}, nil
}

func (s *Service) WatchConfig(ctx context.Context, req *api.WatchConfigReq) (*api.WatchConfigResp, error) {
	for k := range req.Keys {
		if strings.Contains(k, ".") || strings.Contains(k, "$") {
			return nil, ecode.ErrReq
		}
	}
	ch, cancel, e := config.Sdk.GetNoticeByProjectName(req.ProjectName, req.GName, req.AName)
	if e != nil {
		slog.ErrorContext(ctx, "[WatchConfig] get notice failed",
			slog.String("project_name", req.ProjectName),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil, cerror.Convert(ctx.Err())
		case <-ch:
			app, e := config.Sdk.GetAppConfigByProjectName(req.ProjectName, req.GName, req.AName)
			if e != nil {
				if e != ecode.ErrServerClosing {
					slog.ErrorContext(ctx, "[WatchConfig] get config failed",
						slog.String("project_name", req.ProjectName),
						slog.String("group", req.GName),
						slog.String("app", req.AName),
						slog.String("error", e.Error()))
				}
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			resp := &api.WatchConfigResp{
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

func (s *Service) WatchDiscover(ctx context.Context, req *api.WatchDiscoverReq) (*api.WatchDiscoverResp, error) {
	if req.CurDiscoverMode == "dns" && (req.CurDnsHost == "" || req.CurDnsInterval == 0) {
		return nil, ecode.ErrReq
	}
	if req.CurDiscoverMode == "static" && len(req.CurStaticAddrs) == 0 {
		return nil, ecode.ErrReq
	}
	if req.CurDiscoverMode == "kubernetes" && (req.CurKubernetesNamespace == "" || (req.CurKubernetesFieldselector == "" && req.CurKubernetesLabelselector == "")) {
		return nil, ecode.ErrReq
	}
	ch, cancel, e := config.Sdk.GetNoticeByProjectName(req.ProjectName, req.GName, req.AName)
	if e != nil {
		slog.ErrorContext(ctx, "[WatchConfig] get notice failed",
			slog.String("project_name", req.ProjectName),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil, cerror.Convert(ctx.Err())
		case <-ch:
			app, e := config.Sdk.GetAppConfigByProjectName(req.ProjectName, req.GName, req.AName)
			if e != nil {
				if e != ecode.ErrServerClosing {
					slog.ErrorContext(ctx, "[WatchDiscover] get config failed",
						slog.String("project_name", req.ProjectName),
						slog.String("group", req.GName),
						slog.String("app", req.AName),
						slog.String("error", e.Error()))
				}
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			needreturn := app.DiscoverMode != req.CurDiscoverMode || app.CrpcPort != req.CurCrpcPort || app.WebPort != req.CurWebPort || app.CGrpcPort != req.CurCgrpcPort
			if !needreturn {
				if app.DiscoverMode == "dns" {
					needreturn = app.DnsHost != req.CurDnsHost || app.DnsInterval != req.CurDnsInterval
				} else if app.DiscoverMode == "static" {
					for _, addr := range app.StaticAddrs {
						if !slices.Contains(req.CurStaticAddrs, addr) {
							needreturn = true
						}
					}
					for _, addr := range req.CurStaticAddrs {
						if !slices.Contains(app.StaticAddrs, addr) {
							needreturn = true
						}
					}
				} else if app.DiscoverMode == "kubernetes" {
					needreturn = app.KubernetesNs != req.CurKubernetesNamespace ||
						app.KubernetesFS != req.CurKubernetesFieldselector ||
						app.KubernetesLS != req.CurKubernetesLabelselector
				}
			}
			if needreturn {
				return &api.WatchDiscoverResp{
					DiscoverMode:            app.DiscoverMode,
					DnsHost:                 app.DnsHost,
					DnsInterval:             app.DnsInterval,
					StaticAddrs:             app.StaticAddrs,
					KubernetesNamespace:     app.KubernetesNs,
					KubernetesLabelselector: app.KubernetesLS,
					KubernetesFieldselector: app.KubernetesFS,
					CrpcPort:                app.CrpcPort,
					WebPort:                 app.WebPort,
					CgrpcPort:               app.CGrpcPort,
				}, nil
			}
		}
	}
}

func (s *Service) GetInstances(ctx context.Context, req *api.GetInstancesReq) (*api.GetInstancesResp, error) {
	md := metadata.GetMetadata(ctx)

	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	if e := s.appDao.MongoCheckSecret(ctx, projectid, req.GName, req.AName, req.Secret); e != nil {
		slog.ErrorContext(ctx, "[GetInstances] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}

	addrs, e := config.Sdk.GetAppAddrsByProjectID(ctx, projectid, req.GName, req.AName)
	if e != nil {
		slog.ErrorContext(ctx, "[GetInstances] get addrs failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.GetInstancesResp{
		Instances: make(map[string]*api.InstanceInfo, len(addrs)),
	}
	for _, addr := range addrs {
		resp.Instances[addr] = nil
	}
	if !req.WithInfo {
		return resp, nil
	}
	eg := egroup.GetGroup(ctx)
	for _, v := range addrs {
		addr := v
		eg.Go(func(gctx context.Context) error {
			r, e := config.Sdk.PingByPrjoectID(gctx, projectid, req.GName, req.AName, addr)
			if e != nil {
				slog.ErrorContext(ctx, "[GetInstances] get info failed",
					slog.String("operator", md["Token-User"]),
					slog.String("project_id", projectid),
					slog.String("group", req.GName),
					slog.String("app", req.AName),
					slog.String("addr", addr),
					slog.String("error", e.Error()))
				return nil
			}
			resp.Instances[addr] = &api.InstanceInfo{
				Name:        r.Host,
				TotalMem:    r.TotalMem,
				CurMemUsage: r.CurMemUsage,
				CpuNum:      r.CpuNum,
				CurCpuUsage: r.CurCpuUsage,
			}
			return nil
		})
	}
	egroup.PutGroup(eg)
	return resp, nil
}
func (s *Service) GetInstanceInfo(ctx context.Context, req *api.GetInstanceInfoReq) (*api.GetInstanceInfoResp, error) {
	md := metadata.GetMetadata(ctx)

	buf := bpool.Get(0)
	defer bpool.Put(&buf)
	for i, v := range req.ProjectId {
		if i != 0 {
			buf = append(buf, ',')
		}
		buf = strconv.AppendUint(buf, uint64(v), 10)
	}
	projectid := common.BTS(buf)

	if e := s.appDao.MongoCheckSecret(ctx, projectid, req.GName, req.AName, req.Secret); e != nil {
		slog.ErrorContext(ctx, "[GetInstanceInfo] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	r, e := config.Sdk.PingByPrjoectID(ctx, projectid, req.GName, req.AName, req.Addr)
	if e != nil {
		slog.ErrorContext(ctx, "[GetInstanceInfo] get info failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GName),
			slog.String("app", req.AName),
			slog.String("addr", req.Addr),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	return &api.GetInstanceInfoResp{
		Info: &api.InstanceInfo{
			Name:        r.Host,
			TotalMem:    r.TotalMem,
			CurMemUsage: r.CurMemUsage,
			CpuNum:      r.CpuNum,
			CurCpuUsage: r.CurCpuUsage,
		},
	}, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
