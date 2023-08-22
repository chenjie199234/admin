package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// AppSummary and Log exist in same collection
// key=="" && index==0 => AppSummary
// key!="" && index!=0 => Log
// project+group+app+key+index add unique index
// permission_node_id add sparse index
type AppSummary struct {
	ID               primitive.ObjectID     `bson:"_id,omitempty"`
	ProjectID        string                 `bson:"project_id"`
	Group            string                 `bson:"group"`
	App              string                 `bson:"app"`
	Key              string                 `bson:"key"`   //this is always empty for Summary
	Index            uint32                 `bson:"index"` //this is always 0 for Summary
	Paths            map[string]*ProxyPath  `bson:"paths"` //map's key is the base64(proxy path)
	Keys             map[string]*KeySummary `bson:"keys"`  //map's key is config's key name
	Value            string                 `bson:"value"`
	PermissionNodeID string                 `bson:"permission_node_id"`
}

func (a *AppSummary) GetFullName() string {
	return a.ProjectID + "-" + a.Group + "." + a.App
}

type KeySummary struct {
	CurIndex     uint32 `bson:"cur_index"`
	MaxIndex     uint32 `bson:"max_index"`
	CurVersion   uint32 `bson:"cur_version"`
	CurValue     string `bson:"cur_value"`
	CurValueType string `bson:"cur_value_type"`
}
type ProxyPath struct {
	PermissionNodeID string `bson:"permission_node_id"`
	PermissionRead   bool   `bson:"permission_read"`
	PermissionWrite  bool   `bson:"permission_write"`
	PermissionAdmin  bool   `bson:"permission_admin"`
}
type Log struct {
	ProjectID string `bson:"project_id"`
	Group     string `bson:"group"`
	App       string `bson:"app"`
	Key       string `bson:"key"`   //this is always not empty for Log
	Index     uint32 `bson:"index"` //this is always > 0  for Log
	Value     string `bson:"value"`
	ValueType string `bson:"value_type"`
}
