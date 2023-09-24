package model

import (
	"sort"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Node struct {
	NodeId       string `bson:"node_id"` //self's id
	NodeName     string `bson:"node_name"`
	NodeData     string `bson:"node_data"`
	CurNodeIndex uint32 `bson:"cur_node_index"` //auto increment,this is for child's node_id
}
type ProjectIndex struct {
	ProjectName string `bson:"project_name"`
	ProjectID   string `bson:"project_id"`
}
type UserNode struct {
	UserId primitive.ObjectID `bson:"user_id"`
	NodeId string             `bson:"node_id"`
	R      bool               `bson:"r"`
	W      bool               `bson:"w"`
	X      bool               `bson:"x"`
}
type RoleNode struct {
	ProjectID string `bson:"project_id"`
	RoleName  string `bson:"role_name"`
	NodeId    string `bson:"node_id"`
	R         bool   `bson:"r"`
	W         bool   `bson:"w"`
	X         bool   `bson:"x"`
}
type UserNodes []*UserNode

func (u UserNodes) Sort() {
	sort.Slice(u, func(i, j int) bool {
		return strings.Count(u[i].NodeId, ",") < strings.Count(u[j].NodeId, ",")
	})
}

func (u UserNodes) CheckNode(nodeid string) (canread, canwrite, admin bool) {
	for _, usernode := range u {
		if strings.Count(usernode.NodeId, ",") > strings.Count(nodeid, ",") {
			return false, false, false
		}
		if !strings.HasPrefix(nodeid+",", usernode.NodeId+",") {
			continue
		}
		//check admin
		if usernode.X {
			return true, true, true
		}
		if usernode.NodeId == nodeid {
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

func (r RoleNodes) Sort() {
	sort.Slice(r, func(i, j int) bool {
		return strings.Count(r[i].NodeId, ",") < strings.Count(r[j].NodeId, ",")
	})
}
func (r RoleNodes) CheckNode(nodeid string) (canread, canwrite, admin bool) {
	for _, rolenode := range r {
		if strings.Count(rolenode.NodeId, ",") > strings.Count(nodeid, ",") {
			return false, false, false
		}
		if !strings.HasPrefix(nodeid+",", rolenode.NodeId+",") {
			continue
		}
		//check admin
		if rolenode.X {
			return true, true, true
		}
		if rolenode.NodeId == nodeid {
			//this is the target node
			if !rolenode.R {
				return false, false, false
			}
			return rolenode.R, rolenode.W, false
		}
	}
	return false, false, false
}
