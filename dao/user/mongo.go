package user

import (
	"context"
	"strings"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func (d *Dao) MongoUserLogin(ctx context.Context, mobile, oauth2username, oauth2type string) (bson.ObjectID, error) {
	filter := bson.M{"mobile": mobile}
	updater := bson.M{}
	if oauth2type == "FeiShu" {
		updater["feishu_user_name"] = oauth2username
	} else if oauth2type == "DingDing" {
		updater["dingding_user_name"] = oauth2username
	} else if oauth2type == "WXWork" {
		updater["wxwork_user_name"] = oauth2username
	} else {
		return bson.NilObjectID, ecode.ErrReq
	}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetProjection(bson.M{"_id": 1})
	user := &model.User{}
	e := d.mongo.Database("user").Collection("user").FindOneAndUpdate(ctx, filter, bson.M{"$set": updater}, opts).Decode(user)
	if e != nil && e != mongo.ErrNoDocuments {
		return bson.NilObjectID, e
	}
	if e == nil {
		return user.ID, nil
	}
	e = d.mongo.Database("user").Collection("user").FindOne(ctx, filter, options.FindOne().SetProjection(bson.M{"_id": 1})).Decode(user)
	if e != nil {
		return bson.NilObjectID, e
	}
	return user.ID, nil
}

func (d *Dao) MongoInvite(ctx context.Context, operator bson.ObjectID, projectid string, target bson.ObjectID) error {
	if target == bson.NilObjectID {
		return ecode.ErrReq
	}
	if !strings.HasPrefix(projectid, "0,") || strings.Count(projectid, ",") != 1 {
		return ecode.ErrReq
	}
	filter := bson.M{"_id": target, "projects." + projectid: bson.M{"$exists": false}}
	updater := bson.M{"$set": bson.M{"projects." + projectid: []string{}}}
	_, e := d.mongo.Database("user").Collection("user").UpdateOne(ctx, filter, updater)
	return e
}
func (d *Dao) MongoKick(ctx context.Context, operator bson.ObjectID, projectid string, target bson.ObjectID) (e error) {
	if target == bson.NilObjectID {
		return ecode.ErrReq
	}
	if !strings.HasPrefix(projectid, "0,") || strings.Count(projectid, ",") != 1 {
		return ecode.ErrReq
	}
	var s *mongo.Session
	if s, e = d.mongo.StartSession(); e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(options.Transaction().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	var r *mongo.UpdateResult
	r, e = d.mongo.Database("user").Collection("user").UpdateOne(sctx, bson.M{"_id": target}, bson.M{"$unset": bson.M{"projects." + projectid: 1}})
	if e != nil {
		return
	}
	if r.MatchedCount == 0 {
		e = ecode.ErrUserNotExist
		return
	}
	if r.ModifiedCount == 0 {
		e = ecode.ErrUserNotInProject
		return
	}
	_, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, bson.M{"user_id": target, "node_id": bson.M{"$regex": "^" + projectid}})
	return
}

func (d *Dao) MongoGetUsers(ctx context.Context, userids []bson.ObjectID) (map[bson.ObjectID]*model.User, error) {
	cursor, e := d.mongo.Database("user").Collection("user").Find(ctx, bson.M{"_id": bson.M{"$in": userids}})
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	result := make(map[bson.ObjectID]*model.User, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		tmp := &model.User{}
		if e := cursor.Decode(tmp); e != nil {
			return nil, e
		}
		result[tmp.ID] = tmp
	}
	return result, cursor.Err()
}

// page: 0:means return all,>0:means return the required page,if page overflow,the last page will return
func (d *Dao) MongoSearchUsers(ctx context.Context, projectid, name string, pagesize, page int64) (map[bson.ObjectID]*model.User, int64, int64, error) {
	filter := bson.M{}
	if name != "" {
		filter["$or"] = bson.A{
			bson.M{"feishu_user_name": bson.M{"$regex": name}},
			bson.M{"dingding_user_name": bson.M{"$regex": name}},
			bson.M{"wxwork_user_name": bson.M{"$regex": name}},
		}
	}
	if projectid != "" {
		filter["projects."+projectid] = bson.M{"$exists": true}
	}
	totalsize, e := d.mongo.Database("user").Collection("user").CountDocuments(ctx, filter)
	if e != nil {
		return nil, 0, 0, e
	}
	if totalsize == 0 {
		return make(map[bson.ObjectID]*model.User), 0, 0, nil
	}
	opts := options.Find().SetSort(bson.M{"_id": -1})
	if page != 0 {
		skip := (page - 1) * pagesize
		if skip >= totalsize {
			if totalsize%pagesize > 0 {
				page = totalsize/pagesize + 1
			} else {
				page = totalsize / pagesize
			}
			skip = (page - 1) * pagesize
		}
		opts = opts.SetSkip(skip).SetLimit(pagesize)
	}
	cursor, e := d.mongo.Database("user").Collection("user").Find(ctx, filter, opts)
	if e != nil {
		return nil, 0, 0, e
	}
	defer cursor.Close(ctx)
	result := make(map[bson.ObjectID]*model.User, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		tmp := &model.User{}
		if e := cursor.Decode(tmp); e != nil {
			return nil, 0, 0, e
		}
		result[tmp.ID] = tmp
	}
	return result, page, totalsize, cursor.Err()
}

func (d *Dao) MongoUpdateUser(ctx context.Context, userid bson.ObjectID, newname string) (*model.User, error) {
	user := &model.User{}
	filter := bson.M{"_id": userid}
	updater := bson.M{"$set": bson.M{"oauth2_user_name": newname}}
	e := d.mongo.Database("user").Collection("user").FindOneAndUpdate(ctx, filter, updater).Decode(user)
	if e == mongo.ErrNoDocuments {
		e = ecode.ErrUserNotExist
	}
	return user, e
}

func (d *Dao) MongoDelUsers(ctx context.Context, userids []bson.ObjectID) (e error) {
	var s *mongo.Session
	if s, e = d.mongo.StartSession(); e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(options.Transaction().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	if _, e = d.mongo.Database("user").Collection("user").DeleteMany(sctx, bson.M{"_id": bson.M{"$in": userids}}); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, bson.M{"user_id": bson.M{"$in": userids}})
	return
}

func (d *Dao) MongoCreateRole(ctx context.Context, projectid, name, comment string) (e error) {
	if _, e = d.mongo.Database("user").Collection("role").InsertOne(ctx, &model.Role{
		ProjectID: projectid,
		RoleName:  name,
		Comment:   comment,
	}); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			e = ecode.ErrRoleAlreadyExist
		}
		return e
	}
	return
}

