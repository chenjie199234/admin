package permission

import (
	"context"
	"math"
	"strconv"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (d *Dao) MongoGetUserPermission(ctx context.Context, userid primitive.ObjectID, nodeid []uint32) (canread, canwrite, admin bool, e error) {
	noderoute := make([][]uint32, 0, len(nodeid))
	for i := range nodeid {
		noderoute = append(noderoute, nodeid[:i+1])
	}
	usernodes, e := d.MongoGetUserNodes(ctx, userid, noderoute)
	if e != nil {
		return
	}
	canread, canwrite, admin = usernodes.CheckNode(nodeid)
	if admin {
		return
	}
	userrolenodes, e := d.MongoGetUserRoleNodes(ctx, userid, noderoute)
	if e != nil {
		return
	}
	for _, userrolenode := range userrolenodes {
		tmpread, tmpwrite, tmpadmin := userrolenode.CheckNode(nodeid)
		if tmpread {
			canread = tmpread
		}
		if tmpwrite {
			canwrite = tmpwrite
		}
		if tmpadmin {
			admin = tmpadmin
		}
		if admin {
			return
		}
	}
	return
}

func (d *Dao) MongoGetRolePermission(ctx context.Context, rolename string, nodeid []uint32) (canread, canwrite, admin bool, e error) {
	noderoute := make([][]uint32, 0, len(nodeid))
	for i := range nodeid {
		noderoute = append(noderoute, nodeid[:i+1])
	}
	rolenodes, e := d.MongoGetRoleNodes(ctx, rolename, noderoute)
	if e != nil {
		return
	}
	canread, canwrite, admin = rolenodes.CheckNode(nodeid)
	return
}

// if admin is true,canread and canwrite will be ignore
// if admin is false and canread is false too,means delete this user from this node
// if admin is false and canwrite is true,then canread must be tree too
func (d *Dao) MongoUpdateUserPermission(ctx context.Context, operateUserid, targetUserid primitive.ObjectID, nodeid []uint32, admin, canread, canwrite bool) (e error) {
	if len(nodeid) <= 1 {
		return ecode.ErrReq
	}
	if admin {
		//ignore
		canread = true
		canwrite = true
	} else if !canread && canwrite {
		e = ecode.ErrReq
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
	//first get target user permission on this node
	target := &model.UserNode{}
	if e = d.mongo.Database("permission").Collection("usernode").FindOne(sctx, bson.M{"user_id": targetUserid, "node_id": nodeid}).Decode(target); e != nil && e != mongo.ErrNoDocuments {
		return
	}
	if target.X == admin && target.R == canread && target.W == canwrite {
		//nothing need to do
		return
	}
	if !admin && !canread {
		//delete
		var x bool
		if target.X {
			//target is admin on this node
			//operator must be admin on this node's parent node
			_, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid[:len(nodeid)-1])
		} else {
			//target is not admin on this node
			//operator must be admin on this node
			_, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid)
		}
		if e != nil || !x {
			if e == nil {
				e = ecode.ErrPermission
			}
			return
		}
		//all check success,delete database
		_, e = d.mongo.Database("permission").Collection("usernode").DeleteOne(sctx, bson.M{
			"user_id": targetUserid,
			"node_id": nodeid,
		})
		return
	}
	//update

	//check target user exist
	var num int64
	if num, e = d.mongo.Database("user").Collection("user").CountDocuments(sctx, bson.M{"_id": targetUserid}); e != nil || num == 0 {
		if num == 0 {
			e = ecode.ErrUserNotExist
		}
		return
	}
	//get target user permission on parent path
	//if target is admin on parent path,nothing need to do
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, targetUserid, nodeid[:len(nodeid)-1]); e != nil || x {
		return
	}
	if admin || target.X {
		//want to give target admin or remove target's admin permission
		//operator must be admin to this node's parent node
		_, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid[:len(nodeid)-1])
	} else {
		//want to change target's R or W
		//operator must be admin on this node
		_, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid)
	}
	if e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success
	filter := bson.M{"user_id": targetUserid, "node_id": nodeid}
	updater := bson.M{"$set": bson.M{"r": canread, "w": canwrite, "x": admin}}
	if _, e = d.mongo.Database("permission").Collection("usernode").UpdateOne(sctx, filter, updater, options.Update().SetUpsert(true)); e != nil {
		return
	}
	if admin {
		//if target is admin on this node
		//clean all children permission
		filter = bson.M{"user_id": targetUserid}
		for i, v := range nodeid {
			filter["node_id."+strconv.Itoa(i)] = v
		}
		filter["node_id."+strconv.Itoa(len(nodeid))] = bson.M{"$exists": true}
		_, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, filter)
	}
	return
}

