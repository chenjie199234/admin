package config

import (
	"context"
	"crypto/rand"
	"strconv"
	"time"

	"github.com/chenjie199234/admin/ecode"
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

func (d *Dao) MongoGetAllGroups(ctx context.Context, searchfilter string) ([]string, error) {
	regex := "^config_"
	if searchfilter != "" {
		regex += ".*" + searchfilter + ".*"
	}
	r, e := d.mongo.ListDatabaseNames(ctx, bson.M{"name": bson.M{"$regex": regex}})
	if e != nil {
		return nil, e
	}
	for i := range r {
		r[i] = r[i][7:]
	}
	return r, nil
}
func (d *Dao) MongoGetAllApps(ctx context.Context, groupname, searchfilter string) ([]string, error) {
	return d.mongo.Database("config_"+groupname).ListCollectionNames(ctx, bson.M{"name": bson.M{"$regex": searchfilter}})
}
func (d *Dao) MongoGetPermissionNodeID(ctx context.Context, groupname, appname string) (string, error) {
	appsummary := &model.AppSummary{}
	if e := d.mongo.Database("config_"+groupname).Collection(appname).FindOne(ctx, bson.M{"key": "", "index": 0}, options.FindOne().SetProjection(bson.M{"permission_node_id": 1})).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return "", e
	}
	if appsummary.PermissionNodeID == "" {
		return "", ecode.ErrNotInited
	}
	return appsummary.PermissionNodeID, nil
}
func (d *Dao) MongoCreateApp(ctx context.Context, projectid, groupname, appname, secret string) (e error) {
	if len(secret) >= 32 {
		return ecode.ErrSecretLength
	}
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	parent := &model.Node{}
	if e = d.mongo.Database("permission").Collection("node").FindOneAndUpdate(sctx, bson.M{"node_id": projectid + model.ConfigControl}, bson.M{"$inc": bson.M{"cur_node_index": 1}}).Decode(parent); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrProjectNotExist
		}
		return
	}
	if _, e = d.mongo.Database("permission").Collection("node").InsertOne(sctx, &model.Node{
		NodeId:       parent.NodeId + "," + strconv.FormatUint(uint64(parent.CurNodeIndex+1), 10),
		NodeName:     groupname + "." + appname,
		NodeData:     "",
		CurNodeIndex: 0,
	}); e != nil {
		return
	}
	col := d.mongo.Database("config_" + groupname).Collection(appname)
	index := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "key", Value: 1}, primitive.E{Key: "index", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, e = col.Indexes().CreateOne(sctx, index); e != nil && !mongo.IsDuplicateKeyError(e) {
		return
	}
	nonce := make([]byte, 32)
	rand.Read(nonce)
	if _, e = col.InsertOne(sctx, bson.M{
		"key":                "",
		"index":              0,
		"keys":               bson.M{},
		"value":              util.SignMake(secret, nonce),
		"permission_node_id": parent.NodeId + "," + strconv.FormatUint(uint64(parent.CurNodeIndex+1), 10),
	}); e != nil && mongo.IsDuplicateKeyError(e) {
		e = ecode.ErrAppAlreadyExist
	}
	return
}
func (d *Dao) MongoDelApp(ctx context.Context, groupname, appname, secret string) (e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		} else {
			//drop can't be in the multi collection transaction
			e = d.mongo.Database("config_" + groupname).Collection(appname).Drop(ctx)
		}
	}()
	appsummary := &model.AppSummary{}
	if e = d.mongo.Database("config_"+groupname).Collection(appname).FindOne(sctx, bson.M{"key": "", "index": 0}, options.FindOne().SetProjection(bson.M{"value": 1, "permission_node_id": 1})).Decode(appsummary); e != nil {
		return
	}
	if e = util.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	delfilter := bson.M{"node_id": bson.M{"$regex": "^" + appsummary.PermissionNodeID}}
	if _, e = d.mongo.Database("permission").Collection("node").DeleteMany(sctx, delfilter); e != nil {
		return
	}
	if _, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, delfilter); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, delfilter)
	return
}
func (d *Dao) MongoUpdateAppSecret(ctx context.Context, groupname, appname, oldsecret, newsecret string) (e error) {
	if len(oldsecret) >= 32 || len(newsecret) >= 32 {
		return ecode.ErrSecretLength
	}
	if oldsecret == newsecret {
		return
	}
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	col := d.mongo.Database("config_" + groupname).Collection(appname)
	appsummary := &model.AppSummary{}
	if e = col.FindOne(sctx, bson.M{"key": "", "index": 0}).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	//check oldsecret
	if e = util.SignCheck(oldsecret, appsummary.Value); e != nil {
		return
	}
	nonce := make([]byte, 32)
	rand.Read(nonce)
	updater := bson.M{
		"value": util.SignMake(newsecret, nonce),
	}
	for key, keysummary := range appsummary.Keys {
		if oldsecret != "" {
			var plaintext []byte
			plaintext, e = util.Decrypt(oldsecret, keysummary.CurValue)
			if e != nil {
				return
			}
			keysummary.CurValue = common.Byte2str(plaintext)
		}
		if newsecret != "" {
			updater["keys."+key+".cur_value"], _ = util.Encrypt(newsecret, common.Str2byte(keysummary.CurValue))
		} else {
			updater["keys."+key+".cur_value"] = keysummary.CurValue
		}
	}
	if _, e = col.UpdateOne(sctx, bson.M{"key": "", "index": 0}, bson.M{"$set": updater}); e != nil {
		return
	}
	var cursor *mongo.Cursor
	if cursor, e = col.Find(sctx, bson.M{"key": bson.M{"$exists": true, "$nin": bson.A{nil, ""}}, "index": bson.M{"$gt": 0}}); e != nil {
		return
	}
	defer cursor.Close(sctx)
	for cursor.Next(sctx) {
		log := &model.Log{}
		if e = cursor.Decode(log); e != nil {
			return
		}
		if oldsecret != "" {
			var plaintext []byte
			plaintext, e = util.Decrypt(oldsecret, log.Value)
			if e != nil {
				return
			}
			log.Value = common.Byte2str(plaintext)
		}
		if newsecret != "" {
			log.Value, _ = util.Encrypt(newsecret, common.Str2byte(log.Value))
		}
		if _, e = col.UpdateOne(sctx, bson.M{"key": log.Key, "index": log.Index}, bson.M{"$set": bson.M{"value": log.Value}}); e != nil {
			return
		}
	}
	e = cursor.Err()
	return
}

