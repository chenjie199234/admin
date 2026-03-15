package app

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"slices"
	"strings"

	"github.com/chenjie199234/admin/api"
	"github.com/chenjie199234/admin/config"
	appdao "github.com/chenjie199234/admin/dao/app"
	initializedao "github.com/chenjie199234/admin/dao/initialize"
	permissiondao "github.com/chenjie199234/admin/dao/permission"
	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	// "github.com/chenjie199234/Corelib/web"
	// "github.com/chenjie199234/Corelib/crpc"
	// "github.com/chenjie199234/Corelib/cgrpc"
	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/metadata"
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/egroup"
	"github.com/chenjie199234/Corelib/util/graceful"
	"github.com/chenjie199234/Corelib/util/name"
	"go.mongodb.org/mongo-driver/v2/bson"
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
		slog.ErrorContext(ctx, "[GetApp] operator's token format wrong",
			slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GetGName(), req.GetAName())
		if e != nil {
			slog.ErrorContext(ctx, "[GetApp] get app's permission nodeid failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			slog.ErrorContext(ctx, "[GetApp] app's permission nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
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
	app, e := s.appDao.MongoGetApp(ctx, projectid, req.GetGName(), req.GetAName(), req.GetSecret())
	if e != nil {
		slog.ErrorContext(ctx, "[GetApp] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	keys := make(map[string]*api.KeyConfigInfo)
	for k, v := range app.Keys {
		kinfo := &api.KeyConfigInfo{}
		kinfo.SetCurIndex(v.CurIndex)
		kinfo.SetMaxIndex(v.MaxIndex)
		kinfo.SetCurVersion(v.CurVersion)
		kinfo.SetCurValue(v.CurValue)
		kinfo.SetCurValueType(v.CurValueType)
		keys[k] = kinfo
	}
	resp := &api.GetAppResp{}
	resp.SetDiscoverMode(app.DiscoverMode)
	resp.SetKubernetesNamespace(app.KubernetesNs)
	resp.SetKubernetesLabelselector(app.KubernetesLS)
	resp.SetKubernetesFieldselector(app.KubernetesFS)
	resp.SetDnsHost(app.DnsHost)
	resp.SetDnsInterval(app.DnsInterval)
	resp.SetStaticAddrs(app.StaticAddrs)
	resp.SetCrpcPort(app.CrpcPort)
	resp.SetCgrpcPort(app.CGrpcPort)
	resp.SetWebPort(app.WebPort)
	resp.SetKeys(keys)
	return resp, nil
}

func (s *Service) SetApp(ctx context.Context, req *api.SetAppReq) (*api.SetAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[SetApp] operator's token format wrong",
			slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}
	if e := name.SingleCheck(req.GetGName(), false); e != nil {
		slog.ErrorContext(ctx, "[SetApp] group name format wrong",
			slog.String("operator", md["Token-User"]), slog.String("group", req.GetGName()))
		return nil, ecode.ErrReq
	}
	if e := name.SingleCheck(req.GetAName(), false); e != nil {
		slog.ErrorContext(ctx, "[SetApp] app name format wrong",
			slog.String("operator", md["Token-User"]), slog.String("app", req.GetAName()))
		return nil, ecode.ErrReq
	}
	switch req.GetDiscoverMode() {
	case "kubernetes":
		if req.GetKubernetesNamespace() == "" {
			slog.ErrorContext(ctx, "[SetApp] kubernetes namesapce empty", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
		if req.GetKubernetesLabelselector() == "" && req.GetKubernetesFieldselector() == "" {
			slog.ErrorContext(ctx, "[SetApp] kubernetes labelselector and fieldselector empty", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
	case "dns":
		if req.GetDnsHost() == "" {
			slog.ErrorContext(ctx, "[SetApp] dns host empty", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
		if req.GetDnsInterval() == 0 {
			slog.ErrorContext(ctx, "[SetApp] dns interval must be set", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
	case "static":
		if len(req.GetStaticAddrs()) == 0 {
			slog.ErrorContext(ctx, "[SetApp] static addrs empty", slog.String("operator", md["Token-User"]))
			return nil, ecode.ErrReq
		}
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//config control permission check
		if req.GetNewApp() {
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
			nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GetGName(), req.GetAName())
			if e != nil {
				slog.ErrorContext(ctx, "[SetApp] get app's permission nodeid failed",
					slog.String("operator", md["Token-User"]),
					slog.String("project_id", projectid),
					slog.String("group", req.GetGName()),
					slog.String("app", req.GetAName()),
					slog.String("error", e.Error()))
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			nodeids := strings.Split(nodeid, ",")
			if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
				slog.ErrorContext(ctx, "[SetApp] app's permission nodeid format wrong",
					slog.String("operator", md["Token-User"]),
					slog.String("project_id", projectid),
					slog.String("group", req.GetGName()),
					slog.String("app", req.GetAName()),
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
	if req.GetNewApp() {
		nodeidstr, e = s.appDao.MongoCreateApp(
			ctx,
			projectid,
			req.GetGName(),
			req.GetAName(),
			req.GetSecret(),
			req.GetDiscoverMode(),
			req.GetKubernetesNamespace(),
			req.GetKubernetesLabelselector(),
			req.GetKubernetesFieldselector(),
			req.GetDnsHost(),
			req.GetDnsInterval(),
			req.GetStaticAddrs(),
			req.GetCrpcPort(),
			req.GetCgrpcPort(),
			req.GetWebPort())
	} else {
		nodeidstr, e = s.appDao.MongoUpdateApp(
			ctx,
			projectid,
			req.GetGName(),
			req.GetAName(),
			req.GetSecret(),
			req.GetDiscoverMode(),
			req.GetKubernetesNamespace(),
			req.GetKubernetesLabelselector(),
			req.GetKubernetesFieldselector(),
			req.GetDnsHost(),
			req.GetDnsInterval(),
			req.GetStaticAddrs(),
			req.GetCrpcPort(),
			req.GetCgrpcPort(),
			req.GetWebPort())
	}
	if e != nil {
		slog.ErrorContext(ctx, "[SetApp] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeid, e := util.ParseID(nodeidstr)
	if e != nil {
		slog.ErrorContext(ctx, "[SetApp] nodeid format wrong",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", ecode.ErrDBDataBroken.Error()))
	}
	slog.InfoContext(ctx, "[SetApp] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GetGName()),
		slog.String("app", req.GetAName()))
	resp := &api.SetAppResp{}
	resp.SetNodeId(nodeid)
	return resp, nil
}

func (s *Service) DelApp(ctx context.Context, req *api.DelAppReq) (*api.DelAppResp, error) {
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[DelApp] operator's token format wrong",
			slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GetGName(), req.GetAName())
	if e != nil {
		slog.ErrorContext(ctx, "[DelApp] get app's permission nodeid failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		slog.ErrorContext(ctx, "[DelApp] app's permission nodeid format wrong",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("nodeid", nodeid))
		return nil, ecode.ErrDBDataBroken
	}
	//self can't be deleted
	if nodeids[1] == "1" && nodeids[3] == "1" {
		slog.ErrorContext(ctx, "[DelApp] can't delete self",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()))
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
	if e := s.appDao.MongoDelApp(ctx, projectid, req.GetGName(), req.GetAName(), req.GetSecret()); e != nil {
		slog.ErrorContext(ctx, "[DelApp] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[DelApp] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GetGName()),
		slog.String("app", req.GetAName()))
	return &api.DelAppResp{}, nil
}

func (s *Service) UpdateAppSecret(ctx context.Context, req *api.UpdateAppSecretReq) (*api.UpdateAppSecretResp, error) {
	if req.GetOldSecret() == req.GetNewSecret() {
		return &api.UpdateAppSecretResp{}, nil
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[UpdateAppSecret] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GetGName(), req.GetAName())
		if e != nil {
			slog.ErrorContext(ctx, "[UpdateAppSecret] get app's permission nodeid failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			slog.ErrorContext(ctx, "[UpdateAppSecret] app's permission nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
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
	if e := s.appDao.MongoUpdateAppSecret(ctx, projectid, req.GetGName(), req.GetAName(), req.GetOldSecret(), req.GetNewSecret()); e != nil {
		slog.ErrorContext(ctx, "[UpdateAppSecret] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[UpdateAppSecret] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GetGName()),
		slog.String("app", req.GetAName()))
	return &api.UpdateAppSecretResp{}, nil
}

func (s *Service) DelKey(ctx context.Context, req *api.DelKeyReq) (*api.DelKeyResp, error) {
	if strings.Contains(req.GetKey(), ".") || strings.Contains(req.GetKey(), "$") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[DelKey] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GetGName(), req.GetAName())
	if e != nil {
		slog.ErrorContext(ctx, "[DelKey] get app's permission nodeid failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	nodeids := strings.Split(nodeid, ",")
	if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
		slog.ErrorContext(ctx, "[DelKey] app's permission nodeid format wrong",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("nodeid", nodeid))
		return nil, ecode.ErrDBDataBroken
	}
	if nodeids[1] == "1" && nodeids[3] == "1" && (req.GetKey() == "AppConfig" || req.GetKey() == "SourceConfig") {
		slog.ErrorContext(ctx, "[DelKey] can't delete self's 'AppConfig' or 'SourceConfig' key",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()))
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

	if e := s.appDao.MongoDelKey(ctx, projectid, req.GetGName(), req.GetAName(), req.GetKey(), req.GetSecret()); e != nil {
		slog.ErrorContext(ctx, "[DelKey] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("key", req.GetKey()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[DelKey] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GetGName()),
		slog.String("app", req.GetAName()),
		slog.String("key", req.GetKey()))
	return &api.DelKeyResp{}, nil
}

func (s *Service) GetKeyConfig(ctx context.Context, req *api.GetKeyConfigReq) (*api.GetKeyConfigResp, error) {
	if strings.Contains(req.GetKey(), ".") || strings.Contains(req.GetKey(), "$") {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[GetKeyConfig] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GetGName(), req.GetAName())
		if e != nil {
			slog.ErrorContext(ctx, "[GetKeyConfig] get app's permission nodeid failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			slog.ErrorContext(ctx, "[GetKeyConfig] app's permission nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
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
	keysummary, configlog, e := s.appDao.MongoGetKeyConfig(ctx, projectid, req.GetGName(), req.GetAName(), req.GetKey(), req.GetIndex(), req.GetSecret())
	if e != nil {
		slog.ErrorContext(ctx, "[GetKeyConfig] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("key", req.GetKey()),
			slog.Uint64("index", uint64(req.GetIndex())),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.GetKeyConfigResp{}
	resp.SetCurIndex(keysummary.CurIndex)
	resp.SetMaxIndex(keysummary.MaxIndex)
	resp.SetCurVersion(keysummary.CurVersion)
	resp.SetThisIndex(configlog.Index)
	resp.SetValue(configlog.Value)
	resp.SetValueType(configlog.ValueType)
	return resp, nil
}

func (s *Service) SetKeyConfig(ctx context.Context, req *api.SetKeyConfigReq) (*api.SetKeyConfigResp, error) {
	if strings.Contains(req.GetKey(), ".") || strings.Contains(req.GetKey(), "$") {
		return nil, ecode.ErrReq
	}

	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[SetKeyConfig] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	req.SetKey(strings.TrimSpace(req.GetKey()))
	if req.GetKey() == "" {
		slog.ErrorContext(ctx, "[SetKeyConfig] key empty",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()))
		return nil, ecode.ErrReq
	}
	req.SetValue(strings.TrimSpace(req.GetValue()))
	if req.GetValue() == "" {
		slog.ErrorContext(ctx, "[SetKeyConfig] value empty",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("key", req.GetKey()))
		return nil, ecode.ErrReq
	}
	switch req.GetValueType() {
	case "json":
		buf := bytes.NewBuffer(nil)
		if e := json.Compact(buf, common.STB(req.GetValue())); e != nil {
			slog.ErrorContext(ctx, "[SetKeyConfig] json value format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
				slog.String("key", req.GetKey()),
				slog.String("error", e.Error()))
			return nil, ecode.ErrReq
		}
		req.SetValue(common.BTS(buf.Bytes()))
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
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("key", req.GetKey()),
			slog.String("valuetype", req.GetValueType()))
		return nil, ecode.ErrReq
	}

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GetGName(), req.GetAName())
		if e != nil {
			slog.ErrorContext(ctx, "[SetKeyConfig] get app's permission nodeid failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			slog.ErrorContext(ctx, "[SetKeyConfig] app's permission nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
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
	index, version, e := s.appDao.MongoSetKeyConfig(ctx, projectid, req.GetGName(), req.GetAName(), req.GetKey(), req.GetSecret(), req.GetValue(), req.GetValueType(), req.GetNewKey())
	if e != nil {
		slog.ErrorContext(ctx, "[SetKeyConfig] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("key", req.GetKey()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[SetKeyConfig] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GetGName()),
		slog.String("app", req.GetAName()),
		slog.String("key", req.GetKey()),
		slog.Uint64("new_index", uint64(index)),
		slog.Uint64("new_version", uint64(version)))
	return &api.SetKeyConfigResp{}, nil
}

func (s *Service) Rollback(ctx context.Context, req *api.RollbackReq) (*api.RollbackResp, error) {
	if strings.Contains(req.GetKey(), ".") || strings.Contains(req.GetKey(), "$") {
		return nil, ecode.ErrReq
	}
	md := metadata.GetMetadata(ctx)
	operator, e := bson.ObjectIDFromHex(md["Token-User"])
	if e != nil {
		slog.ErrorContext(ctx, "[Rollback] operator's token format wrong", slog.String("operator", md["Token-User"]), slog.String("error", e.Error()))
		return nil, ecode.ErrToken
	}

	projectid := util.FormID(req.GetProjectId())

	if !operator.IsZero() {
		//config control permission check
		nodeid, e := s.appDao.MongoGetPermissionNodeID(ctx, projectid, req.GetGName(), req.GetAName())
		if e != nil {
			slog.ErrorContext(ctx, "[Rollback] get app's permission nodeid failed",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
				slog.String("error", e.Error()))
			return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
		}
		nodeids := strings.Split(nodeid, ",")
		if len(nodeids) != 4 || nodeids[0] != "0" || nodeids[2] != "2" {
			slog.ErrorContext(ctx, "[Rollback] app's permission nodeid format wrong",
				slog.String("operator", md["Token-User"]),
				slog.String("project_id", projectid),
				slog.String("group", req.GetGName()),
				slog.String("app", req.GetAName()),
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
	if e := s.appDao.MongoRollbackKeyConfig(ctx, projectid, req.GetGName(), req.GetAName(), req.GetKey(), req.GetSecret(), req.GetIndex()); e != nil {
		slog.ErrorContext(ctx, "[Rollback] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("key", req.GetKey()),
			slog.Uint64("index", uint64(req.GetIndex())),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	slog.InfoContext(ctx, "[Rollback] success",
		slog.String("operator", md["Token-User"]),
		slog.String("project_id", projectid),
		slog.String("group", req.GetGName()),
		slog.String("app", req.GetAName()),
		slog.String("key", req.GetKey()),
		slog.Uint64("index", uint64(req.GetIndex())))
	return &api.RollbackResp{}, nil
}

func (s *Service) WatchConfig(ctx context.Context, req *api.WatchConfigReq) (*api.WatchConfigResp, error) {
	for k := range req.GetKeys() {
		if strings.Contains(k, ".") || strings.Contains(k, "$") {
			return nil, ecode.ErrReq
		}
	}
	ch, cancel, e := config.Sdk.GetNoticeByProjectName(req.GetProjectName(), req.GetGName(), req.GetAName())
	if e != nil {
		slog.ErrorContext(ctx, "[WatchConfig] get notice failed",
			slog.String("project_name", req.GetProjectName()),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil, cerror.Convert(ctx.Err())
		case <-ch:
			app, e := config.Sdk.GetAppConfigByProjectName(req.GetProjectName(), req.GetGName(), req.GetAName())
			if e != nil {
				if e != ecode.ErrServerClosing {
					slog.ErrorContext(ctx, "[WatchConfig] get config failed",
						slog.String("project_name", req.GetProjectName()),
						slog.String("group", req.GetGName()),
						slog.String("app", req.GetAName()),
						slog.String("error", e.Error()))
				}
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			datas := make(map[string]*api.WatchData, len(req.GetKeys())+3)
			needreturn := false
			for key, clientversion := range req.GetKeys() {
				k, ok := app.Keys[key]
				if !ok || k == nil || k.CurVersion == 0 {
					return nil, ecode.ErrKeyNotExist
				}
				if clientversion != k.CurVersion {
					needreturn = true
					d := &api.WatchData{}
					d.SetKey(key)
					d.SetValue(k.CurValue)
					d.SetValueType(k.CurValueType)
					d.SetVersion(k.CurVersion)
					datas[key] = d
				} else {
					d := &api.WatchData{}
					d.SetKey(key)
					d.SetValue("")
					d.SetValueType("")
					d.SetVersion(k.CurVersion)
					datas[key] = d
				}
			}
			if needreturn {
				resp := &api.WatchConfigResp{}
				resp.SetDatas(datas)
				return resp, nil
			}
		}
	}
}

func (s *Service) WatchDiscover(ctx context.Context, req *api.WatchDiscoverReq) (*api.WatchDiscoverResp, error) {
	if req.GetCurDiscoverMode() == "dns" && (req.GetCurDnsHost() == "" || req.GetCurDnsInterval() == 0) {
		return nil, ecode.ErrReq
	}
	if req.GetCurDiscoverMode() == "static" && len(req.GetCurStaticAddrs()) == 0 {
		return nil, ecode.ErrReq
	}
	if req.GetCurDiscoverMode() == "kubernetes" && (req.GetCurKubernetesNamespace() == "" || (req.GetCurKubernetesFieldselector() == "" && req.GetCurKubernetesLabelselector() == "")) {
		return nil, ecode.ErrReq
	}
	ch, cancel, e := config.Sdk.GetNoticeByProjectName(req.GetProjectName(), req.GetGName(), req.GetAName())
	if e != nil {
		slog.ErrorContext(ctx, "[WatchConfig] get notice failed",
			slog.String("project_name", req.GetProjectName()),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return nil, cerror.Convert(ctx.Err())
		case <-ch:
			app, e := config.Sdk.GetAppConfigByProjectName(req.GetProjectName(), req.GetGName(), req.GetAName())
			if e != nil {
				if e != ecode.ErrServerClosing {
					slog.ErrorContext(ctx, "[WatchDiscover] get config failed",
						slog.String("project_name", req.GetProjectName()),
						slog.String("group", req.GetGName()),
						slog.String("app", req.GetAName()),
						slog.String("error", e.Error()))
				}
				return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
			}
			needreturn := app.DiscoverMode != req.GetCurDiscoverMode() ||
				app.CrpcPort != req.GetCurCrpcPort() ||
				app.WebPort != req.GetCurWebPort() ||
				app.CGrpcPort != req.GetCurCgrpcPort()
			if !needreturn {
				switch app.DiscoverMode {
				case "dns":
					needreturn = app.DnsHost != req.GetCurDnsHost() ||
						app.DnsInterval != req.GetCurDnsInterval()
				case "static":
					for _, addr := range app.StaticAddrs {
						if !slices.Contains(req.GetCurStaticAddrs(), addr) {
							needreturn = true
						}
					}
					for _, addr := range req.GetCurStaticAddrs() {
						if !slices.Contains(app.StaticAddrs, addr) {
							needreturn = true
						}
					}
				case "kubernetes":
					needreturn = app.KubernetesNs != req.GetCurKubernetesNamespace() ||
						app.KubernetesFS != req.GetCurKubernetesFieldselector() ||
						app.KubernetesLS != req.GetCurKubernetesLabelselector()
				}
			}
			if needreturn {
				resp := &api.WatchDiscoverResp{}
				resp.SetDiscoverMode(app.DiscoverMode)
				resp.SetDnsHost(app.DnsHost)
				resp.SetDnsInterval(app.DnsInterval)
				resp.SetStaticAddrs(app.StaticAddrs)
				resp.SetKubernetesNamespace(app.KubernetesNs)
				resp.SetKubernetesLabelselector(app.KubernetesLS)
				resp.SetKubernetesFieldselector(app.KubernetesFS)
				resp.SetCrpcPort(app.CrpcPort)
				resp.SetCgrpcPort(app.CGrpcPort)
				resp.SetWebPort(app.WebPort)
				return resp, nil
			}
		}
	}
}

func (s *Service) GetInstances(ctx context.Context, req *api.GetInstancesReq) (*api.GetInstancesResp, error) {
	md := metadata.GetMetadata(ctx)

	projectid := util.FormID(req.GetProjectId())

	if e := s.appDao.MongoCheckSecret(ctx, projectid, req.GetGName(), req.GetAName(), req.GetSecret()); e != nil {
		slog.ErrorContext(ctx, "[GetInstances] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}

	addrs, e := config.Sdk.GetAppAddrsByProjectID(ctx, projectid, req.GetGName(), req.GetAName())
	if e != nil {
		slog.ErrorContext(ctx, "[GetInstances] get addrs failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.GetInstancesResp{}
	instances := make(map[string]*api.InstanceInfo, len(addrs))
	for _, addr := range addrs {
		instances[addr] = nil
	}
	if !req.GetWithInfo() {
		resp.SetInstances(instances)
		return resp, nil
	}
	eg := egroup.GetGroup(ctx)
	for _, v := range addrs {
		addr := v
		eg.Go(func(gctx context.Context) error {
			r, e := config.Sdk.PingByPrjoectID(gctx, projectid, req.GetGName(), req.GetAName(), addr)
			if e != nil {
				slog.ErrorContext(ctx, "[GetInstances] get info failed",
					slog.String("operator", md["Token-User"]),
					slog.String("project_id", projectid),
					slog.String("group", req.GetGName()),
					slog.String("app", req.GetAName()),
					slog.String("addr", addr),
					slog.String("error", e.Error()))
				return nil
			}
			info := &api.InstanceInfo{}
			info.SetName(r.GetHost())
			info.SetCpuNum(r.GetCpuNum())
			info.SetCpuUsage(r.GetCpuUsage())
			info.SetCpuType(r.GetCpuType())
			info.SetMemTotal(r.GetMemTotal())
			info.SetMemUsage(r.GetMemUsage())
			info.SetMemType(r.GetMemType())
			instances[addr] = info
			return nil
		})
	}
	egroup.PutGroup(eg)
	resp.SetInstances(instances)
	return resp, nil
}
func (s *Service) GetInstanceInfo(ctx context.Context, req *api.GetInstanceInfoReq) (*api.GetInstanceInfoResp, error) {
	md := metadata.GetMetadata(ctx)

	projectid := util.FormID(req.GetProjectId())

	if e := s.appDao.MongoCheckSecret(ctx, projectid, req.GetGName(), req.GetAName(), req.GetSecret()); e != nil {
		slog.ErrorContext(ctx, "[GetInstanceInfo] db op failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	r, e := config.Sdk.PingByPrjoectID(ctx, projectid, req.GetGName(), req.GetAName(), req.GetAddr())
	if e != nil {
		slog.ErrorContext(ctx, "[GetInstanceInfo] get info failed",
			slog.String("operator", md["Token-User"]),
			slog.String("project_id", projectid),
			slog.String("group", req.GetGName()),
			slog.String("app", req.GetAName()),
			slog.String("addr", req.GetAddr()),
			slog.String("error", e.Error()))
		return nil, ecode.ReturnEcode(e, ecode.ErrSystem)
	}
	resp := &api.GetInstanceInfoResp{}
	info := &api.InstanceInfo{}
	info.SetName(r.GetHost())
	info.SetCpuNum(r.GetCpuNum())
	info.SetCpuUsage(r.GetCpuUsage())
	info.SetCpuType(r.GetCpuType())
	info.SetMemTotal(r.GetMemTotal())
	info.SetMemUsage(r.GetMemUsage())
	info.SetMemType(r.GetMemType())
	resp.SetInfo(info)
	return resp, nil
}

// Stop -
func (s *Service) Stop() {
	s.stop.Close(nil, nil)
}
