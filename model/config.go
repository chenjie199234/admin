package model

import "go.mongodb.org/mongo-driver/bson/primitive"

/*
| config_groupname1(database)
|      appname1(collection)
|      appname2
|      appname3
| config_groupname2(database)
|      appnameN(collection)
*/
//every collection has two kinds of data
type AppSummary struct {
	ID               primitive.ObjectID     `bson:"_id"`
	Key              string                 `bson:"key"`   //this is always empty
	Index            uint32                 `bson:"index"` //this is always 0
	Paths            map[string]*ProxyPath  `bson:"paths"` //map's key is the base64(proxy path)
	Keys             map[string]*KeySummary `bson:"keys"`  //map's key is config's key name
	Value            string                 `bson:"value"`
	PermissionNodeID string                 `bson:"permission_node_id"`
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
}
type Log struct {
	Key       string `bson:"key"`   //this is always not empty
	Index     uint32 `bson:"index"` //this is always > 0  for Config
	Value     string `bson:"value"`
	ValueType string `bson:"value_type"`
}
