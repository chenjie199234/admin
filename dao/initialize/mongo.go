package initialize

import (
	"context"
	"strconv"
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

func (d *Dao) MongoInit(ctx context.Context, password string) (e error) {
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
	if _, e = d.mongo.Database("user").Collection("user").InsertOne(sctx, &model.User{
		ID:         primitive.NilObjectID,
		UserName:   "superadmin",
		Password:   password,
		Department: []string{},
		Ctime:      uint32(time.Now().Unix()),
		Roles:      []string{},
	}); e != nil && !mongo.IsDuplicateKeyError(e) {
		return
	} else if e != nil {
		e = ecode.ErrAlreadyInited
		return
	}
	if _, e = d.mongo.Database("permission").Collection("node").InsertOne(sctx, &model.Node{
		NodeId:       "0",
		NodeName:     "root",
		NodeData:     "",
		CurNodeIndex: 0,
	}); e != nil && !mongo.IsDuplicateKeyError(e) {
		return
	} else if e != nil {
		e = ecode.ErrAlreadyInited
		return
	}
	if _, e = d.mongo.Database("permission").Collection("usernode").InsertOne(sctx, &model.UserNode{
		UserId: primitive.NilObjectID,
		NodeId: "0",
		R:      true,
		W:      true,
		X:      true,
	}); e != nil && !mongo.IsDuplicateKeyError(e) {
		return
	} else if e != nil {
		e = ecode.ErrAlreadyInited
		return
	}
	return
}
func (d *Dao) MongoRootLogin(ctx context.Context) (*model.User, error) {
	user := &model.User{}
	if e := d.mongo.Database("user").Collection("user").FindOne(ctx, bson.M{"_id": primitive.NilObjectID}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrNotInited
		}
		return nil, e
	}
	return user, nil
}
func (d *Dao) MongoRootPassword(ctx context.Context, oldpassword, newpassword string) error {
	r, e := d.mongo.Database("user").Collection("user").UpdateOne(ctx, bson.M{"_id": primitive.NilObjectID, "password": oldpassword}, bson.M{"$set": bson.M{"password": newpassword}})
	if e == nil && r.MatchedCount == 0 {
		e = ecode.ErrOldPasswordWrong
	}
	return e
}
func (d *Dao) MongoCreateProject(ctx context.Context, projectname, projectdata string) (e error) {
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
	root := &model.Node{}
	if e = d.mongo.Database("permission").Collection("node").FindOneAndUpdate(sctx, bson.M{"node_id": "0"}, bson.M{"$inc": bson.M{"cur_node_index": 1}}).Decode(root); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrNotInited
		}
		return
	}
	nodeid := "0," + strconv.FormatUint(uint64(root.CurNodeIndex+1), 10)
	docs := bson.A{}
	docs = append(docs, &model.Node{
		NodeId:       nodeid,
		NodeName:     projectname,
		NodeData:     projectdata,
		CurNodeIndex: 2,
	})
	docs = append(docs, &model.Node{
		NodeId:       nodeid + model.UserControl,
		NodeName:     "UserControl",
		NodeData:     "",
		CurNodeIndex: 0,
	})
	docs = append(docs, &model.Node{
		NodeId:       nodeid + model.RoleControl,
		NodeName:     "RoleControl",
		NodeData:     "",
		CurNodeIndex: 0,
	})
	_, e = d.mongo.Database("permission").Collection("node").InsertMany(sctx, docs)
	return
}
func (d *Dao) MongoListProject(ctx context.Context) ([]*model.Node, error) {
	cur, e := d.mongo.Database("permission").Collection("node").Find(ctx, bson.M{"node_id": bson.M{"$regex": "^0,[1-9][0-9]*$"}})
	if e != nil {
		return nil, e
	}
	defer cur.Close(ctx)
	result := make([]*model.Node, 0, cur.RemainingBatchLength())
	e = cur.All(ctx, &result)
	return result, e
}
