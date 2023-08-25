package internal

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Sdk struct {
	client     *mongo.Client
	secret     string
	lker       sync.Mutex
	appsummary *model.AppSummary
	notices    map[string]NoticeHandler
	start      *primitive.Timestamp
}

// keyvalue: map's key is the key name,map's value is the key's data
// keytype: map's key is the key name,map's value is the type of the key's data
type NoticeHandler func(key, keyvalue, keytype string)

var (
	ErrMissingEnvMongo = errors.New("env REMOTE_CONFIG_MONGO_URL missing")
	ErrWrongEnvSecret  = errors.New("env REMOTE_CONFIG_SECRET too long")
)

// must set below env:
// REMOTE_CONFIG_MONGO_URL,format [mongodb/mongodb+srv]://[username:password@]host1,...,hostN/[dbname][?param1=value1&...&paramN=valueN]
// REMOTE_CONFIG_SECRET
func NewDirectSdk(AppConfigTemplate, SourceConfigTemplate []byte) (*Sdk, error) {
	mongourl, secret, e := env()
	if e != nil {
		return nil, e
	}
	client, e := newMongo(mongourl)
	if e != nil {
		return nil, e
	}
	if e = initself(client, mongourl, secret, AppConfigTemplate, SourceConfigTemplate); e != nil {
		return nil, e
	}
	instance := &Sdk{
		client:     client,
		secret:     secret,
		appsummary: &model.AppSummary{},
		notices:    make(map[string]NoticeHandler),
		start:      &primitive.Timestamp{T: uint32(time.Now().Unix() - 1), I: 0},
	}
	instance.first()
	go instance.watch()
	return instance, nil
}
func env() (mongourl string, secret string, e error) {
	if str, ok := os.LookupEnv("REMOTE_CONFIG_MONGO_URL"); ok && str != "<REMOTE_CONFIG_MONGO_URL>" && str != "" {
		mongourl = str
	} else {
		return "", "", ErrMissingEnvMongo
	}
	if str, ok := os.LookupEnv("REMOTE_CONFIG_SECRET"); ok && str != "<REMOTE_CONFIG_SECRET>" && str != "" {
		secret = str
	}
	if len(secret) >= 32 {
		return "", "", ErrWrongEnvSecret
	}
	return
}
func (instance *Sdk) first() error {
	if e := instance.client.Database("app").Collection("config").FindOne(context.Background(), bson.M{"project_id": model.AdminProjectID, "group": model.Group, "app": model.Name, "key": "", "index": 0}).Decode(instance.appsummary); e != nil {
		return e
	}
	//sign check
	if e := util.SignCheck(instance.secret, instance.appsummary.Value); e != nil {
		return e
	}
	if instance.secret == "" {
		return nil
	}
	for _, keysummary := range instance.appsummary.Keys {
		plaintext, e := util.Decrypt(instance.secret, keysummary.CurValue)
		if e != nil {
			return e
		}
		keysummary.CurValue = common.Byte2str(plaintext)
	}
	return nil
}
func (instance *Sdk) watch() {
	watchfilter := mongo.Pipeline{bson.D{bson.E{Key: "$match", Value: bson.M{"ns.db": "app", "ns.coll": "config"}}}}
	var stream *mongo.ChangeStream
	for {
		for stream == nil {
			//connect
			var e error
			if stream, e = instance.client.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup).SetStartAtOperationTime(instance.start)); e != nil {
				log.Error(nil, "[config.directsdk.watch] connect failed", map[string]interface{}{"error": e})
				stream = nil
				time.Sleep(time.Millisecond * 100)
				continue
			}
		}
		for stream.Next(context.Background()) {
			instance.start.T, instance.start.I = stream.Current.Lookup("clusterTime").Timestamp()
			instance.start.I++
			switch stream.Current.Lookup("operationType").StringValue() {
			case "drop":
				//drop collection
				log.Error(nil, "[config.directsdk.watch] all configs deleted", nil)
				instance.lker.Lock()
				for key, notice := range instance.notices {
					notice(key, "", "raw")
				}
				instance.appsummary = &model.AppSummary{}
				instance.lker.Unlock()
			case "insert":
				//insert
				fallthrough
			case "update":
				//update
				projectid, pok := stream.Current.Lookup("fullDocument").Document().Lookup("project_id").StringValueOK()
				group, gok := stream.Current.Lookup("fullDocument").Document().Lookup("group").StringValueOK()
				app, aok := stream.Current.Lookup("fullDocument").Document().Lookup("app").StringValueOK()
				key, kok := stream.Current.Lookup("fullDocument").Document().Lookup("key").StringValueOK()
				index, iok := stream.Current.Lookup("fullDocument").Document().Lookup("index").AsInt64OK()
				if !pok || !gok || !aok || !kok || !iok {
					//unknown doc
					continue
				}
				if projectid != model.AdminProjectID || group != model.Group || app != model.Name || key != "" || index != 0 {
					//this is not the needed appsummary
					continue
				}
				tmp := &model.AppSummary{}
				if e := stream.Current.Lookup("fullDocument").Unmarshal(tmp); e != nil {
					log.Error(nil, "[config.directsdk.watch] document format wrong", map[string]interface{}{"project_id": model.AdminProjectID, "group": model.Group, "app": model.Name, "error": e})
					continue
				}
				//check sign
				if e := util.SignCheck(instance.secret, tmp.Value); e != nil {
					log.Error(nil, "[config.directsdk.watch] sign check failed", map[string]interface{}{"project_id": model.AdminProjectID, "group": model.Group, "app": model.Name, "error": e})
					continue
				}
				if instance.secret != "" {
					success := true
					for _, keysummary := range tmp.Keys {
						plaintext, e := util.Decrypt(instance.secret, keysummary.CurValue)
						if e != nil {
							success = false
							log.Error(nil, "[config.directsdk.watch] decrypt failed", map[string]interface{}{"project_id": model.AdminProjectID, "group": model.Group, "app": model.Name, "error": e})
							break
						}
						keysummary.CurValue = common.Byte2str(plaintext)
					}
					if !success {
						continue
					}
				}
				instance.lker.Lock()
				for key, notice := range instance.notices {
					old, oldok := instance.appsummary.Keys[key]
					new, newok := tmp.Keys[key]
					if oldok {
						if !newok {
							notice(key, "", "raw")
						} else if old.CurVersion == new.CurVersion {
							continue
						} else {
							notice(key, new.CurValue, new.CurValueType)
						}
					} else if !newok {
						continue
					} else {
						notice(key, new.CurValue, new.CurValueType)
					}
				}
				instance.appsummary = tmp
				instance.lker.Unlock()
			case "delete":
				if instance.appsummary.ID.IsZero() || stream.Current.Lookup("documentKey").Document().Lookup("_id").ObjectID().Hex() != instance.appsummary.ID.Hex() {
					//this is not the needed appsummary
					continue
				}
				//this is same as delete the app
				log.Error(nil, "[config.directsdk.watch] app deleted", map[string]interface{}{"project_id": model.AdminProjectID, "group": model.Group, "app": model.Name})
				instance.lker.Lock()
				for key, notice := range instance.notices {
					notice(key, "", "raw")
				}
				instance.appsummary = &model.AppSummary{}
				instance.lker.Unlock()
			}
		}
		if stream.Err() != nil {
			log.Error(nil, "[config.selfsdk.watch] stream disconnected", map[string]interface{}{"error": stream.Err()})
		}
		stream.Close(context.Background())
		stream = nil
	}
}

