package config

import (
	"context"
	"time"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/log"
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
func (d *Dao) MongoDelGroup(ctx context.Context, groupname string) error {
	return d.mongo.Database("config_" + groupname).Drop(ctx)
}
func (d *Dao) MongoGetAllApps(ctx context.Context, groupname, searchfilter string) ([]string, error) {
	return d.mongo.Database("config_"+groupname).ListCollectionNames(ctx, bson.M{"name": bson.M{"$regex": searchfilter}})
}
func (d *Dao) MongoDelApp(ctx context.Context, groupname, appname string) error {
	return d.mongo.Database("config_" + groupname).Collection(appname).Drop(ctx)
}
func (d *Dao) MongoGetAllKeys(ctx context.Context, groupname, appname string) ([]string, error) {
	appsummary := &model.AppSummary{}
	e := d.mongo.Database("config_"+groupname).Collection(appname).FindOne(ctx, bson.M{"key": "", "index": 0}).Decode(appsummary)
	if e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return nil, e
	}
	result := make([]string, 0, len(appsummary.Keys))
	for k := range appsummary.Keys {
		result = append(result, k)
	}
	return result, nil
}
func (d *Dao) MongoDelKey(ctx context.Context, groupname, appname, key string) (e error) {
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
	_, e = d.mongo.Database("config_"+groupname).Collection(appname).UpdateOne(sctx, bson.M{"key": "", "index": 0}, bson.M{"$unset": bson.M{"keys." + key: 1}})
	if e != nil {
		return
	}
	_, e = d.mongo.Database("config_"+groupname).Collection(appname).DeleteMany(sctx, bson.M{"key": key})
	return
}

func (d *Dao) MongoCreate(ctx context.Context, groupname, appname, cipher string, encrypt datahandler) (e error) {
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
			if mongo.IsDuplicateKeyError(e) {
				e = ecode.ErrAppAlreadyExist
			}
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	col := d.mongo.Database("config_" + groupname).Collection(appname)
	index := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "key", Value: 1}, primitive.E{Key: "index", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, e = col.Indexes().CreateOne(sctx, index); e != nil {
		return
	}
	if _, e = col.InsertOne(sctx, bson.M{
		"key":    "",
		"index":  0,
		"cipher": cipher,
		"keys":   bson.M{},
	}); e != nil {
		return
	}
	return
}

type datahandler func(cipher string, origindata string) (newdata string)

func (d *Dao) MongoUpdateCipher(ctx context.Context, groupname, appname, oldcipher, newcipher string, decrypt, encrypt datahandler) (e error) {
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
	if appsummary.Cipher != oldcipher {
		e = ecode.ErrWrongCipher
		return
	}
	updater := bson.M{
		"cipher": newcipher,
	}
	for key, keysummary := range appsummary.Keys {
		if oldcipher != "" {
			keysummary.CurValue = decrypt(oldcipher, keysummary.CurValue)
		}
		if newcipher != "" {
			keysummary.CurValue = encrypt(newcipher, keysummary.CurValue)
		}
		updater["keys."+key+".cur_value"] = keysummary.CurValue
	}
	if _, e = col.UpdateOne(sctx, bson.M{"key": "", "index": 0}, bson.M{"$set": updater}); e != nil {
		return
	}
	var cursor *mongo.Cursor
	if cursor, e = col.Find(sctx, bson.M{"key": bson.M{"$ne": ""}, "index": bson.M{"$gt": 0}}, options.Find().SetSort(bson.M{"key": -1})); e != nil {
		return
	}
	for cursor.Next(sctx) {
		log := &model.Log{}
		if e = cursor.Decode(log); e != nil {
			return
		}
		if oldcipher != "" {
			log.Value = decrypt(oldcipher, log.Value)
		}
		if newcipher != "" {
			log.Value = encrypt(newcipher, log.Value)
		}
		if _, e = col.UpdateOne(sctx, bson.M{"key": log.Key, "index": log.Index}, bson.M{"$set": bson.M{"value": log.Value}}); e != nil {
			return
		}
	}
	e = cursor.Err()
	return
}

