package model

import (
	"sort"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Node struct {
	NodeId       []uint32 `bson:"node_id"` //self's id
	NodeName     string   `bson:"node_name"`
	NodeData     string   `bson:"node_data"`
	CurNodeIndex uint32   `bson:"cur_node_index"` //auto increment,this is for child's node_id
}
type UserNode struct {
	UserId primitive.ObjectID `bson:"user_id"`
	NodeId []uint32           `bson:"node_id"`
	R      bool               `bson:"r"`
	W      bool               `bson:"w"`
	X      bool               `bson:"x"`
}
type RoleNode struct {
	RoleName string   `bson:"role_name"`
	NodeId   []uint32 `bson:"node_id"`
	R        bool     `bson:"r"`
	W        bool     `bson:"w"`
	X        bool     `bson:"x"`
}
type UserNodes []*UserNode

func (u UserNodes) CheckNode(nodeid []uint32) (canread, canwrite, admin bool) {
	sort.Slice(u, func(i, j int) bool {
		return len(u[i].NodeId) < len(u[j].NodeId)
	})
	for _, usernode := range u {
		if len(usernode.NodeId) > len(nodeid) {
			return false, false, false
		}
		isprefix := true
		for i := range usernode.NodeId {
			if usernode.NodeId[i] != nodeid[i] {
				isprefix = false
				break
			}
		}
		if !isprefix {
			continue
		}
		//check admin
		if usernode.X {
			return true, true, true
		}
		if len(usernode.NodeId) == len(nodeid) {
			//this is the target node
			if !usernode.R {
				return false, false, false
			}
			return usernode.R, usernode.W, false
		}
	}
	return false, false, false
}

type RoleNodes []*RoleNode

func (r RoleNodes) CheckNode(nodeid []uint32) (canread, canwrite, admin bool) {
	sort.Slice(r, func(i, j int) bool {
		return len(r[i].NodeId) < len(r[j].NodeId)
	})
	for _, rolenode := range r {
		if len(rolenode.NodeId) > len(nodeid) {
			return false, false, false
		}
		isprefix := true
		for i := range rolenode.NodeId {
			if rolenode.NodeId[i] != nodeid[i] {
				isprefix = false
				break
			}
		}
		if !isprefix {
			continue
		}
		//check admin
		if rolenode.X {
			return true, true, true
		}
		if len(rolenode.NodeId) == len(nodeid) {
			//this is the target node
			if !rolenode.R {
				return false, false, false
			}
			return rolenode.R, rolenode.W, false
		}
	}
	return false, false, false
}