func (d *Dao) MongoGetAllKeys(ctx context.Context, groupname, appname, secret string) ([]string, error) {
	appsummary := &model.AppSummary{}

	if e := d.mongo.Database("config_"+groupname).Collection(appname).FindOne(ctx, bson.M{"key": "", "index": 0}).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return nil, e
	}
	// check sign
	if e := util.SignCheck(secret, appsummary.Value); e != nil {
		return nil, e
	}
	keys := make([]string, 0, len(appsummary.Keys))
	for k := range appsummary.Keys {
		keys = append(keys, k)
	}
	return keys, nil
}

// index == 0 get the current index's config
// index != 0 get the specific index's config
func (d *Dao) MongoGetKeyConfig(ctx context.Context, groupname, appname, key string, index uint32, secret string) (*model.KeySummary, *model.Log, error) {
	col := d.mongo.Database("config_"+groupname, options.Database().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())).Collection(appname)
	var appsummary *model.AppSummary
	var log *model.Log
	if index == 0 {
		//get tge current index's config
		appsummary = &model.AppSummary{}
		if e := col.FindOne(ctx, bson.M{"key": "", "index": 0}, options.FindOne().SetProjection(bson.M{"value": 1, "keys." + key: 1})).Decode(appsummary); e != nil {
			if e == mongo.ErrNoDocuments {
				e = ecode.ErrAppNotExist
			}
			return nil, nil, e
		}
		if appsummary.Keys == nil {
			return nil, nil, ecode.ErrKeyNotExist
		}
		keysummary, ok := appsummary.Keys[key]
		if !ok {
			return nil, nil, ecode.ErrKeyNotExist
		}
		//check secret
		if e := util.SignCheck(secret, appsummary.Value); e != nil {
			return nil, nil, e
		}
		if secret != "" {
			plaintext, e := util.Decrypt(secret, keysummary.CurValue)
			if e != nil {
				return nil, nil, e
			}
			keysummary.CurValue = common.Byte2str(plaintext)
		}
		log = &model.Log{
			Key:       key,
			Index:     keysummary.CurIndex,
			Value:     keysummary.CurValue,
			ValueType: keysummary.CurValueType,
		}
		return keysummary, log, nil
	}
	//get the specific index's config and the current status
	filter := bson.M{"$or": bson.A{bson.M{"key": "", "index": 0}, bson.M{"key": key, "index": index}}}
	cursor, e := col.Find(ctx, filter, options.Find().SetProjection(bson.M{"key": 1, "index": 1, "value": 1, "value_type": 1, "keys." + key: 1}).SetSort(bson.M{"index": 1}))
	if e != nil {
		return nil, nil, e
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		if appsummary == nil {
			tmp := &model.AppSummary{}
			if e = cursor.Decode(tmp); e != nil {
				return nil, nil, e
			}
			appsummary = tmp
		} else {
			tmp := &model.Log{}
			if e = cursor.Decode(tmp); e != nil {
				return nil, nil, e
			}
			log = tmp
		}
	}
	if e := cursor.Err(); e != nil {
		return nil, nil, e
	}
	if appsummary == nil {
		return nil, nil, ecode.ErrAppNotExist
	}
	if appsummary.Keys == nil {
		return nil, nil, ecode.ErrKeyNotExist
	}
	keysummary, ok := appsummary.Keys[key]
	if !ok {
		return nil, nil, ecode.ErrKeyNotExist
	}
	if log == nil {
		return nil, nil, ecode.ErrIndexNotExist
	}
	//check secret
	if e := util.SignCheck(secret, appsummary.Value); e != nil {
		return nil, nil, e
	}
	if secret != "" {
		plaintext, e := util.Decrypt(secret, keysummary.CurValue)
		if e != nil {
			return nil, nil, e
		}
		keysummary.CurValue = common.Byte2str(plaintext)
		plaintext, e = util.Decrypt(secret, log.Value)
		if e != nil {
			return nil, nil, e
		}
		log.Value = common.Byte2str(plaintext)
	}
	return keysummary, log, nil
}
func (d *Dao) MongoSetKeyConfig(ctx context.Context, groupname, appname, key, secret, value, valuetype string) (newindex, newversion uint32, e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	col := d.mongo.Database("config_" + groupname).Collection(appname)
	appsummary := &model.AppSummary{}
	if e = col.FindOne(sctx, bson.M{"key": "", "index": 0}).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	//check secret
	if e = util.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	if secret != "" {
		if value, e = util.Encrypt(secret, common.Str2byte(value)); e != nil {
			return
		}
	}
	keysummary, ok := appsummary.Keys[key]
	if !ok {
		keysummary = &model.KeySummary{
			CurIndex:   0,
			MaxIndex:   0,
			CurVersion: 0,
			CurValue:   "",
		}
	}
	keysummary.MaxIndex += 1
	keysummary.CurIndex = keysummary.MaxIndex
	keysummary.CurVersion += 1
	keysummary.CurValue = value
	keysummary.CurValueType = valuetype
	if _, e = col.UpdateOne(sctx, bson.M{"key": "", "index": 0}, bson.M{"$set": bson.M{"keys." + key: keysummary}}); e != nil {
		return
	}
	if _, e = col.UpdateOne(sctx, bson.M{"key": key, "index": keysummary.CurIndex}, bson.M{"$set": bson.M{"value": value, "value_type": valuetype}}, options.Update().SetUpsert(true)); e != nil {
		return
	}
	newindex = keysummary.CurIndex
	newversion = keysummary.CurVersion
	return
}
func (d *Dao) MongoDelKey(ctx context.Context, groupname, appname, key, secret string) (e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	col := d.mongo.Database("config_" + groupname).Collection(appname)
	appsummary := &model.AppSummary{}
	if e = col.FindOneAndUpdate(sctx, bson.M{"key": "", "index": 0}, bson.M{"$unset": bson.M{"keys." + key: 1}}, options.FindOneAndUpdate().SetProjection(bson.M{"value": 1})).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = util.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	_, e = col.DeleteMany(sctx, bson.M{"key": key})
	return
}
func (d *Dao) MongoRollbackKeyConfig(ctx context.Context, groupname, appname, key, secret string, index uint32) (e error) {
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	col := d.mongo.Database("config_" + groupname).Collection(appname)
	log := &model.Log{}
	if e = col.FindOne(sctx, bson.M{"key": key, "index": index}).Decode(log); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrIndexNotExist
		}
		return
	}
	updateSummary := bson.M{
		"$set": bson.M{
			"keys." + key + ".cur_index":      index,
			"keys." + key + ".cur_value":      log.Value,
			"keys." + key + ".cur_value_type": log.ValueType,
		},
		"$inc": bson.M{
			"keys." + key + ".cur_version": 1,
		},
	}
	appsummary := &model.AppSummary{}
	if e = col.FindOneAndUpdate(sctx, bson.M{"key": "", "index": 0}, updateSummary, options.FindOneAndUpdate().SetProjection(bson.M{"value": 1})).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	e = util.SignCheck(secret, appsummary.Value)
	return
}

