package permission

import (
	"context"
	"strconv"
	"strings"

	"github.com/chenjie199234/admin/ecode"
	"github.com/chenjie199234/admin/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readconcern"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func (d *Dao) MongoGetUserPermission(ctx context.Context, userid bson.ObjectID, nodeid string, withrole bool) (canread, canwrite, admin bool, e error) {
	if userid.IsZero() {
		return true, true, true, nil
	}
	if !strings.HasPrefix(nodeid, "0,") {
		return false, false, false, nil
	}
	copynodeid := nodeid
	noderoute := make([]string, 0, strings.Count(nodeid, ","))
	for copynodeid != "" {
		index := strings.LastIndex(copynodeid, ",")
		if index == -1 {
			break
		}
		noderoute = append(noderoute, copynodeid)
		copynodeid = copynodeid[:index]
	}
	usernodes, e := d.MongoGetUserNodes(ctx, userid, noderoute[len(noderoute)-1], noderoute)
	if e != nil {
		return
	}
	canread, canwrite, admin = usernodes.CheckNode(nodeid)
	if admin {
		return
	}
	if !withrole {
		return
	}
	userrolenodes, e := d.MongoGetUserRoleNodes(ctx, userid, noderoute[len(noderoute)-1], noderoute)
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

func (d *Dao) MongoGetRolePermission(ctx context.Context, projectid, rolename, nodeid string) (canread, canwrite, admin bool, e error) {
	if strings.HasPrefix(projectid, "0,") || strings.Count(projectid, ",") != 1 || strings.HasPrefix(nodeid, "0,") {
		return false, false, false, nil
	}
	noderoute := make([]string, 0, strings.Count(nodeid, ","))
	for nodeid != "" {
		index := strings.LastIndex(nodeid, ",")
		if index == -1 {
			break
		}
		noderoute = append(noderoute, nodeid)
		nodeid = nodeid[:index]
	}
	rolenodes, e := d.MongoGetRoleNodes(ctx, projectid, rolename, noderoute)
	if e != nil {
		return
	}
	canread, canwrite, admin = rolenodes.CheckNode(nodeid)
	return
}

// if admin is true,canread and canwrite will be ignore
// if admin is false and canread is false too,means delete this user from this node
// if admin is false and canwrite is true,then canread must be tree too
func (d *Dao) MongoUpdateUserPermission(ctx context.Context, operator, target bson.ObjectID, nodeid string, admin, canread, canwrite bool) (e error) {
	if !strings.HasPrefix(nodeid, "0,") {
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
	//first get target user permission on this node
	targetnode := &model.UserNode{}
	if e = d.mongo.Database("permission").Collection("usernode").FindOne(sctx, bson.M{"user_id": target, "node_id": nodeid}).Decode(targetnode); e != nil && e != mongo.ErrNoDocuments {
		return
	}
	if targetnode.X == admin && targetnode.R == canread && targetnode.W == canwrite {
		//nothing need to do
		return
	}
	if !admin && !canread {
		//delete
		var x bool
		if targetnode.X {
			//target is admin on this node
			//operator must be admin on this node's parent node
			lastindex := strings.LastIndex(nodeid, ",")
			_, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid[:lastindex], true)
		} else {
			//target is not admin on this node
			//operator must be admin on this node
			_, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid, true)
		}
		if e != nil || !x {
			if e == nil {
				e = ecode.ErrPermission
			}
			return
		}
		//all check success,delete database
		_, e = d.mongo.Database("permission").Collection("usernode").DeleteOne(sctx, bson.M{
			"user_id": target,
			"node_id": nodeid,
		})
		return
	}
	//update

	//check target user exist
	targetuser := &model.User{}
	if e = d.mongo.Database("user").Collection("user").FindOne(sctx, bson.M{"_id": target}).Decode(targetuser); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrUserNotExist
		}
		return
	}
	inproject := false
	for userprojectid := range targetuser.Projects {
		if strings.HasPrefix(nodeid, userprojectid) {
			inproject = true
			break
		}
	}
	if !inproject {
		e = ecode.ErrUserNotInProject
		return
	}
	//get target user permission on parent path
	//if target is admin on parent path,nothing need to do
	lastindex := strings.LastIndex(nodeid, ",")
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, target, nodeid[:lastindex], false); e != nil || x {
		return
	}
	if admin || targetnode.X {
		//want to give target admin or remove target's admin permission
		//operator must be admin to this node's parent node
		_, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid[:lastindex], true)
	} else {
		//want to change target's R or W
		//operator must be admin on this node
		_, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid, true)
	}
	if e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success
	filter := bson.M{"user_id": target, "node_id": nodeid}
	updater := bson.M{"$set": bson.M{"r": canread, "w": canwrite, "x": admin}}
	if _, e = d.mongo.Database("permission").Collection("usernode").UpdateOne(sctx, filter, updater, options.UpdateOne().SetUpsert(true)); e != nil {
		return
	}
	if admin {
		//if target is admin on this node
		//clean all children permission
		filter = bson.M{"user_id": target, "node_id": bson.M{"$regex": "^" + nodeid + ","}}
		_, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, filter)
	}
	return
}

