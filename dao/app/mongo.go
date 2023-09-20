package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strconv"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/secure"
	"github.com/chenjie199234/Corelib/util/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (d *Dao) MongoGetApp(ctx context.Context, projectid, gname, aname, secret string) (*model.AppSummary, error) {
	appsummary := &model.AppSummary{}
	e := d.mongo.Database("app").Collection("config").FindOne(ctx, bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}).Decode(appsummary)
	if e != nil {
		return nil, e
	}
	// check sign
	if e := secure.SignCheck(secret, appsummary.Value); e != nil {
		return nil, e
	}
	if secret != "" {
		for _, keysummary := range appsummary.Keys {
			plaintext, e := secure.AesDecrypt(secret, keysummary.CurValue)
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
func (d *Dao) MongoGetAppWithoutDecrypt(ctx context.Context, projectid, gname, aname string) (*model.AppSummary, error) {
	app := &model.AppSummary{}
	filter := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
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
func (d *Dao) MongoGetPermissionNodeID(ctx context.Context, projectid, gname, aname string) (string, error) {
	appsummary := &model.AppSummary{}
	filterSummary := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
	opts := options.FindOne().SetProjection(bson.M{"permission_node_id": 1})
	if e := d.mongo.Database("app").Collection("config").FindOne(ctx, filterSummary, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return "", e
	}
	return appsummary.PermissionNodeID, nil
}
func (d *Dao) MongoCreateApp(ctx context.Context, projectid, gname, aname, secret, discovermode, kubernetesns, kubernetesls, dnshost string, dnsinterval uint32, staticaddrs []string) (nodeid string, e error) {
	var sign string
	sign, e = secure.SignMake(secret)
	if e != nil {
		return "", e
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
	projectindex := &model.ProjectIndex{}
	if e = d.mongo.Database("permission").Collection("projectindex").FindOne(sctx, bson.M{"project_id": projectid}).Decode(projectindex); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrProjectNotExist
		}
		return
	}
	parent := &model.Node{}
	if e = d.mongo.Database("permission").Collection("node").FindOneAndUpdate(sctx, bson.M{"node_id": projectid + model.AppControl}, bson.M{"$inc": bson.M{"cur_node_index": 1}}).Decode(parent); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrProjectNotExist
		}
		return
	}
	nodeid = parent.NodeId + "," + strconv.FormatUint(uint64(parent.CurNodeIndex+1), 10)
	if _, e = d.mongo.Database("permission").Collection("node").InsertOne(sctx, &model.Node{
		NodeId:       nodeid,
		NodeName:     gname + "." + aname,
		NodeData:     "",
		CurNodeIndex: 0,
	}); e != nil {
		return
	}
	if _, e = d.mongo.Database("app").Collection("config").InsertOne(sctx, &model.AppSummary{
		ProjectID:        projectid,
		ProjectName:      projectindex.ProjectName,
		Group:            gname,
		App:              aname,
		Key:              "",
		Index:            0,
		DiscoverMode:     discovermode,
		KubernetesNs:     kubernetesns,
		KubernetesLS:     kubernetesls,
		DnsHost:          dnshost,
		DnsInterval:      dnsinterval,
		StaticAddrs:      staticaddrs,
		Paths:            map[string]*model.ProxyPath{},
		Keys:             map[string]*model.KeySummary{},
		Value:            sign,
		PermissionNodeID: nodeid,
	}); e != nil && mongo.IsDuplicateKeyError(e) {
		e = ecode.ErrAppAlreadyExist
	}
	return
}
func (d *Dao) MongoUpdateApp(ctx context.Context, projectid, gname, aname, secret, discovermode, kubernetesns, kubernetesls, dnshost string, dnsinterval uint32, staticaddrs []string) (nodeid string, e error) {
	filter := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
	updater := bson.M{
		"$set": bson.M{
			"discover_mode": discovermode,
			"kubernetes_ns": kubernetesns,
			"kubernetes_ls": kubernetesls,
			"dns_host":      dnshost,
			"dns_interval":  dnsinterval,
			"static_addrs":  staticaddrs,
		},
	}
	app := &model.AppSummary{}
	e = d.mongo.Database("app").Collection("config").FindOneAndUpdate(ctx, filter, updater, options.FindOneAndUpdate().SetProjection(bson.M{"permission_node_id": 1})).Decode(app)
	if e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	nodeid = app.PermissionNodeID
	return
}
func (d *Dao) MongoDelApp(ctx context.Context, projectid, gname, aname, secret string) (e error) {
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
	filterSummary := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
	opts := options.FindOne().SetProjection(bson.M{"value": 1, "permission_node_id": 1})
	if e = d.mongo.Database("app").Collection("config").FindOne(sctx, filterSummary, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = secure.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	if _, e = d.mongo.Database("app").Collection("config").DeleteMany(sctx, bson.M{"project_id": projectid, "group": gname, "app": aname}); e != nil {
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
func (d *Dao) MongoUpdateAppSecret(ctx context.Context, projectid, gname, aname, oldsecret, newsecret string) (e error) {
	if oldsecret == newsecret {
		return
	}
	var sign string
	sign, e = secure.SignMake(newsecret)
	if e != nil {
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
	filterSummary := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
	opts := options.FindOne().SetProjection(bson.M{"value": 1, "keys": 1})
	if e = d.mongo.Database("app").Collection("config").FindOne(sctx, filterSummary, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	//check oldsecret
	if e = secure.SignCheck(oldsecret, appsummary.Value); e != nil {
		return
	}
	nonce := make([]byte, 32)
	rand.Read(nonce)
	updaterSummary := bson.M{
		"value": sign,
	}
	for key, keysummary := range appsummary.Keys {
		if oldsecret != "" {
			var plaintext []byte
			plaintext, e = secure.AesDecrypt(oldsecret, keysummary.CurValue)
			if e != nil {
				return
			}
			keysummary.CurValue = common.Byte2str(plaintext)
		}
		if newsecret != "" {
			updaterSummary["keys."+key+".cur_value"], _ = secure.AesEncrypt(newsecret, common.Str2byte(keysummary.CurValue))
		} else {
			updaterSummary["keys."+key+".cur_value"] = keysummary.CurValue
		}
	}
	if _, e = d.mongo.Database("app").Collection("config").UpdateOne(sctx, filterSummary, bson.M{"$set": updaterSummary}); e != nil {
		return
	}
	filterlog := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": bson.M{"$exists": true, "$type": "string", "$ne": ""}, "index": bson.M{"$gt": 0}}
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
			plaintext, e = secure.AesDecrypt(oldsecret, log.Value)
			if e != nil {
				return
			}
			log.Value = common.Byte2str(plaintext)
		}
		if newsecret != "" {
			log.Value, _ = secure.AesEncrypt(newsecret, common.Str2byte(log.Value))
		}
		filter := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": log.Key, "index": log.Index}
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
func (d *Dao) MongoGetKeyConfig(ctx context.Context, projectid, gname, aname, key string, index uint32, secret string) (*model.KeySummary, *model.Log, error) {
	col := d.mongo.Database("app", options.Database().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())).Collection("config")
	var appsummary *model.AppSummary
	var log *model.Log
	if index == 0 {
		//get tge current index's config
		appsummary = &model.AppSummary{}
		filterSummary := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
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
		if e := secure.SignCheck(secret, appsummary.Value); e != nil {
			return nil, nil, e
		}
		if secret != "" {
			plaintext, e := secure.AesDecrypt(secret, keysummary.CurValue)
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
	filter := bson.M{"project_id": projectid, "group": gname, "app": aname, "$or": bson.A{bson.M{"key": "", "index": 0}, bson.M{"key": key, "index": index}}}
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
	if e := secure.SignCheck(secret, appsummary.Value); e != nil {
		return nil, nil, e
	}
	if secret != "" {
		plaintext, e := secure.AesDecrypt(secret, keysummary.CurValue)
		if e != nil {
			return nil, nil, e
		}
		keysummary.CurValue = common.Byte2str(plaintext)
		plaintext, e = secure.AesDecrypt(secret, log.Value)
		if e != nil {
			return nil, nil, e
		}
		log.Value = common.Byte2str(plaintext)
	}
	return keysummary, log, nil
}
func (d *Dao) MongoSetKeyConfig(ctx context.Context, projectid, gname, aname, key, secret, value, valuetype string, newkey bool) (newindex, newversion uint32, e error) {
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
	filterSummary := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
	opts := options.FindOne().SetProjection(bson.M{"value": 1, "keys." + key: 1})
	if e = d.mongo.Database("app").Collection("config").FindOne(sctx, filterSummary, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	//check secret
	if e = secure.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	if secret != "" {
		if value, e = secure.AesEncrypt(secret, common.Str2byte(value)); e != nil {
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
	} else if newkey {
		e = ecode.ErrKeyAlreadyExist
		return
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
	filterLog := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": key, "index": keysummary.CurIndex}
	updaterLog := bson.M{"$set": bson.M{"value": value, "value_type": valuetype}}
	if _, e = d.mongo.Database("app").Collection("config").UpdateOne(sctx, filterLog, updaterLog, options.Update().SetUpsert(true)); e != nil {
		return
	}
	newindex = keysummary.CurIndex
	newversion = keysummary.CurVersion
	return
}
func (d *Dao) MongoDelKey(ctx context.Context, projectid, gname, aname, key, secret string) (e error) {
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
	filterSummary := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
	updaterSummary := bson.M{"$unset": bson.M{"keys." + key: 1}}
	opts := options.FindOneAndUpdate().SetProjection(bson.M{"value": 1})
	if e = d.mongo.Database("app").Collection("config").FindOneAndUpdate(sctx, filterSummary, updaterSummary, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = secure.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	delfilter := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": key}
	_, e = d.mongo.Database("app").Collection("config").DeleteMany(sctx, delfilter)
	return
}
func (d *Dao) MongoRollbackKeyConfig(ctx context.Context, projectid, gname, aname, key, secret string, index uint32) (e error) {
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
	filterSummary := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
	if e = d.mongo.Database("app").Collection("config").FindOne(sctx, filterSummary).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = secure.SignCheck(secret, appsummary.Value); e != nil {
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
	filterLog := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": key, "index": index}
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
func (d *Dao) MongoSetProxyPath(ctx context.Context, projectid, gname, aname, secret, path string, read, write, admin, newpath bool) (nodeid string, e error) {
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
	filter := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
	updater1 := bson.M{
		"$set": bson.M{
			"paths." + b64path + ".permission_read":  read,
			"paths." + b64path + ".permission_write": write,
			"paths." + b64path + ".permission_admin": admin,
		},
	}
	opts := options.FindOneAndUpdate().SetProjection(bson.M{"value": 1, "paths." + b64path: 1, "permission_node_id": 1})
	if e = d.mongo.Database("app").Collection("config").FindOneAndUpdate(sctx, filter, updater1, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = secure.SignCheck(secret, appsummary.Value); e != nil {
		return
	}
	addpermission := false
	if len(appsummary.Paths) == 0 {
		addpermission = true
	} else if path, ok := appsummary.Paths[b64path]; !ok {
		addpermission = true
	} else {
		nodeid = path.PermissionNodeID
	}
	if !addpermission {
		if newpath {
			e = ecode.ErrProxyPathAlreadyExist
		}
		return
	}
	parent := &model.Node{}
	if e = d.mongo.Database("permission").Collection("node").FindOneAndUpdate(sctx, bson.M{"node_id": appsummary.PermissionNodeID}, bson.M{"$inc": bson.M{"cur_node_index": 1}}).Decode(parent); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppPermissionMissing
		}
		return
	}
	nodeid = parent.NodeId + "," + strconv.FormatUint(uint64(parent.CurNodeIndex+1), 10)
	if _, e = d.mongo.Database("permission").Collection("node").InsertOne(sctx, &model.Node{
		NodeId:       nodeid,
		NodeName:     path,
		NodeData:     "",
		CurNodeIndex: 0,
	}); e != nil {
		return
	}
	updater2 := bson.M{"$set": bson.M{"paths." + b64path + ".permission_node_id": nodeid}}
	_, e = d.mongo.Database("app").Collection("config").UpdateOne(sctx, filter, updater2)
	return
}
func (d *Dao) MongoDelProxyPath(ctx context.Context, projectid, gname, aname, secret, path string) (e error) {
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
	filter := bson.M{"project_id": projectid, "group": gname, "app": aname, "key": "", "index": 0}
	updater := bson.M{"$unset": bson.M{"paths." + b64path: 1}}
	opts := options.FindOneAndUpdate().SetProjection(bson.M{"value": 1, "paths." + b64path: 1})
	if e = d.mongo.Database("app").Collection("config").FindOneAndUpdate(sctx, filter, updater, opts).Decode(appsummary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if e = secure.SignCheck(secret, appsummary.Value); e != nil {
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
