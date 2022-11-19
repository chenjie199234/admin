package initialize

import (
	"context"
	"crypto/rand"
	"strconv"
	"time"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (d *Dao) MongoInit(ctx context.Context, password string) (e error) {
	if len(password) >= 32 {
		return ecode.ErrPasswordLength
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
	nonce := make([]byte, 16)
	rand.Read(nonce)
	if _, e = d.mongo.Database("user").Collection("user").InsertOne(sctx, &model.User{
		ID:         primitive.NilObjectID,
		UserName:   "superadmin",
		Password:   util.SignMake(password, nonce),
		Department: []string{},
		Ctime:      uint32(time.Now().Unix()),
		Roles:      []string{},
		Projects:   []string{},
	}); e != nil && !mongo.IsDuplicateKeyError(e) {
		return
	} else if e != nil {
		e = ecode.ErrAlreadyInited
		return
	}
	docs := bson.A{}
	//root node
	docs = append(docs, &model.Node{
		NodeId:       "0",
		NodeName:     "root",
		NodeData:     "",
		CurNodeIndex: 1,
	})
	//project admin's node
	docs = append(docs, &model.Node{
		NodeId:       "0,1",
		NodeName:     "admin",
		NodeData:     "",
		CurNodeIndex: 3,
	})
	//project admin's user control node
	docs = append(docs, &model.Node{
		NodeId:       "0,1" + model.UserControl,
		NodeName:     "UserControl",
		NodeData:     "",
		CurNodeIndex: 0,
	})
	//project admin's role control node
	docs = append(docs, &model.Node{
		NodeId:       "0,1" + model.RoleControl,
		NodeName:     "RoleControl",
		NodeData:     "",
		CurNodeIndex: 0,
	})
	//project admin's config control node
	docs = append(docs, &model.Node{
		NodeId:       "0,1" + model.ConfigControl,
		NodeName:     "ConfigControl",
		NodeData:     "",
		CurNodeIndex: 1,
	})
	docs = append(docs, &model.Node{
		NodeId:       "0,1" + model.ConfigControl + ",1",
		NodeName:     model.Group + "." + model.Name,
		NodeData:     "",
		CurNodeIndex: 0,
	})
	if _, e = d.mongo.Database("permission").Collection("node").InsertMany(sctx, docs); e != nil && !mongo.IsDuplicateKeyError(e) {
		return
	} else if e != nil {
		e = ecode.ErrAlreadyInited
		return
	}
	if _, e = d.mongo.Database("config_"+model.Group).Collection(model.Name).UpdateOne(sctx, bson.M{"key": "", "index": 0}, bson.M{"$set": bson.M{"permission_node_id": "0,1" + model.ConfigControl + ",1"}}); e != nil {
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
func (d *Dao) MongoRootPassword(ctx context.Context, oldpassword, newpassword string) (e error) {
	if len(oldpassword) >= 32 || len(newpassword) >= 32 {
		return ecode.ErrPasswordLength
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
	nonce := make([]byte, 16)
	rand.Read(nonce)
	user := &model.User{}
	filter := bson.M{"_id": primitive.NilObjectID}
	updater := bson.M{"password": util.SignMake(newpassword, nonce)}
	if e = d.mongo.Database("user").Collection("user").FindOneAndUpdate(ctx, filter, bson.M{"$set": updater}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrNotInited
		}
		return
	}
	if e = util.SignCheck(oldpassword, user.Password); e != nil {
		e = ecode.ErrPasswordWrong
	}
	return
}
func (d *Dao) MongoCreateProject(ctx context.Context, projectname, projectdata string) (nodeid string, e error) {
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
	nodeid = "0," + strconv.FormatUint(uint64(root.CurNodeIndex+1), 10)
	docs := bson.A{}
	docs = append(docs, &model.Node{
		NodeId:       nodeid,
		NodeName:     projectname,
		NodeData:     projectdata,
		CurNodeIndex: 3,
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
	docs = append(docs, &model.Node{
		NodeId:       nodeid + model.ConfigControl,
		NodeName:     "ConfigControl",
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
func (d *Dao) MongoDelProject(ctx context.Context, projectid string) (e error) {
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
	if _, e = d.mongo.Database("user").Collection("user").UpdateMany(sctx, bson.M{}, bson.M{"$pull": bson.M{"projects": projectid, "roles": bson.M{"$regex": "^" + projectid}}}); e != nil {
		return
	}
	if _, e = d.mongo.Database("user").Collection("role").DeleteMany(sctx, bson.M{"project": projectid}); e != nil {
		return
	}
	if _, e = d.mongo.Database("permission").Collection("node").DeleteMany(sctx, bson.M{"node_id": bson.M{"$regex": "^" + projectid}}); e != nil {
		return
	}
	if _, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, bson.M{"node_id": bson.M{"$regex": "^" + projectid}}); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, bson.M{"project": projectid})
	return
}
