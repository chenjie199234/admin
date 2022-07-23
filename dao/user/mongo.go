package user

import (
	"context"

	"github.com/chenjie199234/admin/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (d *Dao) MongoSearchUsers(ctx context.Context, name string, limit int64) (map[primitive.ObjectID]*model.User, error) {
	cursor, e := d.mongo.Database("user").Collection("user").Find(ctx, bson.M{"name": bson.M{"$regex": name}}, options.Find().SetLimit(limit))
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
	return result, nil
}
