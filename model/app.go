package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

// AppSummary and Log exist in same collection
// key=="" && index==0 => AppSummary
// key!="" && index!=0 => Log
// project+group+app+key+index add unique index
// permission_node_id add sparse index
type AppSummary struct {
	ID               bson.ObjectID          `bson:"_id,omitempty"`
	ProjectID        string                 `bson:"project_id"`
	ProjectName      string                 `bson:"project_name"`
	Group            string                 `bson:"group"`
	App              string                 `bson:"app"`
	Key              string                 `bson:"key"`           //this is always empty for Summary
	Index            uint32                 `bson:"index"`         //this is always 0 for Summary
	DiscoverMode     string                 `bson:"discover_mode"` //kubernetes,dns,static
	KubernetesNs     string                 `bson:"kubernetes_ns"`
	KubernetesLS     string                 `bson:"kubernetes_ls"`
	KubernetesFS     string                 `bson:"kubernetes_fs"`
	DnsHost          string                 `bson:"dns_host"`
	DnsInterval      uint32                 `bson:"dns_interval"` //unit second
	StaticAddrs      []string               `bson:"static_addrs"`
	CrpcPort         uint32                 `bson:"crpc_port"`
	CGrpcPort        uint32                 `bson:"cgrpc_port"`
	WebPort          uint32                 `bson:"web_port"`
	Keys             map[string]*KeySummary `bson:"keys"` //map's key is config's key name
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
type Log struct {
	ProjectID string `bson:"project_id"`
	Group     string `bson:"group"`
	App       string `bson:"app"`
	Key       string `bson:"key"`   //this is always not empty for Log
	Index     uint32 `bson:"index"` //this is always > 0  for Log
	Value     string `bson:"value"`
	ValueType string `bson:"value_type"`
}
