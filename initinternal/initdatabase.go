package initinternal

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/secure"
	"github.com/chenjie199234/Corelib/util/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func InitDatabase(secret string, db *mongo.Client) (e error) {
	var ac []byte
	var sc []byte
	if ac, e = os.ReadFile("./AppConfig.json"); e != nil {
		log.Error(nil, "[InitSelf] read ./AppConfig.json failed", map[string]interface{}{"error": e})
		return
	}
	if sc, e = os.ReadFile("./SourceConfig.json"); e != nil {
		log.Error(nil, "[InitSelf] read ./SourceConfig.json failed", map[string]interface{}{"error": e})
		return
	}
	bufapp := bytes.NewBuffer(nil)
	if e = json.Compact(bufapp, ac); e != nil {
		log.Error(nil, "[InitSelf] ./AppConfig.json format wrong", map[string]interface{}{"error": e})
		return
	}
	bufsource := bytes.NewBuffer(nil)
	if e = json.Compact(bufsource, sc); e != nil {
		log.Error(nil, "[InitSelf] ./SourceConfig.json format wrong", map[string]interface{}{"error": e})
		return e
	}
	appconfig := ""
	sourceconfig := ""
	if secret != "" {
		appconfig, _ = secure.AesEncrypt(secret, bufapp.Bytes())
		sourceconfig, _ = secure.AesEncrypt(secret, bufsource.Bytes())
	} else {
		appconfig = common.Byte2str(bufapp.Bytes())
		sourceconfig = common.Byte2str(bufsource.Bytes())
	}

	var needcommit bool
	var s mongo.Session
	if s, e = db.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local())); e != nil {
		log.Error(nil, "[InitSelf] start mongo session failed", map[string]interface{}{"mongo": "admin_mongo", "error": e})
		return
	}
	defer s.EndSession(context.Background())
	sctx := mongo.NewSessionContext(context.Background(), s)
	if e = s.StartTransaction(); e != nil {
		log.Error(nil, "[InitSelf] start mongo transaction failed", map[string]interface{}{"mongo": "admin_mongo", "error": e})
		return
	}
	defer func() {
		if !needcommit || e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
			log.Error(nil, "[InitSelf] commit mongo failed", map[string]interface{}{"mongo": "admin_mongo", "error": e})
		}
	}()

	//check init status

	//check project index
	existProjectIndex := &model.ProjectIndex{}
	e = db.Database("permission").Collection("projectindex").FindOne(sctx, bson.M{"project_id": model.AdminProjectID}).Decode(existProjectIndex)
	if e != nil && e != mongo.ErrNoDocuments {
		log.Error(nil, "[InitSelf] get project index failed", map[string]interface{}{"project_id": model.AdminProjectID, "error": e})
		return
	}
	if e == nil && existProjectIndex.ProjectName != model.Project {
		log.Error(nil, "[InitSelf] already inited with other project name", map[string]interface{}{
			"project_id":           model.AdminProjectID,
			"exist_project_name":   existProjectIndex.ProjectName,
			"current_project_name": model.Project})
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
			log.Error(nil, "[InitSelf] get nodes failed", map[string]interface{}{"error": e})
			return
		}
		nodes := make([]*model.Node, 0, c.RemainingBatchLength())
		if e = c.All(context.Background(), &nodes); e != nil {
			log.Error(nil, "[InitSelf] get nodes failed", map[string]interface{}{"error": e})
			return
		}
		if len(nodes) != 5 {
			log.Error(nil, "[InitSelf] basic nodes missing", nil)
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
				log.Error(nil, "[InitSelf] basic node data dirty", nil)
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
			log.Error(nil, "[InitSelf] project not exist but app already exist", map[string]interface{}{
				"project_id": model.AdminProjectID,
				"group":      model.Group,
				"app":        model.Name})
			e = errors.New("dirty")
		} else {
			log.Error(nil, "[InitSelf] get app failed", map[string]interface{}{
				"project_id": model.AdminProjectID,
				"group":      model.Group,
				"app":        model.Name,
				"error":      e})
		}
		return
	}
	if existProjectIndex.ProjectName != "" && e != nil {
		//project exist,the app should exist too
		if e == mongo.ErrNoDocuments {
			log.Error(nil, "[InitSelf] project exist but app not exist", map[string]interface{}{
				"project_id": model.AdminProjectID,
				"group":      model.Group,
				"app":        model.Name})
			e = errors.New("dirty")
		} else {
			log.Error(nil, "[InitSelf] get app failed", map[string]interface{}{
				"project_id": model.AdminProjectID,
				"group":      model.Group,
				"app":        model.Name,
				"error":      e})
		}
		return
	}
	if existProjectIndex.ProjectName != "" {
		//project exist,the app should exist too
		//check secret
		if e = secure.SignCheck(secret, existAppSummary.Value); e != nil {
			log.Error(nil, "[InitSelf] secret check failed", map[string]interface{}{
				"project_id": model.AdminProjectID,
				"group":      model.Group,
				"app":        model.Name,
				"error":      e})
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
		log.Error(nil, "[InitSelf] init project index failed", map[string]interface{}{
			"project_id":   model.AdminProjectID,
			"project_name": model.Project,
			"error":        e})
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
			log.Error(nil, "[InitSelf] project and app not exist but some permission nodes already exist", map[string]interface{}{"node_ids": nodeids, "error": e})
			e = errors.New("dirty")
		} else {
			log.Error(nil, "[InitSelf] init permission nodes failed", map[string]interface{}{"node_ids": nodeids, "error": e})
		}
		return
	}
	//init app
	docs = docs[0:0]
	sign, _ := secure.SignMake(secret)
	//summary
	docs = append(docs, &model.AppSummary{
		ProjectID: model.AdminProjectID,
		Group:     model.Group,
		App:       model.Name,
		Key:       "",
		Index:     0,
		Paths:     map[string]*model.ProxyPath{},
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
			log.Error(nil, "[InitSelf] project not exist but app already exist", map[string]interface{}{
				"project_id": model.AdminProjectID,
				"group":      model.Group,
				"app":        model.Name})
			e = errors.New("dirty")
		} else {
			log.Error(nil, "[InitSelf] init app failed", map[string]interface{}{
				"project_id": model.AdminProjectID,
				"group":      model.Group,
				"app":        model.Name,
				"error":      e})
		}
	}
	return
}
