package initialize

import (
	"context"
	"crypto/rand"
	"strconv"
	"strings"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/secure"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func (d *Dao) MongoInit(ctx context.Context, password string) (e error) {
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
	sign, _ := secure.SignMake(password)
	if _, e = d.mongo.Database("user").Collection("user").InsertOne(sctx, bson.M{"_id": bson.NilObjectID, "password": sign, "projects": map[string][]string{}}); e != nil && !mongo.IsDuplicateKeyError(e) {
		return
	} else if e != nil {
		e = ecode.ErrAlreadyInited
		return
	}
	if _, e = d.mongo.Database("permission").Collection("usernode").InsertOne(sctx, &model.UserNode{
		UserId: bson.NilObjectID,
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
	if e := d.mongo.Database("user").Collection("user").FindOne(ctx, bson.M{"_id": bson.NilObjectID}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrNotInited
		}
		return nil, e
	}
	return user, nil
}
func (d *Dao) MongoUpdateRootPassword(ctx context.Context, oldpassword, newpassword string) (e error) {
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
	nonce := make([]byte, 16)
	rand.Read(nonce)
	user := &model.User{}
	filter := bson.M{"_id": bson.NilObjectID}
	sign, _ := secure.SignMake(newpassword)
	updater := bson.M{"password": sign}
	if e = d.mongo.Database("user").Collection("user").FindOneAndUpdate(sctx, filter, bson.M{"$set": updater}).Decode(user); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrNotInited
		}
		return
	}
	e = secure.SignCheck(oldpassword, user.Password)
	return
}
func (d *Dao) MongoCreateProject(ctx context.Context, projectname, projectdata string) (projectid string, e error) {
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
	root := &model.Node{}
	if e = d.mongo.Database("permission").Collection("node").FindOneAndUpdate(sctx, bson.M{"node_id": "0"}, bson.M{"$inc": bson.M{"cur_node_index": 1}}).Decode(root); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrNotInited
		}
		return
	}
	projectid = "0," + strconv.FormatUint(uint64(root.CurNodeIndex+1), 10)
	if _, e = d.mongo.Database("permission").Collection("projectindex").InsertOne(sctx, bson.M{"project_name": projectname, "project_id": projectid}); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			e = ecode.ErrProjectAlreadyExist
		}
		return
	}
	docs := bson.A{}
	docs = append(docs, &model.Node{
		NodeId:       projectid,
		NodeName:     projectname,
		NodeData:     projectdata,
		CurNodeIndex: 100,
	})
	docs = append(docs, &model.Node{
		NodeId:       projectid + model.UserAndRoleControl,
		NodeName:     "UserAndRoleControl",
		NodeData:     "",
		CurNodeIndex: 0,
	})
	docs = append(docs, &model.Node{
		NodeId:       projectid + model.AppControl,
		NodeName:     "AppControl",
		NodeData:     "",
		CurNodeIndex: 0,
	})
	_, e = d.mongo.Database("permission").Collection("node").InsertMany(sctx, docs)
	return
}
func (d *Dao) MongoGetProjectIDByName(ctx context.Context, projectname string) (string, error) {
	projectindex := &model.ProjectIndex{}
	if e := d.mongo.Database("permission").Collection("projectindex").FindOne(ctx, bson.M{"project_name": projectname}).Decode(projectindex); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrProjectNotExist
		}
		return "", e
	}
	return projectindex.ProjectID, nil
}
func (d *Dao) MongoGetProjectNameByID(ctx context.Context, projectid string) (string, error) {
	projectindex := &model.ProjectIndex{}
	if e := d.mongo.Database("permission").Collection("projectindex").FindOne(ctx, bson.M{"project_id": projectid}).Decode(projectindex); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrProjectNotExist
		}
		return "", e
	}
	return projectindex.ProjectName, nil
}
func (d *Dao) MongoUpdateProject(ctx context.Context, projectid, newname, newdata string) (node *model.Node, e error) {
	if !strings.HasPrefix(projectid, "0,") || strings.Count(projectid, ",") != 1 {
		return nil, ecode.ErrReq
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
	nameupdater := bson.M{"$set": bson.M{"project_name": newname}}
	if r, e = d.mongo.Database("permission").Collection("projectindex").UpdateOne(sctx, bson.M{"project_id": projectid}, nameupdater); e != nil {
		if mongo.IsDuplicateKeyError(e) {
			e = ecode.ErrProjectAlreadyExist
		}
		return
	}
	if r.MatchedCount == 0 {
		e = ecode.ErrProjectNotExist
		return
	}
	var samename bool
	if r.ModifiedCount == 0 {
		samename = true
	}
	node = &model.Node{}
	nodeupdater := bson.M{"$set": bson.M{"node_name": newname, "node_data": newdata}}
	if e = d.mongo.Database("permission").Collection("node").FindOneAndUpdate(sctx, bson.M{"node_id": projectid}, nodeupdater).Decode(node); e != nil {
		return
	}
	if !samename {
		if _, e = d.mongo.Database("app").Collection("config").UpdateOne(sctx, bson.M{"key": "", "index": 0, "project_id": projectid}, nameupdater); e != nil {
			return
		}
	}
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
func (d *Dao) MongoDelProject(ctx context.Context, projectid string) (node *model.Node, e error) {
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
	node = &model.Node{}
	if e = d.mongo.Database("permission").Collection("node").FindOneAndDelete(sctx, bson.M{"node_id": projectid}).Decode(node); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrProjectNotExist
		}
		return
	}
	if _, e = d.mongo.Database("permission").Collection("projectindex").DeleteOne(sctx, bson.M{"project_id": projectid}); e != nil {
		return
	}
	if _, e = d.mongo.Database("permission").Collection("node").DeleteMany(sctx, bson.M{"node_id": bson.M{"$regex": "^" + projectid}}); e != nil {
		return
	}
	if _, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, bson.M{"node_id": bson.M{"$regex": "^" + projectid}}); e != nil {
		return
	}
	if _, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, bson.M{"project_id": projectid}); e != nil {
		return
	}
	if _, e = d.mongo.Database("user").Collection("user").UpdateMany(sctx, bson.M{}, bson.M{"$unset": bson.M{"projects." + projectid: 1}}); e != nil {
		return
	}
	if _, e = d.mongo.Database("user").Collection("role").DeleteMany(sctx, bson.M{"project_id": projectid}); e != nil {
		return
	}
	_, e = d.mongo.Database("app").Collection("config").DeleteMany(sctx, bson.M{"project_id": projectid})
	return
}