// if nodeids are not empty or nil,only the node in the required nodeids will return
func (d *Dao) MongoGetUserNodes(ctx context.Context, userid bson.ObjectID, projectid string, nodeids []string) (model.UserNodes, error) {
	filter := bson.M{"user_id": userid}
	nodeidfilter := bson.M{"$regex": "^" + projectid}
	if len(nodeids) > 0 {
		nodeidfilter["$in"] = nodeids
	}
	filter["node_id"] = nodeidfilter
	cursor, e := d.mongo.Database("permission").Collection("usernode").Find(ctx, filter)
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	result := make(model.UserNodes, 0, cursor.RemainingBatchLength())
	e = cursor.All(ctx, &result)
	result.Sort()
	return result, e
}

// if admin is true,canread and canwrite will be ignore
// if admin is false and canread is false too,means delete this user from this node
// if admin is false and canwrite is true,then canread must be tree too
func (d *Dao) MongoUpdateRolePermission(ctx context.Context, operator bson.ObjectID, projectid, rolename string, nodeid string, admin, canread, canwrite bool) (e error) {
	//role belong's to project,so the nodeid must belong to this project
	if !strings.HasPrefix(projectid, "0,") || strings.Count(projectid, ",") != 1 || strings.Count(nodeid, ",") < 2 || !strings.HasPrefix(nodeid+",", projectid+",") {
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
	//first get target role permission on this node
	targetnode := &model.UserNode{}
	if e = d.mongo.Database("permission").Collection("rolenode").FindOne(sctx, bson.M{"project_id": projectid, "role_name": rolename, "node_id": nodeid}).Decode(targetnode); e != nil && e != mongo.ErrNoDocuments {
		return
	}
	if targetnode.R == canread && targetnode.W == canwrite && targetnode.X == admin {
		//nothing need todo
		return
	}
	if !admin && !canread {
		//delete
		var x bool
		if targetnode.X {
			//target is admin on this node
			//operator must be admin on this node's parent node
			lastindex := strings.LastIndex(nodeid, ",")
			_, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid[:lastindex], true)
		} else {
			//target is not admin on this node
			//operator must be admin on this node
			_, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid, true)
		}
		if e != nil || !x {
			if e == nil {
				e = ecode.ErrPermission
			}
			return
		}
		_, e = d.mongo.Database("permission").Collection("rolenode").DeleteOne(sctx, bson.M{
			"project_id": projectid,
			"role_name":  rolename,
			"node_id":    nodeid,
		})
		return
	}
	//update

	//check target role exist
	var num int64
	if num, e = d.mongo.Database("user").Collection("role").CountDocuments(sctx, bson.M{"project_id": projectid, "role_name": rolename}); e != nil || num == 0 {
		if num == 0 {
			e = ecode.ErrRoleNotExist
		}
		return
	}
	//get target role permission on parent path
	//if target is admin on parent path,nothing need to do
	lastindex := strings.LastIndex(nodeid, ",")
	var x bool
	if _, _, x, e = d.MongoGetRolePermission(sctx, projectid, rolename, nodeid[:lastindex]); e != nil || x {
		return
	}
	if admin || targetnode.X {
		//want to give target admin or remove target's admin permission
		//operator must be admin to this node's parent node
		_, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid[:lastindex], true)
	} else {
		//want to change target's R or W
		//operator must be admin on this node
		_, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid, true)
	}
	if e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success
	filter := bson.M{"project_id": projectid, "role_name": rolename, "node_id": nodeid}
	updater := bson.M{"$set": bson.M{"r": canread, "w": canwrite, "x": admin}}
	if _, e = d.mongo.Database("permission").Collection("rolenode").UpdateOne(sctx, filter, updater, options.UpdateOne().SetUpsert(true)); e != nil {
		return
	}
	if admin {
		//if target is admin on this node
		//clean all children permission
		filter = bson.M{"project_id": projectid, "role_name": rolename, "node_id": bson.M{"$regex": "^" + nodeid + ","}}
		_, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, filter)
	}
	return
}

