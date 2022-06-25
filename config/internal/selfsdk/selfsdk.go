package selfsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	lker       sync.Mutex
	appsummary *model.AppSummary
	notices    map[string]NoticeHandler
}

//keyvalue: map's key is the key name,map's value is the key's data
//keytype: map's key is the key name,map's value is the type of the key's data
type NoticeHandler func(key, keyvalue, keytype string)

//url format [mongodb/mongodb+srv]://[username:password@]host1,...,hostN/[dbname][?param1=value1&...&paramN=valueN]
func NewDirectSdk(selfgroup, selfname string, mongodburl string) (*Sdk, error) {
	client, e := newMongo(mongodburl, selfgroup, selfname)
	if e != nil {
		return nil, e
	}
	starttime := &primitive.Timestamp{T: uint32(time.Now().Unix()), I: 0}
	watchfilter := mongo.Pipeline{bson.D{bson.E{Key: "$match", Value: bson.M{"ns.db": "config_" + selfgroup, "ns.coll": selfname}}}}
	stream, e := client.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup).SetStartAtOperationTime(starttime))
	if e != nil {
		return nil, e
	}
	col := client.Database("config_"+selfgroup, options.Database().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())).Collection(selfname)
	//get first,then watch change stream
	appsummary := &model.AppSummary{}
	if e = col.FindOne(context.Background(), bson.M{"index": 0}).Decode(appsummary); e != nil && e != mongo.ErrNoDocuments {
		return nil, e
	}
	if appsummary.Cipher != "" {
		for _, keysummary := range appsummary.Keys {
			keysummary.CurValue = util.Decrypt(appsummary.Cipher, keysummary.CurValue)
		}
	}
	instance := &Sdk{
		client:     client,
		appsummary: appsummary,
		notices:    make(map[string]NoticeHandler),
	}
	go func() {
		for {
			for stream == nil {
				//reconnect
				time.Sleep(time.Millisecond * 100)
				if stream, e = client.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup).SetStartAtOperationTime(starttime)); e != nil {
					log.Error(nil, "[config.selfsdk.watch] reconnect error:", e)
					stream = nil
					continue
				}
			}
			for stream.Next(context.Background()) {
				starttime.T, starttime.I = stream.Current.Lookup("clusterTime").Timestamp()
				starttime.I++
				switch stream.Current.Lookup("operationType").StringValue() {
				case "drop":
					//drop collection
					instance.lker.Lock()
					for key, notice := range instance.notices {
						old, ok := instance.appsummary.Keys[key]
						if ok && old.CurVersion != 0 {
							notice(key, "", "raw")
						}
					}
					instance.appsummary = &model.AppSummary{}
					instance.lker.Unlock()
				case "insert":
					//insert
					fallthrough
				case "update":
					//update
					key, kok := stream.Current.Lookup("fullDocument").Document().Lookup("key").StringValueOK()
					index, iok := stream.Current.Lookup("fullDocument").Document().Lookup("index").Int32OK()
					if !kok || !iok {
						//unknown doc
						continue
					}
					if key != "" || index != 0 {
						//this is not the appsummary
						continue
					}
					tmp := &model.AppSummary{}
					if e := stream.Current.Lookup("fullDocument").Unmarshal(tmp); e != nil {
						log.Error(nil, "[config.selfsdk.watch] group:", selfgroup, "app:", selfname, "appsummary data broken,error:", e)
						continue
					}
					if tmp.Cipher != "" {
						for _, keysummary := range tmp.Keys {
							keysummary.CurValue = util.Decrypt(tmp.Cipher, keysummary.CurValue)
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
						} else {
							if !newok {
								continue
							} else {
								notice(key, new.CurValue, new.CurValueType)
							}
						}
					}
					instance.appsummary = tmp
					instance.lker.Unlock()
				case "delete":
					if stream.Current.Lookup("documentKey").Document().Lookup("_id").ObjectID().Hex() != appsummary.ID.Hex() {
						//this is not the appsummary
						continue
					}
					instance.lker.Lock()
					for key, notice := range instance.notices {
						old, ok := instance.appsummary.Keys[key]
						if ok && old.CurVersion != 0 {
							notice(key, "", "raw")
						}
					}
					appsummary = &model.AppSummary{}
					instance.lker.Unlock()
				}
			}
			if stream.Err() != nil {
				log.Error(nil, "[config.selfsdk.watch] error:", stream.Err())
			}
			stream.Close(context.Background())
			stream = nil
		}
	}()
	return instance, nil
}

//watch the same key will overwrite the old one's notice function
//but the old's cancel function can still work
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
	}
	return func() {
		instance.lker.Lock()
		delete(instance.notices, key)
		instance.lker.Unlock()
	}
}