// watch the same key will overwrite the old one's notice function
// but the old's cancel function can still work
func (instance *Sdk) Watch(key string, notice NoticeHandler) (cancel func()) {
	instance.lker.Lock()
	defer instance.lker.Unlock()
	if _, ok := instance.notices[key]; ok {
		instance.notices[key] = notice
		return func() {
			instance.lker.Lock()
			delete(instance.notices, key)
			instance.lker.Unlock()
		}
	}
	instance.notices[key] = notice
	if keysummary, ok := instance.appsummary.Keys[key]; ok {
		go notice(key, keysummary.CurValue, keysummary.CurValueType)
	} else {
		go notice(key, "", "raw")
	}
	return func() {
		instance.lker.Lock()
		delete(instance.notices, key)
		instance.lker.Unlock()
	}
}

func newMongo(url string) (db *mongo.Client, e error) {
	op := options.Client().ApplyURI(url)
	op = op.SetMaxPoolSize(3)
	op = op.SetConnectTimeout(time.Second * 3)
	op = op.SetTimeout(time.Second * 10)
	op = op.SetHeartbeatInterval(time.Second * 5)
	op = op.SetMaxConnIdleTime(time.Minute)
	op = op.SetReadPreference(readpref.SecondaryPreferred())
	op = op.SetReadConcern(readconcern.Majority())
	if db, e = mongo.Connect(context.Background(), op); e != nil {
		return
	}
	if e = db.Ping(context.Background(), nil); e != nil {
		return nil, e
	}
	return db, nil
}
func initself(db *mongo.Client, mongourl, secret string, AppConfigTemplate, SourceConfigTemplate []byte) (e error) {
	bufapp := bytes.NewBuffer(nil)
	if e = json.Compact(bufapp, AppConfigTemplate); e != nil {
		return
	}
	bufsource := bytes.NewBuffer(nil)
	SourceConfigTemplate = bytes.ReplaceAll(SourceConfigTemplate, []byte("example_mongo"), []byte("admin_mongo"))
	SourceConfigTemplate = bytes.ReplaceAll(SourceConfigTemplate, []byte("[mongodb/mongodb+srv]://[username:password@]host1,...,hostN[/dbname][?param1=value1&...&paramN=valueN]"), []byte(mongourl))
	if e = json.Compact(bufsource, SourceConfigTemplate); e != nil {
		return e
	}
	appconfig := ""
	sourceconfig := ""
	if secret != "" {
		appconfig, _ = util.Encrypt(secret, bufapp.Bytes())
		sourceconfig, _ = util.Encrypt(secret, bufsource.Bytes())
	} else {
		appconfig = common.Byte2str(bufapp.Bytes())
		sourceconfig = common.Byte2str(bufsource.Bytes())
	}
	nonce := make([]byte, 32)
	rand.Read(nonce)

	var s mongo.Session
	if s, e = db.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local())); e != nil {
		return e
	}
	defer s.EndSession(context.Background())
	sctx := mongo.NewSessionContext(context.Background(), s)
	if e = s.StartTransaction(); e != nil {
		return e
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()

	//check init status

	//check project index
	existProjectIndex := &model.ProjectIndex{}
	if e = db.Database("permission").Collection("projectindex").FindOne(sctx, bson.M{"project_id": model.AdminProjectID}).Decode(existProjectIndex); e != nil && e != mongo.ErrNoDocuments {
		return
	} else if e == nil && existProjectIndex.ProjectName != model.Project {
		e = errors.New("init conflict:already inited with other project name")
		return
	}
	//check app
	existAppSummary := &model.AppSummary{}
	e = db.Database("app").Collection("config").FindOne(sctx, bson.M{"project_id": model.AdminProjectID, "group": model.Group, "app": model.Name, "key": "", "index": 0}).Decode(existAppSummary)
	if existProjectIndex.ProjectName == "" && e != mongo.ErrNoDocuments {
		//project not exist,the app should not exist too
		if e == nil {
			e = errors.New("init conflict:db data dirty")
		}
		return
	}
	if existProjectIndex.ProjectName != "" && e != nil {
		//project exist,the app should exist too
		if e == mongo.ErrNoDocuments {
			e = errors.New("init conflict:db data dirty")
		}
		return
	}
	if existProjectIndex.ProjectName != "" {
		//project exist,the app should exist too
		//check secret
		if e = util.SignCheck(secret, existAppSummary.Value); e != nil {
			return
		}
	}
	if existProjectIndex.ProjectName != "" {
		return
	}

	//init now

	//init project index
	if _, e = db.Database("permission").Collection("projectindex").InsertOne(sctx, bson.M{"project_name": model.Project, "project_id": model.AdminProjectID}); e != nil {
		return
	}
	//init node
	docs := bson.A{}
	docs = append(docs, model.Node{
		NodeId:       "0",
		NodeName:     "root",
		NodeData:     "",
		CurNodeIndex: 1,
	})
	//first project's node
	docs = append(docs, &model.Node{
		NodeId:       model.AdminProjectID,
		NodeName:     model.Project,
		NodeData:     "",
		CurNodeIndex: 100,
	})
	//first project's user and role control node
	docs = append(docs, &model.Node{
		NodeId:       model.AdminProjectID + model.UserAndRoleControl,
		NodeName:     "UserAndRoleControl",
		NodeData:     "",
		CurNodeIndex: 0,
	})
	//first project's config control node
	docs = append(docs, &model.Node{
		NodeId:       model.AdminProjectID + model.AppControl,
		NodeName:     "AppControl",
		NodeData:     "",
		CurNodeIndex: 1,
	})
	//first project's first app(this app: admin)'s config node
	docs = append(docs, &model.Node{
		NodeId:       model.AdminProjectID + model.AppControl + ",1",
		NodeName:     model.Group + "." + model.Name,
		NodeData:     "",
		CurNodeIndex: 0,
	})
	if _, e = db.Database("permission").Collection("node").InsertMany(sctx, docs); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			e = errors.New("init conflict:permission node already exist")
		}
		return
	}
	//init app
	docs = docs[0:0]
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
		Value:            util.SignMake(secret, nonce),
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
			e = errors.New("init conflict:config already exist")
		}
	}
	return
}