// if nodeids are not empty or nil,only the node in the required nodeids will return
func (d *Dao) MongoGetRoleNodes(ctx context.Context, projectid, rolename string, nodeids []string) (model.RoleNodes, error) {
	filter := bson.M{
		"project_id": projectid,
		"role_name":  rolename,
	}
	exist, e := d.mongo.Database("user").Collection("role").CountDocuments(ctx, filter)
	if e != nil {
		return nil, e
	}
	if exist == 0 {
		return nil, ecode.ErrRoleNotExist
	}
	if len(nodeids) > 0 {
		filter["node_id"] = bson.M{"$in": nodeids}
	}
	cursor, e := d.mongo.Database("permission").Collection("rolenode").Find(ctx, filter)
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	result := make(model.RoleNodes, 0, cursor.RemainingBatchLength())
	e = cursor.All(ctx, &result)
	result.Sort()
	return result, e
}

// if nodeids are not empty or nil,only the node in the required nodeids will return
func (d *Dao) MongoGetUserRoleNodes(ctx context.Context, userid bson.ObjectID, projectid string, nodeids []string) (map[string]model.RoleNodes, error) {
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
	if len(userinfo.Projects) == 0 {
		return nil, nil
	}
	userprojectroles, ok := userinfo.Projects[projectid]
	if !ok || len(userprojectroles) == 0 {
		return nil, nil
	}
	or := bson.A{}
	for _, userprojectrole := range userprojectroles {
		or = append(or, bson.M{
			"project_id": projectid,
			"role_name":  userprojectrole,
		})
	}
	filter := bson.M{"$or": or}
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
	for _, rolenodes := range result {
		rolenodes.Sort()
	}
	return result, nil
}

