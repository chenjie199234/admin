package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
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

func (d *Dao) MongoGetApp(ctx context.Context, gname, aname, secret string) (*model.AppSummary, error) {
	appsummary := &model.AppSummary{}
	e := d.mongo.Database("app").Collection("config").FindOne(ctx, bson.M{"group": gname, "app": aname, "key": "", "index": 0}).Decode(appsummary)
	if e != nil {
		return nil, e
	}
	// check sign
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
	if e := decodeProxyPath(appsummary); e != nil {
		return nil, e
	}
	return appsummary, nil
}
func (d *Dao) MongoGetPermissionNodeID(ctx context.Context, gname, aname string) (string, error) {
	appsummary := &model.AppSummary{}
	filterSummary := bson.M{"group": gname, "app": aname, "key": "", "index": 0}
	opts := options.FindOne().SetProjection(bson.M{"permission_node_id": 1})
	if e := d.mongo.Database("app").Collection("config").FindOne(ctx, filterSummary, opts).Decode(appsummary); e != nil {
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
func (d *Dao) MongoCreateApp(ctx context.Context, projectid, gname, aname, secret string) (e error) {
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
	if e = d.mongo.Database("permission").Collection("node").FindOneAndUpdate(sctx, bson.M{"node_id": projectid + model.AppControl}, bson.M{"$inc": bson.M{"cur_node_index": 1}}).Decode(parent); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrProjectNotExist
		}
		return
	}
	if _, e = d.mongo.Database("permission").Collection("node").InsertOne(sctx, &model.Node{
		NodeId:       parent.NodeId + "," + strconv.FormatUint(uint64(parent.CurNodeIndex+1), 10),
		NodeName:     gname + "." + aname,
		NodeData:     "",
		CurNodeIndex: 0,
	}); e != nil {
		return
	}
	nonce := make([]byte, 32)
	rand.Read(nonce)
	if _, e = d.mongo.Database("app").Collection("config").InsertOne(sctx, &model.AppSummary{
		Group:            gname,
		App:              aname,
		Key:              "",
		Index:            0,
		Paths:            map[string]*model.ProxyPath{},
		Keys:             map[string]*model.KeySummary{},
		Value:            util.SignMake(secret, nonce),
		PermissionNodeID: parent.NodeId + "," + strconv.FormatUint(uint64(parent.CurNodeIndex+1), 10),
	}); e != nil && mongo.IsDuplicateKeyError(e) {
		fmt.Println(e)
		e = ecode.ErrAppAlreadyExist
	}
	return
}
func (d *Dao) MongoDelApp(ctx context.Context, gname, aname, secret string) (e error) {
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
	appsummary := &model.AppSummary{}
	filterSummary := bson.M{"group": gname, "app": aname, "key": "", "index": 0}
	opts := options.FindOne().SetProjection(bson.M{"value": 1, "permission_node_id": 1})
	if e = d.mongo.Database("app").Collection("config").FindOne(sctx, filterSummary, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = util.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	if _, e = d.mongo.Database("app").Collection("config").DeleteMany(sctx, bson.M{"group": gname, "app": aname}); e != nil {
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
func (d *Dao) MongoUpdateAppSecret(ctx context.Context, gname, aname, oldsecret, newsecret string) (e error) {
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
	appsummary := &model.AppSummary{}
	filterSummary := bson.M{"group": gname, "app": aname, "key": "", "index": 0}
	opts := options.FindOne().SetProjection(bson.M{"value": 1, "keys": 1})
	if e = d.mongo.Database("app").Collection("config").FindOne(sctx, filterSummary, opts).Decode(appsummary); e != nil {
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
	updaterSummary := bson.M{
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
			updaterSummary["keys."+key+".cur_value"], _ = util.Encrypt(newsecret, common.Str2byte(keysummary.CurValue))
		} else {
			updaterSummary["keys."+key+".cur_value"] = keysummary.CurValue
		}
	}
	if _, e = d.mongo.Database("app").Collection("config").UpdateOne(sctx, filterSummary, bson.M{"$set": updaterSummary}); e != nil {
		return
	}
	filterlog := bson.M{"group": gname, "app": aname, "key": bson.M{"$exists": true, "$type": "string", "$ne": ""}, "index": bson.M{"$gt": 0}}
	var cursor *mongo.Cursor
	if cursor, e = d.mongo.Database("app").Collection("config").Find(sctx, filterlog); e != nil {
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
		filter := bson.M{"group": gname, "app": aname, "key": log.Key, "index": log.Index}
		updater := bson.M{"$set": bson.M{"value": log.Value}}
		if _, e = d.mongo.Database("app").Collection("config").UpdateOne(sctx, filter, updater); e != nil {
			return
		}
	}
	e = cursor.Err()
	return
}

// index == 0 get the current index's config
// index != 0 get the specific index's config
func (d *Dao) MongoGetKeyConfig(ctx context.Context, gname, aname, key string, index uint32, secret string) (*model.KeySummary, *model.Log, error) {
	col := d.mongo.Database("app", options.Database().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())).Collection("config")
	var appsummary *model.AppSummary
	var log *model.Log
	if index == 0 {
		//get tge current index's config
		appsummary = &model.AppSummary{}
		filterSummary := bson.M{"group": gname, "app": aname, "key": "", "index": 0}
		opts := options.FindOne().SetProjection(bson.M{"value": 1, "keys." + key: 1})
		if e := col.FindOne(ctx, filterSummary, opts).Decode(appsummary); e != nil {
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
	filter := bson.M{"group": gname, "app": aname, "$or": bson.A{bson.M{"key": "", "index": 0}, bson.M{"key": key, "index": index}}}
	opts := options.Find().SetProjection(bson.M{"key": 1, "index": 1, "value": 1, "value_type": 1, "keys." + key: 1}).SetSort(bson.M{"index": 1})
	cursor, e := col.Find(ctx, filter, opts)
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
func (d *Dao) MongoSetKeyConfig(ctx context.Context, gname, aname, key, secret, value, valuetype string) (newindex, newversion uint32, e error) {
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
	appsummary := &model.AppSummary{}
	filterSummary := bson.M{"group": gname, "app": aname, "key": "", "index": 0}
	opts := options.FindOne().SetProjection(bson.M{"value": 1, "keys." + key: 1})
	if e = d.mongo.Database("app").Collection("config").FindOne(sctx, filterSummary, opts).Decode(appsummary); e != nil {
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
	updaterSummary := bson.M{"$set": bson.M{"keys." + key: keysummary}}
	if _, e = d.mongo.Database("app").Collection("config").UpdateOne(sctx, filterSummary, updaterSummary); e != nil {
		return
	}
	filterLog := bson.M{"group": gname, "app": aname, "key": key, "index": keysummary.CurIndex}
	updaterLog := bson.M{"$set": bson.M{"value": value, "value_type": valuetype}}
	if _, e = d.mongo.Database("app").Collection("config").UpdateOne(sctx, filterLog, updaterLog, options.Update().SetUpsert(true)); e != nil {
		return
	}
	newindex = keysummary.CurIndex
	newversion = keysummary.CurVersion
	return
}
func (d *Dao) MongoDelKey(ctx context.Context, gname, aname, key, secret string) (e error) {
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
	appsummary := &model.AppSummary{}
	filterSummary := bson.M{"group": gname, "app": aname, "key": "", "index": 0}
	updaterSummary := bson.M{"$unset": bson.M{"keys." + key: 1}}
	opts := options.FindOneAndUpdate().SetProjection(bson.M{"value": 1})
	if e = d.mongo.Database("app").Collection("config").FindOneAndUpdate(sctx, filterSummary, updaterSummary, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = util.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	delfilter := bson.M{"group": gname, "app": aname, "key": key}
	_, e = d.mongo.Database("app").Collection("config").DeleteMany(sctx, delfilter)
	return
}
func (d *Dao) MongoRollbackKeyConfig(ctx context.Context, gname, aname, key, secret string, index uint32) (e error) {
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
	appsummary := &model.AppSummary{}
	filterSummary := bson.M{"group": gname, "app": aname, "key": "", "index": 0}
	if e = d.mongo.Database("app").Collection("config").FindOne(sctx, filterSummary).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = util.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	if len(appsummary.Keys) == 0 {
		e = ecode.ErrKeyNotExist
		return
	}
	if keysummary, ok := appsummary.Keys[key]; !ok {
		e = ecode.ErrKeyNotExist
		return
	} else if keysummary.CurIndex == index {
		return
	}
	log := &model.Log{}
	filterLog := bson.M{"group": gname, "app": aname, "key": key, "index": index}
	if e = d.mongo.Database("app").Collection("config").FindOne(sctx, filterLog).Decode(log); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrIndexNotExist
		}
		return
	}
	updaterSummary := bson.M{
		"$set": bson.M{
			"keys." + key + ".cur_index":      index,
			"keys." + key + ".cur_value":      log.Value,
			"keys." + key + ".cur_value_type": log.ValueType,
		},
		"$inc": bson.M{
			"keys." + key + ".cur_version": 1,
		},
	}
	_, e = d.mongo.Database("app").Collection("config").UpdateOne(sctx, filterSummary, updaterSummary)
	return
}
func (d *Dao) MongoSetProxyPath(ctx context.Context, gname, aname, secret, path string, read, write, admin bool) (e error) {
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
	b64path := encodeProxyPath(path)
	appsummary := &model.AppSummary{}
	filter := bson.M{"group": gname, "app": aname, "key": "", "index": 0}
	updater1 := bson.M{"$set": bson.M{"paths." + b64path + ".permission_read": read, "paths." + b64path + ".permission_write": write, "paths." + b64path + ".permission_admin": admin}}
	opts := options.FindOneAndUpdate().SetProjection(bson.M{"value": 1, "paths." + b64path: 1, "permission_node_id": 1})
	if e = d.mongo.Database("app").Collection("config").FindOneAndUpdate(sctx, filter, updater1, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = util.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	addpermission := false
	if len(appsummary.Paths) == 0 {
		addpermission = true
	} else if _, ok := appsummary.Paths[b64path]; !ok {
		addpermission = true
	}
	if !addpermission {
		return
	}
	parent := &model.Node{}
	if e = d.mongo.Database("permission").Collection("node").FindOneAndUpdate(sctx, bson.M{"node_id": appsummary.PermissionNodeID}, bson.M{"$inc": bson.M{"cur_node_index": 1}}).Decode(parent); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppPermissionMissing
		}
		return
	}
	newnodeid := parent.NodeId + "," + strconv.FormatUint(uint64(parent.CurNodeIndex+1), 10)
	if _, e = d.mongo.Database("permission").Collection("node").InsertOne(sctx, &model.Node{
		NodeId:       newnodeid,
		NodeName:     path,
		NodeData:     "",
		CurNodeIndex: 0,
	}); e != nil {
		return
	}
	updater2 := bson.M{"$set": bson.M{"paths." + b64path + ".permission_node_id": newnodeid}}
	_, e = d.mongo.Database("app").Collection("config").UpdateOne(sctx, filter, updater2)
	return e
}
func (d *Dao) MongoDelProxyPath(ctx context.Context, gname, aname, secret, path string) (e error) {
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
	b64path := base64.StdEncoding.EncodeToString(common.Str2byte(path))
	appsummary := &model.AppSummary{}
	filter := bson.M{"group": gname, "app": aname, "key": "", "index": 0}
	updater := bson.M{"$unset": bson.M{"paths." + b64path: 1}}
	opts := options.FindOneAndUpdate().SetProjection(bson.M{"value": 1, "paths." + b64path: 1})
	if e = d.mongo.Database("app").Collection("config").FindOneAndUpdate(ctx, filter, updater, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = util.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	permissionid := ""
	if len(appsummary.Paths) != 0 {
		if proxypath, ok := appsummary.Paths[b64path]; ok {
			permissionid = proxypath.PermissionNodeID
		}
	}
	if permissionid == "" {
		return
	}
	delfilter := bson.M{"node_id": bson.M{"$regex": "^" + permissionid}}
	if _, e = d.mongo.Database("permission").Collection("node").DeleteMany(sctx, delfilter); e != nil {
		return
	}
	if _, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, delfilter); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, delfilter)
	return
}

// first key groupname,second key appname,value curconfig
type WatchDropCollectionHandler func()
type WatchUpdateHandler func(string, string, *model.AppSummary)
type WatchDeleteConfigHandler func(id string)

func (d *Dao) MongoWatchConfig(drop WatchDropCollectionHandler, update WatchUpdateHandler, delC WatchDeleteConfigHandler, initall map[string]*model.AppSummary) error {
	starttime := &primitive.Timestamp{T: uint32(time.Now().Unix()) - 1, I: uint32(0)}
	watchfilter := mongo.Pipeline{bson.D{primitive.E{Key: "$match", Value: bson.M{"ns.db": "app", "ns.coll": "config"}}}}

	if e := d.mongoGetAll(initall); e != nil {
		return e
	}
	go func() {
		var stream *mongo.ChangeStream
		for {
			for stream == nil {
				//connect
				var e error
				if stream, e = d.mongo.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup).SetStartAtOperationTime(starttime)); e != nil {
					log.Error(nil, "[dao.MongoWatchConfig] connect stream error:", e)
					stream = nil
					time.Sleep(time.Millisecond * 100)
					continue
				}
			}
			for stream.Next(context.Background()) {
				starttime.T, starttime.I = stream.Current.Lookup("clusterTime").Timestamp()
				starttime.I++
				switch stream.Current.Lookup("operationType").StringValue() {
				case "drop":
					//drop collection
					drop()
				case "insert":
					//insert document
					fallthrough
				case "update":
					//update document
					gname, gok := stream.Current.Lookup("fullDocument").Document().Lookup("group").StringValueOK()
					aname, aok := stream.Current.Lookup("fullDocument").Document().Lookup("app").StringValueOK()
					key, kok := stream.Current.Lookup("fullDocument").Document().Lookup("key").StringValueOK()
					index, iok := stream.Current.Lookup("fullDocument").Document().Lookup("index").AsInt32OK()
					if !gok || !aok || !kok || !iok {
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
						log.Error(nil, "[dao.MongoWatchConfig] group:", gname, "app:", aname, "summary data broken:", e)
						continue
					}
					//decode proxy path
					if e := decodeProxyPath(s); e != nil {
						log.Error(nil, "[dao.MongoWatchConfig] group:", gname, "app:", aname, "proxy path broken:", e)
						continue
					}
					update(gname, aname, s)
				case "delete":
					//delete document
					delC(stream.Current.Lookup("documentKey").Document().Lookup("_id").ObjectID().Hex())
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

// this function can't guarantee atomic
// this function can only be used in MongoWatchConfig function
func (d *Dao) mongoGetAll(initall map[string]*model.AppSummary) error {
	if initall == nil {
		return nil
	}
	filter := bson.M{
		"permission_node_id": bson.M{"$exists": true, "$type": "string", "$ne": ""},
		"key":                "",
		"index":              0,
	}
	var cursor *mongo.Cursor
	cursor, e := d.mongo.Database("app").Collection("config").Find(context.Background(), filter)
	if e != nil {
		return e
	}
	defer cursor.Close(context.Background())
	tmp := make([]*model.AppSummary, 0, cursor.RemainingBatchLength())
	if e := cursor.All(context.Background(), &tmp); e != nil {
		return e
	}
	for _, v := range tmp {
		initall[v.Group+"."+v.App] = v
	}
	return nil
}

// this function will not decrypt
func (d *Dao) MongoGetAppConfig(ctx context.Context, gname, aname string) (*model.AppSummary, error) {
	app := &model.AppSummary{}
	filter := bson.M{"group": gname, "app": aname, "key": "", "index": 0}
	if e := d.mongo.Database("app").Collection("config").FindOne(ctx, filter).Decode(app); e != nil {
		if e == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, e
	}
	if e := decodeProxyPath(app); e != nil {
		return nil, e
	}
	return app, nil
}
func encodeProxyPath(path string) string {
	return base64.StdEncoding.EncodeToString(common.Str2byte(path))
}
func decodeProxyPath(app *model.AppSummary) error {
	tmp := make(map[string]*model.ProxyPath)
	for path, info := range app.Paths {
		realpath, e := base64.StdEncoding.DecodeString(path)
		if e != nil {
			return e
		}
		tmp[common.Byte2str(realpath)] = info
	}
	app.Paths = tmp
	return nil
}