// first key groupname,second key appname,value curconfig
type WatchUpdateHandler func(string, string, *model.AppSummary)
type WatchDeleteAppHandler func(groupname, appname string)
type WatchDeleteConfigHandler func(groupname, appname string, id string)

func (d *Dao) MongoWatchConfig(update WatchUpdateHandler, delA WatchDeleteAppHandler, delC WatchDeleteConfigHandler) error {
	starttime := &primitive.Timestamp{T: uint32(time.Now().Unix()), I: uint32(0)}
	watchfilter := mongo.Pipeline{bson.D{primitive.E{Key: "$match", Value: bson.M{"ns.db": bson.M{"$regex": "^config_"}}}}}
	stream, e := d.mongo.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup).SetStartAtOperationTime(starttime))
	if e != nil {
		return e
	}
	go func() {
		for {
			for stream == nil {
				//reconnect
				time.Sleep(time.Millisecond * 5)
				if stream, e = d.mongo.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup).SetStartAtOperationTime(starttime)); e != nil {
					log.Error(nil, "[dao.MongoWatchConfig] reconnect stream error:", e)
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
					groupname := stream.Current.Lookup("ns").Document().Lookup("db").StringValue()[7:]
					appname := stream.Current.Lookup("ns").Document().Lookup("coll").StringValue()
					delA(groupname, appname)
				case "insert":
					//insert document
					fallthrough
				case "update":
					//update document
					groupname := stream.Current.Lookup("ns").Document().Lookup("db").StringValue()[7:]
					appname := stream.Current.Lookup("ns").Document().Lookup("coll").StringValue()
					key, ok1 := stream.Current.Lookup("fullDocument").Document().Lookup("key").StringValueOK()
					index, ok2 := stream.Current.Lookup("fullDocument").Document().Lookup("index").Int32OK()
					if !ok1 || !ok2 {
						//unknown doc
						continue
					}
					if key != "" || index != 0 {
						//this is not the app summary
						continue
					}
					//this is the app summary
					s := &model.AppSummary{}
					if e := stream.Current.Lookup("fullDocument").Unmarshal(s); e != nil {
						log.Error(nil, "[dao.MongoWatchConfig] group:", groupname, "app:", appname, "summary data broken:", e)
						continue
					}
					update(groupname, appname, s)
				case "delete":
					//delete document
					groupname := stream.Current.Lookup("ns").Document().Lookup("db").StringValue()[7:]
					appname := stream.Current.Lookup("ns").Document().Lookup("coll").StringValue()
					id := stream.Current.Lookup("documentKey").Document().Lookup("_id").ObjectID().Hex()
					delC(groupname, appname, id)
				}
			}
			if stream.Err() != nil {
				log.Error(nil, "[dao.MongoWatchConfig]", stream.Err())
			}
			stream.Close(context.Background())
			stream = nil
		}
	}()
	return nil
}

// this function will not decrypt
func (d *Dao) MongoGetAppConfig(ctx context.Context, groupname, appname string) (*model.AppSummary, error) {
	app := &model.AppSummary{}
	if e := d.mongo.Database("config_"+groupname).Collection(appname).FindOne(ctx, bson.M{"key": "", "index": 0}).Decode(app); e != nil {
		if e == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, e
	}
	return app, nil
}