// get one specific node's info
func (d *Dao) MongoGetNode(ctx context.Context, nodeid string) (*model.Node, error) {
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
func (d *Dao) MongoGetNodes(ctx context.Context, nodeids []string) ([]*model.Node, error) {
	if len(nodeids) == 0 {
		return nil, nil
	}
	cursor, e := d.mongo.Database("permission").Collection("node").Find(ctx, bson.M{"node_id": bson.M{"$in": nodeids}})
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	nodes := make([]*model.Node, 0, cursor.RemainingBatchLength())
	e = cursor.All(ctx, &nodes)
	return nodes, e
}

// get children of pnodeid,pnodeid will not return
// if pnodeid is empty,pnodeid is the root's nodeid
// all
//   - true,get all children(include children's children)
//   - false,get the direct children
func (d *Dao) MongoListChildrenNodes(ctx context.Context, pnodeid string, all bool) ([]*model.Node, error) {
	filter := bson.M{}
	regex := ""
	if pnodeid != "" {
		regex = "^" + pnodeid + ","
	} else {
		regex = "^0,"
	}
	if !all {
		regex += "[1-9][0-9]*$"
	}
	filter["node_id"] = bson.M{"$regex": regex}
	cursor, e := d.mongo.Database("permission").Collection("node").Find(ctx, filter)
	if e != nil {
		return nil, e
	}
	defer cursor.Close(ctx)
	nodes := make([]*model.Node, 0, cursor.RemainingBatchLength())
	e = cursor.All(ctx, &nodes)
	return nodes, e
}
func (d *Dao) MongoAddNode(ctx context.Context, operator bson.ObjectID, pnodeid string, name, data string) (nodeid string, e error) {
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
	//check admin
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operator, pnodeid, true); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//check parent exist
	parent := &model.Node{}
	if e = d.mongo.Database("permission").Collection("node").FindOneAndUpdate(sctx, bson.M{"node_id": pnodeid}, bson.M{"$inc": bson.M{"cur_node_index": 1}}).Decode(parent); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrPNodeNotExist
		}
		return
	}
	//all check success,modify database
	nodeid = pnodeid + "," + strconv.FormatUint(uint64(parent.CurNodeIndex+1), 10)
	_, e = d.mongo.Database("permission").Collection("node").InsertOne(sctx, &model.Node{
		NodeId:       nodeid,
		NodeName:     name,
		NodeData:     data,
		CurNodeIndex: 0,
	})
	return
}
func (d *Dao) MongoUpdateNode(ctx context.Context, operator bson.ObjectID, nodeid string, newname, newdata string) (node *model.Node, e error) {
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
	//check admin
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid, true); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success,update database
	node = &model.Node{}
	e = d.mongo.Database("permission").Collection("node").FindOneAndUpdate(sctx, bson.M{"node_id": nodeid}, bson.M{"$set": bson.M{"node_name": newname, "node_data": newdata}}).Decode(node)
	if e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrNodeNotExist
		}
		return
	}
	return
}
func (d *Dao) MongoMoveNode(ctx context.Context, operator bson.ObjectID, nodeid, pnodeid string) (newnodeid string, e error) {
	//nodeid and pnodeid must in same project
	if !strings.HasPrefix(nodeid, "0,") || strings.Count(nodeid, ",") == 1 || !strings.HasPrefix(pnodeid, "0,") {
		return "", ecode.ErrReq
	}
	if index := strings.Index(pnodeid[2:], ","); index == -1 {
		if pnodeid[2:] != nodeid[2:][:strings.Index(nodeid[2:], ",")] {
			return "", ecode.ErrReq
		}
	} else {
		if pnodeid[2:][:strings.Index(pnodeid[2:], ",")] != nodeid[2:][:strings.Index(nodeid[2:], ",")] {
			return "", ecode.ErrReq
		}
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
	if _, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid, true); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//check admin in new path
	if _, _, x, e = d.MongoGetUserPermission(sctx, operator, pnodeid, true); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//update the new parent
	if _, e = d.mongo.Database("permission").Collection("node").UpdateOne(sctx, bson.M{"node_id": pnodeid}, bson.M{"$inc": bson.M{"cur_node_index": 1}}); e != nil {
		return
	}
	newnodeid = parent.NodeId + "," + strconv.FormatUint(uint64(parent.CurNodeIndex+1), 10)
	filter := bson.M{"node_id": bson.M{"$regex": "^" + nodeid}}
	updater := bson.A{bson.M{
		"$set": bson.M{
			"node_id": bson.M{
				"$concat": bson.A{
					newnodeid,
					bson.M{
						"$substrBytes": bson.A{
							"$node_id",
							len(nodeid),
							bson.M{"$strLenBytes": "$node_id"},
						},
					},
				},
			},
		},
	}}
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
func (d *Dao) MongoDeleteNode(ctx context.Context, operator bson.ObjectID, nodeid string) (node *model.Node, e error) {
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
	//check admin
	var x bool
	if _, _, x, e = d.MongoGetUserPermission(sctx, operator, nodeid, true); e != nil || !x {
		if e == nil {
			e = ecode.ErrPermission
		}
		return
	}
	//all check success,delete database
	node = &model.Node{}
	if e = d.mongo.Database("permission").Collection("node").FindOneAndDelete(sctx, bson.M{"node_id": nodeid}).Decode(node); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrNodeNotExist
		}
		return
	}
	delfilter := bson.M{"node_id": bson.M{"$regex": "^" + nodeid}}
	if _, e = d.mongo.Database("permission").Collection("node").DeleteMany(sctx, delfilter); e != nil {
		return
	}
	if _, e = d.mongo.Database("permission").Collection("usernode").DeleteMany(sctx, delfilter); e != nil {
		return
	}
	_, e = d.mongo.Database("permission").Collection("rolenode").DeleteMany(sctx, delfilter)
	return
}
