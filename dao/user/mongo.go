package user

import (
	"context"
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

func (d *Dao) MongoLogin(ctx context.Context) (userid primitive.ObjectID, e error) {
	// TODO
	return primitive.NewObjectID(), nil
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
func (d *Dao) MongoSearchUsers(ctx context.Context, name string, limit, skip int64) (map[primitive.ObjectID]*model.User, int64, error) {
	totalsize, e := d.mongo.Database("user").Collection("user").CountDocuments(ctx, bson.M{"user_name": bson.M{"$regex": name}})
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
	cursor, e := d.mongo.Database("user").Collection("user").Find(ctx, bson.M{"user_name": bson.M{"$regex": name}}, opts)
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

func (d *Dao) MongoUpdateUser(ctx context.Context, userid primitive.ObjectID, name string, department []string) error {
	_, e := d.mongo.Database("user").Collection("user").UpdateOne(ctx, bson.M{"_id": userid}, bson.M{"$set": bson.M{"user_name": name, "department": department}})
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

func (d *Dao) MongoCreateRole(ctx context.Context, name, comment string) error {
	if _, e := d.mongo.Database("user").Collection("role").InsertOne(ctx, &model.Role{
		RoleName: name,
		Comment:  comment,
		Ctime:    uint32(time.Now().Unix()),
	}); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			return ecode.ErrRoleAlreadyExist
		}
		return e
	}
	return nil
}

func (d *Dao) MongoGetRoles(ctx context.Context, names []string) (map[string]*model.Role, error) {
	cursor, e := d.mongo.Database("user").Collection("role").Find(ctx, bson.M{"role_name": bson.M{"$in": names}})
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	result := make(map[string]*model.Role, cursor.RemainingBatchLength())
	for cursor.Next(ctx) {
		tmp := &model.Role{}
		if e := cursor.Decode(tmp); e != nil {
			return nil, e
		}
		result[tmp.RoleName] = tmp
	}
	return result, cursor.Err()
}

// if limit is 0 means all
func (d *Dao) MongoSearchRoles(ctx context.Context, name string, limit, skip int64) (map[string]*model.Role, int64, error) {
	totalsize, e := d.mongo.Database("user").Collection("role").CountDocuments(ctx, bson.M{"role_name": bson.M{"$regex": name}})
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
	cursor, e := d.mongo.Database("user").Collection("role").Find(ctx, bson.M{"role_name": bson.M{"$regex": name}}, opts)
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

func (d *Dao) MongoUpdateRole(ctx context.Context, name, comment string) error {
	_, e := d.mongo.Database("user").Collection("role").UpdateOne(ctx, bson.M{"role_name": name}, bson.M{"$set": bson.M{"comment": comment}})
	return e
}

func (d *Dao) MongoDelRoles(ctx context.Context, rolenames []string) (e error) {
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
	if _, e = d.mongo.Database("user").Collection("role").DeleteMany(sctx, bson.M{"role_name": bson.M{"$in": rolenames}}); e != nil {
		return
	}
	if _, e = d.mongo.Database("user").Collection("user").UpdateMany(sctx, bson.M{"roles": bson.M{"$in": rolenames}}, bson.M{"$pullAll": bson.M{"roles": rolenames}}); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, bson.M{"role_name": bson.M{"$in": rolenames}})
	return
}

func (d *Dao) MongoAddUserRole(ctx context.Context, userid primitive.ObjectID, rolename string) (e error) {
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
	exist, e = d.mongo.Database("user").Collection("role").CountDocuments(sctx, bson.M{"role_name": rolename})
	if e != nil {
		return
	}
	if exist == 0 {
		e = ecode.ErrRoleNotExist
		return
	}
	_, e = d.mongo.Database("user").Collection("user").UpdateOne(sctx, bson.M{"_id": userid}, bson.M{"$addToSet": bson.M{"roles": rolename}})
	return
}
func (d *Dao) MongoDelUserRole(ctx context.Context, userid primitive.ObjectID, rolename string) error {
	_, e := d.mongo.Database("user").Collection("user").UpdateOne(ctx, bson.M{"_id": userid}, bson.M{"$pull": bson.M{"roles": rolename}})
	return e
}