// if limit is 0 means all
func (d *Dao) MongoSearchRoles(ctx context.Context, projectid, name string, pagesize, page int64) (map[string]*model.Role, int64, int64, error) {
	filter := bson.M{"project_id": projectid}
	if name != "" {
		filter["role_name"] = bson.M{"$regex": name}
	}
	totalsize, e := d.mongo.Database("user").Collection("role").CountDocuments(ctx, filter)
	if e != nil {
		return nil, 0, 0, e
	}
	if totalsize == 0 {
		return make(map[string]*model.Role), 0, 0, nil
	}
	opts := options.Find().SetSort(bson.M{"_id": -1})
	if page != 0 {
		skip := (page - 1) * pagesize
		if skip >= totalsize {
			if totalsize%pagesize > 0 {
				page = totalsize/pagesize + 1
			} else {
				page = totalsize / pagesize
			}
			skip = (page - 1) * totalsize
		}
		opts = opts.SetSkip(skip).SetLimit(pagesize)
	}
	cursor, e := d.mongo.Database("user").Collection("role").Find(ctx, filter, opts)
	if e != nil {
		return nil, 0, 0, e
	}
	defer cursor.Close(ctx)
	result := make(map[string]*model.Role, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		tmp := &model.Role{}
		if e := cursor.Decode(tmp); e != nil {
			return nil, 0, 0, e
		}
		result[tmp.RoleName] = tmp
	}
	return result, page, totalsize, cursor.Err()
}

func (d *Dao) MongoUpdateRole(ctx context.Context, projectid, name, newcomment string) (*model.Role, error) {
	role := &model.Role{}
	filter := bson.M{"project_id": projectid, "role_name": name}
	updater := bson.M{"$set": bson.M{"comment": newcomment}}
	e := d.mongo.Database("user").Collection("role").FindOneAndUpdate(ctx, filter, updater).Decode(role)
	if e == mongo.ErrNoDocuments {
		e = ecode.ErrRoleNotExist
	}
	return role, e
}

func (d *Dao) MongoDelRoles(ctx context.Context, projectid string, rolenames []string) (e error) {
	var s *mongo.Session
	if s, e = d.mongo.StartSession(); e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(options.Transaction().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	if _, e = d.mongo.Database("user").Collection("role").DeleteMany(sctx, bson.M{"project_id": projectid, "role_name": bson.M{"$in": rolenames}}); e != nil {
		return
	}
	filter := bson.M{"projects." + projectid: bson.M{"$exists": true}}
	updater := bson.M{"$pullAll": bson.M{"projects." + projectid: rolenames}}
	if _, e = d.mongo.Database("user").Collection("user").UpdateMany(sctx, filter, updater); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, bson.M{"project_id": projectid, "role_name": bson.M{"$in": rolenames}})
	return
}

func (d *Dao) MongoAddUserRole(ctx context.Context, userid bson.ObjectID, projectid, rolename string) (e error) {
	var s *mongo.Session
	if s, e = d.mongo.StartSession(); e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(options.Transaction().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	var exist int64
	exist, e = d.mongo.Database("user").Collection("role").CountDocuments(sctx, bson.M{"project_id": projectid, "role_name": rolename})
	if e != nil {
		return
	}
	if exist == 0 {
		e = ecode.ErrRoleNotExist
		return
	}
	filter := bson.M{"_id": userid, "projects." + projectid: bson.M{"$exists": true}}
	updater := bson.M{"$addToSet": bson.M{"projects." + projectid: rolename}}
	_, e = d.mongo.Database("user").Collection("user").UpdateOne(sctx, filter, updater)
	return
}
func (d *Dao) MongoDelUserRole(ctx context.Context, userid bson.ObjectID, projectid, rolename string) error {
	filter := bson.M{"_id": userid}
	updater := bson.M{"$pull": bson.M{"projects." + projectid: rolename}}
	r, e := d.mongo.Database("user").Collection("user").UpdateOne(ctx, filter, updater)
	if r.MatchedCount == 0 {
		e = ecode.ErrUserNotExist
	}
	return e
}
