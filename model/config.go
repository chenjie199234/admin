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
	ID     primitive.ObjectID     `bson:"_id"`
	Key    string                 `bson:"key"`   //this is always empty
	Index  uint32                 `bson:"index"` //this is always 0
	Cipher string                 `bson:"cipher"`
	Keys   map[string]*KeySummary `bson:"keys"` //map's key is config's key name
}
type KeySummary struct {
	CurIndex     uint32 `bson:"cur_index"`
	MaxIndex     uint32 `bson:"max_index"`
	CurVersion   uint32 `bson:"cur_version"`
	CurValue     string `bson:"cur_value"` //if Cipher is not empty,this field is encrypt
	CurValueType string `bson:"cur_value_type"`
}
type Log struct {
	Key       string `bson:"key"`   //this is always not empty
	Index     uint32 `bson:"index"` //this is always > 0  for Config
	Value     string `bson:"value"`
	ValueType string `bson:"value_type"`
}
