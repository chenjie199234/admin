package internal

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
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
}

// keyvalue: map's key is the key name,map's value is the key's data
// keytype: map's key is the key name,map's value is the type of the key's data
type NoticeHandler func(key, keyvalue, keytype string)

// url format [mongodb/mongodb+srv]://[username:password@]host1,...,hostN/[dbname][?param1=value1&...&paramN=valueN]
func NewDirectSdk(selfgroup, selfname, mongourl, secret string, AppConfigTemplate, SourceConfigTemplate []byte) (*Sdk, error) {
	client, e := newMongo(mongourl, selfgroup, selfname, secret, AppConfigTemplate, SourceConfigTemplate)
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
	if e = col.FindOne(context.Background(), bson.M{"key": "", "index": 0}).Decode(appsummary); e != nil && e != mongo.ErrNoDocuments {
		return nil, e
	} else if e == nil {
		//sign check
		if e := util.SignCheck(secret, appsummary.Value); e != nil {
			return nil, e
		}
		if secret != "" {
			for _, keysummary := range appsummary.Keys {
				plaintext, e := util.Decrypt(secret, keysummary.CurValue)
				if e != nil {
					return nil, e
				}
				keysummary.CurValue = common.Byte2str(plaintext)
			}
		}
	}
	instance := &Sdk{
		appsummary: appsummary,
		notices:    make(map[string]NoticeHandler),
	}
	go func() {
		for {
			for stream == nil {
				//reconnect
				time.Sleep(time.Millisecond * 100)
				if stream, e = client.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup).SetStartAtOperationTime(starttime)); e != nil {
					log.Error(nil, "[config.directsdk.watch] reconnect:", e)
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
						notice(key, "", "raw")
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
						log.Error(nil, "[config.directsdk.watch] group:", selfgroup, "app:", selfname, e)
						continue
					}
					//check sign
					if e := util.SignCheck(secret, tmp.Value); e != nil {
						log.Error(nil, "[config.directsdk.watch] group:", selfgroup, "app:", selfname, e)
						continue
					}
					if secret != "" {
						for _, keysummary := range tmp.Keys {
							plaintext, e := util.Decrypt(secret, keysummary.CurValue)
							if e != nil {
								log.Error(nil, "[config.directsdk.watch] group:", selfgroup, "app:", selfname, e)
								break
							}
							keysummary.CurValue = common.Byte2str(plaintext)
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
					if appsummary.ID.IsZero() || stream.Current.Lookup("documentKey").Document().Lookup("_id").ObjectID().Hex() != appsummary.ID.Hex() {
						//this is not the appsummary
						continue
					}
					instance.lker.Lock()
					for key, notice := range instance.notices {
						notice(key, "", "raw")
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

func newMongo(url, groupname, appname, secret string, AppConfigTemplate, SourceConfigTemplate []byte) (db *mongo.Client, e error) {
	op := options.Client().ApplyURI(url)
	op = op.SetMaxPoolSize(2)
	op = op.SetConnectTimeout(time.Second)
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
	//init self mongo
	col := db.Database("config_" + groupname).Collection(appname)
	if _, e = col.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: bson.D{{Key: "key", Value: 1}, {Key: "index", Value: 1}}, Options: options.Index().SetUnique(true)}); e != nil && !mongo.IsDuplicateKeyError(e) {
		return
	}
	bufapp := bytes.NewBuffer(nil)
	if e = json.Compact(bufapp, AppConfigTemplate); e != nil {
		return
	}
	bufsource := bytes.NewBuffer(nil)
	SourceConfigTemplate = bytes.ReplaceAll(SourceConfigTemplate, []byte("example_mongo"), []byte("admin_mongo"))
	SourceConfigTemplate = bytes.ReplaceAll(SourceConfigTemplate, []byte("[mongodb/mongodb+srv]://[username:password@]host1,...,hostN[/dbname][?param1=value1&...&paramN=valueN]"), []byte(url))
	if e = json.Compact(bufsource, SourceConfigTemplate); e != nil {
		return
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
	if _, e = col.InsertOne(context.Background(), bson.M{
		"key":   "",
		"index": 0,
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
		"value": util.SignMake(secret, nonce),
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
