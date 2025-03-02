package initinternal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"

	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/secure"
	// "github.com/chenjie199234/Corelib/trace"
	"github.com/chenjie199234/Corelib/util/common"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func InitDatabase(secret string, db *mongo.Client) (e error) {
	var ac []byte
	var sc []byte
	if ac, e = os.ReadFile("./AppConfig.json"); e != nil {
		slog.ErrorContext(nil, "[InitDatabase] read ./AppConfig.json failed", slog.String("error", e.Error()))
		return
	}
	if sc, e = os.ReadFile("./SourceConfig.json"); e != nil {
		slog.ErrorContext(nil, "[InitDatabase] read ./SourceConfig.json failed", slog.String("error", e.Error()))
		return
	}
	bufapp := bytes.NewBuffer(nil)
	if e = json.Compact(bufapp, ac); e != nil {
		slog.ErrorContext(nil, "[InitDatabase] ./AppConfig.json format wrong", slog.String("error", e.Error()))
		return
	}
	bufsource := bytes.NewBuffer(nil)
	if e = json.Compact(bufsource, sc); e != nil {
		slog.ErrorContext(nil, "[InitDatabase] ./SourceConfig.json format wrong", slog.String("error", e.Error()))
		return e
	}
	appconfig := ""
	sourceconfig := ""
	if secret != "" {
		appconfig, _ = secure.AesEncrypt(secret, bufapp.Bytes())
		sourceconfig, _ = secure.AesEncrypt(secret, bufsource.Bytes())
	} else {
		appconfig = common.BTS(bufapp.Bytes())
		sourceconfig = common.BTS(bufsource.Bytes())
	}

	var needcommit bool
	var s *mongo.Session
	if s, e = db.StartSession(); e != nil {
		slog.ErrorContext(nil, "[InitDatabase] start mongo session failed", slog.String("mongo", "admin_mongo"), slog.String("error", e.Error()))
		return
	}
	defer s.EndSession(context.Background())
	sctx := mongo.NewSessionContext(context.Background(), s)
	if e = s.StartTransaction(options.Transaction().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())); e != nil {
		slog.ErrorContext(nil, "[InitDatabase] start mongo transaction failed", slog.String("mongo", "admin_mongo"), slog.String("error", e.Error()))
		return
	}
	defer func() {
		if !needcommit || e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
			slog.ErrorContext(nil, "[InitDatabase] commit mongo failed", slog.String("mongo", "admin_mongo"), slog.String("error", e.Error()))
		}
	}()

	//check init status

	//check project index
	existProjectIndex := &model.ProjectIndex{}
	e = db.Database("permission").Collection("projectindex").FindOne(sctx, bson.M{"project_id": model.AdminProjectID}).Decode(existProjectIndex)
	if e != nil && e != mongo.ErrNoDocuments {
		slog.ErrorContext(nil, "[InitDatabase] get project index failed", slog.String("project_id", model.AdminProjectID), slog.String("error", e.Error()))
		return
	}
	if e == nil && existProjectIndex.ProjectName != model.Project {
		slog.ErrorContext(nil, "[InitDatabase] already inited with other project name",
			slog.String("project_id", model.AdminProjectID),
			slog.String("exist_project_name", existProjectIndex.ProjectName),
			slog.String("current_project_name", model.Project))
		e = errors.New("conflict")
		return
	}
	//check node
	nodefilter := bson.M{
		"node_id": bson.M{
			"$in": bson.A{
				"0",
				model.AdminProjectID,
				model.AdminProjectID + model.UserAndRoleControl,
				model.AdminProjectID + model.AppControl,
				model.AdminProjectID + model.AppControl + ",1",
			},
		},
	}
	if existProjectIndex.ProjectName != "" {
		var c *mongo.Cursor
		c, e = db.Database("permission").Collection("node").Find(sctx, nodefilter)
		if e != nil {
			slog.ErrorContext(nil, "[InitDatabase] get nodes failed", slog.String("error", e.Error()))
			return
		}
		nodes := make([]*model.Node, 0, c.RemainingBatchLength())
		if e = c.All(sctx, &nodes); e != nil {
			slog.ErrorContext(nil, "[InitDatabase] get nodes failed", slog.String("error", e.Error()))
			return
		}
		if len(nodes) != 5 {
			slog.ErrorContext(nil, "[InitDatabase] basic nodes missing")
			e = errors.New("dirty")
			return
		}
		for _, node := range nodes {
			dirty := false
			switch node.NodeId {
			case "0":
				if node.NodeName != "root" {
					dirty = true
				}
			case model.AdminProjectID:
				if node.NodeName != model.Project {
					dirty = true
				}
			case model.AdminProjectID + model.UserAndRoleControl:
				if node.NodeName != "UserAndRoleControl" {
					dirty = true
				}
			case model.AdminProjectID + model.AppControl:
				if node.NodeName != "AppControl" {
					dirty = true
				}
			case model.AdminProjectID + model.AppControl + ",1":
				if node.NodeName != model.Group+"."+model.Name {
					dirty = true
				}
			}
			if dirty {
				slog.ErrorContext(nil, "[InitDatabase] basic node data dirty")
				e = errors.New("dirty")
				return
			}
		}
	}
	//check app
	existAppSummary := &model.AppSummary{}
	appSummaryFilter := bson.M{
		"project_id": model.AdminProjectID,
		"group":      model.Group,
		"app":        model.Name,
		"key":        "",
		"index":      0,
	}
	e = db.Database("app").Collection("config").FindOne(sctx, appSummaryFilter).Decode(existAppSummary)
	if existProjectIndex.ProjectName == "" && e != mongo.ErrNoDocuments {
		//project not exist,the app should not exist too
		if e == nil {
			slog.ErrorContext(nil, "[InitDatabase] project not exist but app already exist",
				slog.String("project_id", model.AdminProjectID),
				slog.String("group", model.Group),
				slog.String("app", model.Name))
			e = errors.New("dirty")
		} else {
			slog.ErrorContext(nil, "[InitDatabase] get app failed",
				slog.String("project_id", model.AdminProjectID),
				slog.String("group", model.Group),
				slog.String("app", model.Name),
				slog.String("error", e.Error()))
		}
		return
	}
	if existProjectIndex.ProjectName != "" && e != nil {
		//project exist,the app should exist too
		if e == mongo.ErrNoDocuments {
			slog.ErrorContext(nil, "[InitDatabase] project exist but app not exist",
				slog.String("project_id", model.AdminProjectID),
				slog.String("group", model.Group),
				slog.String("app", model.Name))
			e = errors.New("dirty")
		} else {
			slog.ErrorContext(nil, "[InitDatabase] get app failed",
				slog.String("project_id", model.AdminProjectID),
				slog.String("group", model.Group),
				slog.String("app", model.Name),
				slog.String("error", e.Error()))
		}
		return
	}
	if existProjectIndex.ProjectName != "" {
		//project exist,the app should exist too
		//check secret
		if e = secure.SignCheck(secret, existAppSummary.Value); e != nil {
			slog.ErrorContext(nil, "[InitDatabase] secret check failed",
				slog.String("project_id", model.AdminProjectID),
				slog.String("group", model.Group),
				slog.String("app", model.Name),
				slog.String("error", e.Error()))
			return
		}
	}
	if existProjectIndex.ProjectName != "" {
		//already inited and there is no conflict and the secret is right
		return
	}

	//init now
	needcommit = true

	//init project index
	if _, e = db.Database("permission").Collection("projectindex").InsertOne(sctx, bson.M{"project_name": model.Project, "project_id": model.AdminProjectID}); e != nil {
		slog.ErrorContext(nil, "[InitDatabase] init project index failed",
			slog.String("project_id", model.AdminProjectID),
			slog.String("project_name", model.Project),
			slog.String("error", e.Error()))
		return
	}
	//init node
	nodeids := make([]string, 0, 10)
	docs := bson.A{}
	docs = append(docs, model.Node{
		NodeId:       "0",
		NodeName:     "root",
		NodeData:     "",
		CurNodeIndex: 1,
	})
	nodeids = append(nodeids, "0")
	//first project's node
	docs = append(docs, &model.Node{
		NodeId:       model.AdminProjectID,
		NodeName:     model.Project,
		NodeData:     "",
		CurNodeIndex: 100,
	})
	nodeids = append(nodeids, model.AdminProjectID)
	//first project's user and role control node
	docs = append(docs, &model.Node{
		NodeId:       model.AdminProjectID + model.UserAndRoleControl,
		NodeName:     "UserAndRoleControl",
		NodeData:     "",
		CurNodeIndex: 0,
	})
	nodeids = append(nodeids, model.AdminProjectID+model.UserAndRoleControl)
	//first project's config control node
	docs = append(docs, &model.Node{
		NodeId:       model.AdminProjectID + model.AppControl,
		NodeName:     "AppControl",
		NodeData:     "",
		CurNodeIndex: 1,
	})
	nodeids = append(nodeids, model.AdminProjectID+model.AppControl)
	//first project's first app(this app: admin)'s config node
	docs = append(docs, &model.Node{
		NodeId:       model.AdminProjectID + model.AppControl + ",1",
		NodeName:     model.Group + "." + model.Name,
		NodeData:     "",
		CurNodeIndex: 0,
	})
	nodeids = append(nodeids, model.AdminProjectID+model.AppControl+",1")
	if _, e = db.Database("permission").Collection("node").InsertMany(sctx, docs); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			slog.ErrorContext(nil, "[InitDatabase] project and app not exist but some permission nodes already exist", slog.Any("node_ids", nodeids), slog.String("error", e.Error()))
			e = errors.New("dirty")
		} else {
			slog.ErrorContext(nil, "[InitDatabase] init permission nodes failed", slog.Any("node_ids", nodeids), slog.String("error", e.Error()))
		}
		return
	}
	//init app
	docs = docs[0:0]
	sign, _ := secure.SignMake(secret)
	//summary
	docs = append(docs, &model.AppSummary{
		ProjectID:    model.AdminProjectID,
		ProjectName:  model.Project,
		Group:        model.Group,
		App:          model.Name,
		Key:          "",
		Index:        0,
		DiscoverMode: "",
		KubernetesNs: "",
		KubernetesLS: "",
		KubernetesFS: "",
		DnsHost:      "",
		DnsInterval:  0,
		StaticAddrs:  nil,
		Keys: map[string]*model.KeySummary{
			"AppConfig": {
				CurIndex:     1,
				CurVersion:   1,
				MaxIndex:     1,
				CurValue:     appconfig,
				CurValueType: "json",
			},
			"SourceConfig": {
				CurIndex:     1,
				CurVersion:   1,
				MaxIndex:     1,
				CurValue:     sourceconfig,
				CurValueType: "json",
			},
		},
		Value:            sign,
		PermissionNodeID: model.AdminProjectID + model.AppControl + ",1",
	})
	//AppConfig
	docs = append(docs, &model.Log{
		ProjectID: model.AdminProjectID,
		Group:     model.Group,
		App:       model.Name,
		Key:       "AppConfig",
		Index:     1,
		Value:     appconfig,
		ValueType: "json",
	})
	//SourceConfig
	docs = append(docs, &model.Log{
		ProjectID: model.AdminProjectID,
		Group:     model.Group,
		App:       model.Name,
		Key:       "SourceConfig",
		Index:     1,
		Value:     sourceconfig,
		ValueType: "json",
	})
	if _, e = db.Database("app").Collection("config").InsertMany(sctx, docs); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			slog.ErrorContext(nil, "[InitDatabase] project not exist but app already exist",
				slog.String("project_id", model.AdminProjectID),
				slog.String("group", model.Group),
				slog.String("app", model.Name))
			e = errors.New("dirty")
		} else {
			slog.ErrorContext(nil, "[InitDatabase] init app failed",
				slog.String("project_id", model.AdminProjectID),
				slog.String("group", model.Group),
				slog.String("app", model.Name),
				slog.String("error", e.Error()))
		}
	}
	return
}