//index == 0 get the current index's config
//index != 0 get the specific index's config
func (d *Dao) MongoGetKeyConfig(ctx context.Context, groupname, appname, key string, index uint32, decrypt datahandler) (*model.KeySummary, *model.Log, error) {
	col := d.mongo.Database("config_"+groupname, options.Database().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())).Collection(appname)
	var appsummary *model.AppSummary
	var log *model.Log
	if index != 0 {
		//get the specific index's config and the current status
		filter := bson.M{"$or": bson.A{bson.M{"key": "", "index": 0}, bson.M{"key": key, "index": index}}}
		cursor, e := col.Find(ctx, filter, options.Find().SetSort(bson.M{"index": 1}))
		if e != nil {
			return nil, nil, e
		}
		for cursor.Next(ctx) {
			if appsummary == nil {
				tmps := &model.AppSummary{}
				if e = cursor.Decode(tmps); e != nil {
					return nil, nil, e
				}
				appsummary = tmps
			} else {
				tmpc := &model.Log{}
				if e = cursor.Decode(tmpc); e != nil {
					return nil, nil, e
				}
				log = tmpc
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
		if appsummary.Cipher != "" {
			keysummary.CurValue = decrypt(appsummary.Cipher, keysummary.CurValue)
			log.Value = decrypt(appsummary.Cipher, log.Value)
		}
		return keysummary, log, nil
	}
	//get tge current index's config
	appsummary = &model.AppSummary{}
	if e := col.FindOne(ctx, bson.M{"key": "", "index": 0}).Decode(appsummary); e != nil {
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
	if appsummary.Cipher != "" {
		keysummary.CurValue = decrypt(appsummary.Cipher, keysummary.CurValue)
	}
	log = &model.Log{
		Key:   key,
		Index: keysummary.CurIndex,
		Value: keysummary.CurValue,
	}
	return keysummary, log, nil
}

//get the app's all keys' current config
func (d *Dao) MongoGetAppConfig(ctx context.Context, groupname, appname string, decrypt datahandler) (*model.AppSummary, error) {
	app := &model.AppSummary{}
	if e := d.mongo.Database("config_"+groupname).Collection(appname).FindOne(ctx, bson.M{"key": "", "index": 0}).Decode(app); e != nil {
		if e == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, e
	}
	if app.Cipher != "" {
		for _, keysummary := range app.Keys {
			keysummary.CurValue = decrypt(app.Cipher, keysummary.CurValue)
		}
	}
	return app, nil
}
func (d *Dao) MongoSetConfig(ctx context.Context, groupname, appname, key, value string, encrypt datahandler) (newindex, newversion uint32, e error) {
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
	if appsummary.Cipher != "" {
		value = encrypt(appsummary.Cipher, value)
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
	if _, e = col.UpdateOne(sctx, bson.M{"key": "", "index": 0}, bson.M{"$set": bson.M{"keys." + key: keysummary}}); e != nil {
		return
	}
	if _, e = col.UpdateOne(sctx, bson.M{"key": key, "index": keysummary.CurIndex}, bson.M{"$set": bson.M{"value": value}}, options.Update().SetUpsert(true)); e != nil {
		return
	}
	newindex = keysummary.CurIndex
	newversion = keysummary.CurVersion
	return
}
func (d *Dao) MongoRollbackConfig(ctx context.Context, groupname, appname, key string, index uint32) (e error) {
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
			"keys." + key + ".cur_index": index,
			"keys." + key + ".cur_value": log.Value,
		},
		"$inc": bson.M{
			"keys." + key + ".cur_version": 1,
		},
	}
	if r := col.FindOneAndUpdate(sctx, bson.M{"key": "", "index": 0}, updateSummary); r.Err() != nil {
		if r.Err() == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		} else {
			e = r.Err()
		}
	}
	return
}

//first key groupname,second key appname,value curconfig
type WatchUpdateHandler func(string, string, *model.AppSummary)
type WatchDeleteAppHandler func(groupname, appname string)
type WatchDeleteConfigHandler func(groupname, appname string, id string)

func (d *Dao) getall(decrypt datahandler) (map[string]map[string]*model.AppSummary, error) {
	groups, e := d.MongoGetAllGroups(context.Background(), "")
	if e != nil {
		return nil, e
	}
	result := make(map[string]map[string]*model.AppSummary, len(groups))
	for _, group := range groups {
		tmpgroup := make(map[string]*model.AppSummary)
		apps, e := d.MongoGetAllApps(context.Background(), group, "")
		if e != nil {
			return nil, e
		}
		for _, app := range apps {
			tmpapp := &model.AppSummary{}
			if e := d.mongo.Database("config_"+group).Collection(app).FindOne(context.Background(), bson.M{"key": "", "index": 0}).Decode(tmpapp); e != nil {
				if e == mongo.ErrNoDocuments {
					log.Error(nil, "[MongoWatchConfig.getall] group:", group, "app:", app, "doesn't exist app summary")
					continue
				}
				log.Error(nil, "[MongoWatchConfig.getall] group:", group, "app:", app, "get app summary error:", e)
				continue
			}
			if tmpapp.Cipher != "" {
				for _, keysummary := range tmpapp.Keys {
					keysummary.CurValue = decrypt(tmpapp.Cipher, keysummary.CurValue)
				}
			}
			tmpgroup[app] = tmpapp
		}
		if len(tmpgroup) != 0 {
			result[group] = tmpgroup
		}
	}
	return result, nil
}
func (d *Dao) MongoWatchConfig(update WatchUpdateHandler, delA WatchDeleteAppHandler, delC WatchDeleteConfigHandler, decrypt datahandler) error {
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
					if s.Cipher != "" {
						for _, keysummary := range s.Keys {
							keysummary.CurValue = decrypt(s.Cipher, keysummary.CurValue)
						}
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
