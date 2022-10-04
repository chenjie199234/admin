package initialize

import (
	"context"
	"time"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

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
	docs := make([]interface{}, 0, 3)
	docs = append(docs, &model.Node{
		NodeId:       model.RootNodeId,
		NodeName:     "root",
		NodeData:     "",
		CurNodeIndex: 2,
	})
	docs = append(docs, &model.Node{
		NodeId:       model.UserControlNodeId,
		NodeName:     "UserControl",
		NodeData:     "",
		CurNodeIndex: 0,
	})
	docs = append(docs, &model.Node{
		NodeId:       model.RoleControlNodeId,
		NodeName:     "RoleControl",
		NodeData:     "",
		CurNodeIndex: 0,
	})
	if _, e = d.mongo.Database("permission").Collection("node").InsertMany(sctx, docs); e != nil && !mongo.IsDuplicateKeyError(e) {
		return
	} else if e != nil {
		e = ecode.ErrAlreadyInited
		return
	}
	if _, e = d.mongo.Database("permission").Collection("usernode").InsertOne(sctx, &model.UserNode{
		UserId: primitive.NilObjectID,
		NodeId: []uint32{0},
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