// if nodeids are not empty or nil,only the node in the required nodeids will return
func (d *Dao) MongoGetUserNodes(ctx context.Context, userid primitive.ObjectID, nodeids [][]uint32) (model.UserNodes, error) {
	filter := bson.M{"user_id": userid}
	if len(nodeids) > 0 {
		filter["node_id"] = bson.M{"$in": nodeids}
	}
	cursor, e := d.mongo.Database("permission").Collection("usernode").Find(ctx, filter)
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	result := make(model.UserNodes, 0, cursor.RemainingBatchLength())
	e = cursor.All(ctx, &result)
	return result, e
}

// if admin is true,canread and canwrite will be ignore
// if admin is false and canread is false too,means delete this user from this node
// if admin is false and canwrite is true,then canread must be tree too
func (d *Dao) MongoUpdateRolePermission(ctx context.Context, operateUserid primitive.ObjectID, rolename string, nodeid []uint32, admin, canread, canwrite bool) (e error) {
	if len(nodeid) == 0 {
		return ecode.ErrReq
	}
	if admin {
		//ignore
		canread = true
		canwrite = true
	} else if !canread && canwrite {
		e = ecode.ErrReq
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
	//first get target role permission on this node
	target := &model.UserNode{}
	if e = d.mongo.Database("permission").Collection("rolenode").FindOne(sctx, bson.M{"role_name": rolename, "node_id": nodeid}).Decode(target); e != nil && e != mongo.ErrNoDocuments {
		return
	}
	if target.R == canread && target.W == canwrite && target.X == admin {
		//nothing need todo
		return
	}
	if !admin && !canread {
		//delete
		var x bool
		if target.X {
			//target is admin on this node
			//operator must be admin on this node's parent node
			_, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid[:len(nodeid)-1])
		} else {
			//target is not admin on this node
			//operator must be admin on this node
			_, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid)
		}
		if e != nil || !x {
			if e == nil {
				e = ecode.ErrPermission
			}
			return
		}
		_, e = d.mongo.Database("permission").Collection("rolenode").DeleteOne(sctx, bson.M{
			"role_name": rolename,
			"node_id":   nodeid,
		})
		return
	}
	//update

	//check target role exist
	var num int64
	if num, e = d.mongo.Database("user").Collection("role").CountDocuments(sctx, bson.M{"role_name": rolename}); e != nil || num == 0 {
		if num == 0 {
			e = ecode.ErrRoleNotExist
		}
		return
	}
	//get target role permission on parent path
	//if target is admin on parent path,nothing need to do
	var x bool
	if _, _, x, e = d.MongoGetRolePermission(sctx, rolename, nodeid[:len(nodeid)-1]); e != nil || x {
		return
	}
	if x || target.X {
		//want to give target admin or remove target's admin permission
		//operator must be admin to this node's parent node
		_, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid[:len(nodeid)-1])
	} else {
		//want to change target's R or W
		//operator must be admin on this node
		_, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid)
	}
	if e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success
	filter := bson.M{"role_name": rolename, "node_id": nodeid}
	updater := bson.M{"$set": bson.M{"r": canread, "w": canwrite, "x": admin}}
	if _, e = d.mongo.Database("permission").Collection("rolenode").UpdateOne(sctx, filter, updater, options.Update().SetUpsert(true)); e != nil {
		return
	}
	if admin {
		//if target is admin on this node
		//clean all children permission
		filter = bson.M{"role_name": rolename}
		for i, v := range nodeid {
			filter["node_id."+strconv.Itoa(i)] = v
		}
		filter["node_id."+strconv.Itoa(len(nodeid))] = bson.M{"$exists": true}
		_, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, filter)
	}
	return
}

