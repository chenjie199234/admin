package user

import (
	"context"
	"strings"
	"time"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (d *Dao) MongoUserLogin(ctx context.Context) (userid primitive.ObjectID, e error) {
	// TODO
	return primitive.NewObjectID(), nil
}

func (d *Dao) MongoInvite(ctx context.Context, operator primitive.ObjectID, project string, target primitive.ObjectID) (e error) {
	if target == primitive.NilObjectID {
		return ecode.ErrReq
	}
	if !strings.HasPrefix(project, "0,") || strings.Count(project, ",") != 1 {
		return ecode.ErrReq
	}
	var r *mongo.UpdateResult
	r, e = d.mongo.Database("user").Collection("user").UpdateOne(ctx, bson.M{"_id": target}, bson.M{"$addToSet": bson.M{"projects": project}})
	if e != nil {
		return
	}
	if r.MatchedCount == 0 {
		e = ecode.ErrUserNotExist
		return
	} else if r.ModifiedCount == 0 {
		e = ecode.ErrUserAlreadyInvited
		return
	}
	return
}
func (d *Dao) MongoKick(ctx context.Context, operator primitive.ObjectID, project string, target primitive.ObjectID) (e error) {
	if target == primitive.NilObjectID {
		return ecode.ErrReq
	}
	if !strings.HasPrefix(project, "0,") || strings.Count(project, ",") != 1 {
		return ecode.ErrReq
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
	if _, e = d.mongo.Database("user").Collection("user").UpdateOne(sctx, bson.M{"_id": target}, bson.M{"$pull": bson.M{"projects": project, "roles": bson.M{"$regex": "^" + project}}}); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, bson.M{"user_id": target, "node_id": bson.M{"$regex": "^" + project}})
	return
}

func (d *Dao) MongoGetUsers(ctx context.Context, userids []primitive.ObjectID) (map[primitive.ObjectID]*model.User, error) {
	cursor, e := d.mongo.Database("user").Collection("user").Find(ctx, bson.M{"_id": bson.M{"$in": userids}})
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	result := make(map[primitive.ObjectID]*model.User, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		tmp := &model.User{}
		if e := cursor.Decode(tmp); e != nil {
			return nil, e
		}
		result[tmp.ID] = tmp
	}
	return result, cursor.Err()
}

// if limit is 0 means all
func (d *Dao) MongoSearchUsers(ctx context.Context, project, name string, limit, skip int64) (map[primitive.ObjectID]*model.User, int64, error) {
	filter := bson.M{"user_name": bson.M{"$regex": name}}
	if project != "" {
		filter["projects"] = project
	}
	totalsize, e := d.mongo.Database("user").Collection("user").CountDocuments(ctx, filter)
	if e != nil {
		return nil, 0, e
	}
	opts := options.Find().SetSort(bson.M{"ctime": -1})
	if skip != 0 {
		opts = opts.SetSkip(skip)
	}
	if limit != 0 {
		opts = opts.SetLimit(limit)
	}
	cursor, e := d.mongo.Database("user").Collection("user").Find(ctx, filter, opts)
	if e != nil {
		return nil, 0, e
	}
	defer cursor.Close(ctx)
	result := make(map[primitive.ObjectID]*model.User, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		tmp := &model.User{}
		if e := cursor.Decode(tmp); e != nil {
			return nil, 0, e
		}
		result[tmp.ID] = tmp
	}
	return result, totalsize, cursor.Err()
}

func (d *Dao) MongoUpdateUser(ctx context.Context, userid primitive.ObjectID, newname string, newdepartment []string) error {
	_, e := d.mongo.Database("user").Collection("user").UpdateOne(ctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"user_name": newname, "department": newdepartment}})
	return e
}

func (d *Dao) MongoDelUsers(ctx context.Context, userids []primitive.ObjectID) (e error) {
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
	if _, e = d.mongo.Database("user").Collection("user").DeleteMany(sctx, bson.M{"_id": bson.M{"$in": userids}}); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, bson.M{"user_id": bson.M{"$in": userids}})
	return
}

func (d *Dao) MongoCreateRole(ctx context.Context, project, name, comment string) (e error) {
	if _, e = d.mongo.Database("user").Collection("role").InsertOne(ctx, &model.Role{
		Project:  project,
		RoleName: name,
		Comment:  comment,
		Ctime:    uint32(time.Now().Unix()),
	}); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			e = ecode.ErrRoleAlreadyExist
		}
		return e
	}
	return
}

// if limit is 0 means all
func (d *Dao) MongoSearchRoles(ctx context.Context, project, name string, limit, skip int64) (map[string]*model.Role, int64, error) {
	totalsize, e := d.mongo.Database("user").Collection("role").CountDocuments(ctx, bson.M{"project": project, "role_name": bson.M{"$regex": name}})
	if e != nil {
		return nil, 0, e
	}
	opts := options.Find().SetSort(bson.M{"ctime": -1})
	if skip != 0 {
		opts = opts.SetSkip(skip)
	}
	if limit != 0 {
		opts = opts.SetLimit(limit)
	}
	cursor, e := d.mongo.Database("user").Collection("role").Find(ctx, bson.M{"project": project, "role_name": bson.M{"$regex": name}}, opts)
	if e != nil {
		return nil, 0, e
	}
	defer cursor.Close(ctx)
	result := make(map[string]*model.Role, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		tmp := &model.Role{}
		if e := cursor.Decode(tmp); e != nil {
			return nil, 0, e
		}
		result[tmp.RoleName] = tmp
	}
	return result, totalsize, cursor.Err()
}

func (d *Dao) MongoUpdateRole(ctx context.Context, project, name, comment string) error {
	_, e := d.mongo.Database("user").Collection("role").UpdateOne(ctx, bson.M{"project": project, "role_name": name}, bson.M{"$set": bson.M{"comment": comment}})
	return e
}

func (d *Dao) MongoDelRoles(ctx context.Context, project string, rolenames []string) (e error) {
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
	if _, e = d.mongo.Database("user").Collection("role").DeleteMany(sctx, bson.M{"project": project, "role_name": bson.M{"$in": rolenames}}); e != nil {
		return
	}
	in := []string{}
	for _, rolename := range rolenames {
		in = append(in, project+":"+rolename)
	}
	if _, e = d.mongo.Database("user").Collection("user").UpdateMany(sctx, bson.M{"roles": bson.M{"$in": in}}, bson.M{"$pullAll": bson.M{"roles": in}}); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, bson.M{"project": project, "role_name": bson.M{"$in": rolenames}})
	return
}

func (d *Dao) MongoAddUserRole(ctx context.Context, userid primitive.ObjectID, project, rolename string) (e error) {
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
	var exist int64
	exist, e = d.mongo.Database("user").Collection("role").CountDocuments(sctx, bson.M{"project": project, "role_name": rolename})
	if e != nil {
		return
	}
	if exist == 0 {
		e = ecode.ErrRoleNotExist
		return
	}
	var r *mongo.UpdateResult
	if r, e = d.mongo.Database("user").Collection("user").UpdateOne(sctx, bson.M{"_id": userid, "projects": project}, bson.M{"$addToSet": bson.M{"roles": project + ":" + rolename}}); e != nil {
		return
	}
	if r.MatchedCount == 0 {
		e = ecode.ErrUserNotInProject
	}
	return
}
func (d *Dao) MongoDelUserRole(ctx context.Context, userid primitive.ObjectID, project, rolename string) error {
	_, e := d.mongo.Database("user").Collection("user").UpdateOne(ctx, bson.M{"_id": userid}, bson.M{"$pull": bson.M{"roles": project + ":" + rolename}})
	return e
}