var defaultAppConfig = `{
	"handler_timeout":{
		"/admin.status/ping":{
			"GET":"200ms",
			"CRPC":"200ms",
			"GRPC":"200ms"
		},
		"/admin.config/watch":{
			"POST":"0s"
		},
		"/admin.config/updatechiper":{
			"POST":"0s"
		}
	},
	"handler_rate":[{
		"Path":"/admin.status/ping",
		"Method":["GET","GRPC","CRPC"],
		"MaxPerSec":10
	}],
	"white_ip":[],
	"black_ip":[],
	"web_path_rewrite":{
		"GET":{
			"/example/origin/url":"/example/new/url"
		}
	},
	"access_keys":{
		"default":["default_sec_key"],
		"/admin.status/ping":["specific_sec_key"]
	},
	"token_secret":"test",
	"token_expire":"24h",
	"service":{

	}
}`
var defaultSourceConfig = `{
	"cgrpc_server":{
		"connect_timeout":"200ms",
		"global_timeout":"200ms",
		"heart_probe":"1.5s"
	},
	"cgrpc_client":{
		"connect_timeout":"200ms",
		"global_timeout":"0",
		"heart_probe":"1.5s"
	},
	"crpc_server":{
		"connect_timeout":"200ms",
		"global_timeout":"200ms",
		"heart_probe":"1.5s"
	},
	"crpc_client":{
		"connect_timeout":"200ms",
		"global_timeout":"0",
		"heart_probe":"1.5s"
	},
	"web_server":{
		"close_mode":1,
		"connect_timeout":"200ms",
		"global_timeout":"200ms",
		"idle_timeout":"5s",
		"heart_probe":"1.5s",
		"static_file":"./src",
		"web_cors":{
			"cors_origin":["*"],
			"cors_header":["*"],
			"cors_expose":[]
		}
	},
	"web_client":{
		"connect_timeout":"200ms",
		"global_timeout":"0",
		"idle_timeout":"5s",
		"heart_probe":"1.5s"
	},
	"mongo":{
		"admin_mongo":{
			"url":"%s",
			"max_open":100,
			"max_idletime":"10m",
			"io_timeout":"500ms",
			"conn_timeout":"500ms"
		}
	},
	"sql":{
		"example_sql":{
			"url":"[username:password@][protocol(address)][/dbname][?param1=value1&...&paramN=valueN]",
			"max_open":100,
			"max_idletime":"10m",
			"io_timeout":"200ms",
			"conn_timeout":"200ms"
		}
	},
	"redis":{
		"example_redis":{
			"url":"[redis/rediss]://[[username:]password@]host[/dbindex]",
			"max_open":100,
			"max_idletime":"10m",
			"io_timeout":"200ms",
			"conn_timeout":"200ms"
		}
	},
	"kafka_pub":[
		{
			"addrs":["127.0.0.1:12345"],
			"username":"example",
			"password":"example",
			"auth_method":3,
			"compress_method":2,
			"topic_name":"example_topic",
			"io_timeout":"500ms",
			"conn_timeout":"200ms"
		}
	],
	"kafka_sub":[
		{
			"addrs":["127.0.0.1:12345"],
			"username":"example",
			"password":"example",
			"auth_method":3,
			"topic_name":"example_topic",
			"group_name":"example_group",
			"conn_timeout":"200ms",
			"start_offset":-2,
			"commit_interval":"0s"
		}
	]
}`

func newMongo(url string, groupname, appname string) (db *mongo.Client, e error) {
	op := options.Client().ApplyURI(url)
	op = op.SetMaxPoolSize(2)
	op = op.SetHeartbeatInterval(time.Second * 5)
	op = op.SetReadPreference(readpref.SecondaryPreferred())
	op = op.SetReadConcern(readconcern.Majority())
	if db, e = mongo.Connect(context.Background(), op); e != nil {
		return
	}
	if e = db.Ping(context.Background(), nil); e != nil {
		return nil, e
	}
	//init self mongo
	var s mongo.Session
	if s, e = db.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local())); e != nil {
		return
	}
	defer s.EndSession(context.Background())
	sctx := mongo.NewSessionContext(context.Background(), s)
	if e = s.StartTransaction(); e != nil {
		return nil, e
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
			if mongo.IsDuplicateKeyError(e) {
				e = nil
			}
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	col := db.Database("config_" + groupname).Collection(appname)
	index := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "key", Value: 1}, primitive.E{Key: "index", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, e = col.Indexes().CreateOne(sctx, index); e != nil {
		return
	}
	buf := bytes.NewBuffer(nil)
	if e = json.Compact(buf, common.Str2byte(defaultAppConfig)); e != nil {
		return
	}
	appconfig := buf.String()
	buf.Reset()
	if e = json.Compact(buf, common.Str2byte(fmt.Sprintf(defaultSourceConfig, url))); e != nil {
		return
	}
	sourceconfig := buf.String()
	if _, e = col.InsertOne(sctx, bson.M{
		"key":    "",
		"index":  0,
		"cipher": "",
		"keys": map[string]*model.KeySummary{
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
	}); e != nil {
		return
	}
	if _, e = col.UpdateOne(sctx, bson.M{"key": "AppConfig", "index": 1}, bson.M{"$set": bson.M{"value": appconfig, "value_type": "json"}}, options.Update().SetUpsert(true)); e != nil {
		return
	}
	if _, e = col.UpdateOne(sctx, bson.M{"key": "SourceConfig", "index": 1}, bson.M{"$set": bson.M{"value": sourceconfig, "value_type": "json"}}, options.Update().SetUpsert(true)); e != nil {
		return
	}
	return db, nil
}