// if nodeids are not empty or nil,only the node in the required nodeids will return
func (d *Dao) MongoGetRoleNodes(ctx context.Context, rolename string, nodeids [][]uint32) (model.RoleNodes, error) {
	filter := bson.M{
		"role_name": rolename,
	}
	if len(nodeids) > 0 {
		filter["node_id"] = bson.M{"$in": nodeids}
	}
	cursor, e := d.mongo.Database("permission").Collection("rolename").Find(ctx, filter)
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	result := make(model.RoleNodes, 0, cursor.RemainingBatchLength())
	e = cursor.All(ctx, &result)
	return result, e
}

// if nodeids are not empty or nil,only the node in the required nodeids will return
func (d *Dao) MongoGetUserRoleNodes(ctx context.Context, userid primitive.ObjectID, nodeids [][]uint32) (map[string]model.RoleNodes, error) {
	r := d.mongo.Database("user").Collection("user").FindOne(ctx, bson.M{"_id": userid})
	if e := r.Err(); e != nil {
		if r.Err() == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, e
	}
	userinfo := &model.User{}
	if e := r.Decode(userinfo); e != nil {
		return nil, e
	}
	filter := bson.M{
		"role_name": bson.M{"$in": userinfo.Roles},
	}
	if len(nodeids) > 0 {
		filter["node_id"] = bson.M{"$in": nodeids}
	}
	cursor, e := d.mongo.Database("permission").Collection("rolenode").Find(ctx, filter)
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	tmp := make([]*model.RoleNode, 0, cursor.RemainingBatchLength())
	if e = cursor.All(ctx, &tmp); e != nil {
		return nil, e
	}
	result := make(map[string]model.RoleNodes)
	for _, v := range tmp {
		if _, ok := result[v.RoleName]; !ok {
			result[v.RoleName] = make(model.RoleNodes, 0, 10)
		}
		result[v.RoleName] = append(result[v.RoleName], v)
	}
	return result, nil
}

// get one specific node's info
func (d *Dao) MongoGetNode(ctx context.Context, nodeid []uint32) (*model.Node, error) {
	if len(nodeid) == 0 {
		return nil, nil
	}
	node := &model.Node{}
	if e := d.mongo.Database("permission").Collection("node").FindOne(ctx, bson.M{"node_id": nodeid}).Decode(node); e != nil {
		if e == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, e
	}
	return node, nil
}

// get all specific nodes' info in the nodeids
func (d *Dao) MongoGetNodes(ctx context.Context, nodeids [][]uint32) ([]*model.Node, error) {
	if len(nodeids) == 0 {
		return nil, nil
	}
	cursor, e := d.mongo.Database("permission").Collection("node").Find(ctx, bson.M{"node_id": bson.M{"$in": nodeids}})
	if e != nil {
		return nil, e
	}
	defer cursor.Close(context.Background())
	nodes := make([]*model.Node, 0, cursor.RemainingBatchLength())
	e = cursor.All(ctx, &nodes)
	return nodes, e
}

// get one specific node's children,if pnodeid is empty or nil,return all nodes
func (d *Dao) MongoListNode(ctx context.Context, pnodeid []uint32) ([]*model.Node, error) {
	filter := bson.M{}
	for i, v := range pnodeid {
		filter["node_id."+strconv.Itoa(i)] = v
	}
	filter["node_id."+strconv.Itoa(len(pnodeid))] = bson.M{"$exists": true}
	cursor, e := d.mongo.Database("permission").Collection("node").Find(ctx, filter)
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	nodes := make([]*model.Node, 0, cursor.RemainingBatchLength())
	e = cursor.All(ctx, &nodes)
	return nodes, e
}
func (d *Dao) MongoAddNode(ctx context.Context, operateUserid primitive.ObjectID, pnodeid []uint32, name, data string) (e error) {
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
	//check parent exist
	parent := &model.Node{}
	e = d.mongo.Database("permission").Collection("node").FindOne(sctx, bson.M{"node_id": pnodeid}).Decode(parent)
	if e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrPNodeNotExist
		}
		return
	}
	//check admin
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, pnodeid); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success,modify database
	if _, e = d.mongo.Database("permission").Collection("node").InsertOne(sctx, &model.Node{
		NodeId:       append(pnodeid, parent.CurNodeIndex+1),
		NodeName:     name,
		NodeData:     data,
		CurNodeIndex: 0,
	}); e != nil {
		return
	}
	if _, e = d.mongo.Database("permission").Collection("node").UpdateOne(sctx, bson.M{"node_id": pnodeid}, bson.M{"$inc": bson.M{"cur_node_index": 1}}); e != nil {
		return
	}
	return
}
func (d *Dao) MongoUpdateNode(ctx context.Context, operateUserid primitive.ObjectID, nodeid []uint32, name, data string) (e error) {
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
	//check admin
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success,update database
	r, e := d.mongo.Database("permission").Collection("node").UpdateOne(sctx, bson.M{"node_id": nodeid}, bson.M{"$set": bson.M{"node_name": name, "node_data": data}})
	if e == nil && r.MatchedCount == 0 {
		e = ecode.ErrNodeNotExist
	}
	return
}
func (d *Dao) MongoMoveNode(ctx context.Context, operateUserid primitive.ObjectID, nodeid, pnodeid []uint32) (e error) {
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
	//check self exist
	var self int64
	self, e = d.mongo.Database("permission").Collection("node").CountDocuments(sctx, bson.M{"node_id": nodeid})
	if e != nil {
		return
	}
	if self == 0 {
		e = ecode.ErrNodeNotExist
		return
	}
	//check parent exist
	parent := &model.Node{}
	if e = d.mongo.Database("permission").Collection("node").FindOne(sctx, bson.M{"node_id": pnodeid}).Decode(parent); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrPNodeNotExist
		}
		return
	}
	//check admin in current path
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//check admin in new path
	if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, pnodeid); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//update the new parent
	if _, e = d.mongo.Database("permission").Collection("node").UpdateOne(sctx, bson.M{"node_id": pnodeid}, bson.M{"$inc": bson.M{"cur_node_index": 1}}); e != nil {
		return
	}
	newnodeid := append(parent.NodeId, parent.CurNodeIndex+1)
	filter := bson.M{}
	for i, v := range nodeid {
		filter["node_id."+strconv.Itoa(i)] = v
	}
	updater := bson.A{bson.M{"$set": bson.M{"node_id": bson.M{"$concatArrays": bson.A{newnodeid, bson.M{"$slice": bson.A{"$node_id", len(nodeid), math.MaxInt32}}}}}}}
	//update the node
	if _, e = d.mongo.Database("permission").Collection("node").UpdateMany(sctx, filter, updater); e != nil {
		return
	}
	//update the usernode
	if _, e = d.mongo.Database("permission").Collection("usernode").UpdateMany(sctx, filter, updater); e != nil {
		return
	}
	//update the rolenode
	if _, e = d.mongo.Database("permission").Collection("rolenode").UpdateMany(sctx, filter, updater); e != nil {
		return
	}
	return
}
func (d *Dao) MongoDeleteNode(ctx context.Context, operateUserid primitive.ObjectID, nodeid []uint32) (e error) {
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
	//check admin
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operateUserid, nodeid); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success,delete database
	delfilter := bson.M{}
	for i, v := range nodeid {
		delfilter["node_id."+strconv.Itoa(i)] = v
	}
	if _, e = d.mongo.Database("permission").Collection("node").DeleteMany(sctx, delfilter); e != nil {
		return
	}
	if _, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, delfilter); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, delfilter)
	return
}
