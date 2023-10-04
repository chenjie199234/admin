package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty"` //user's id
	Mobile           string              `json:"mobile"`
	DingTalkUserName string              `json:"dingtalk_user_name"`
	FeiShuUserName   string              `json:"feishu_user_name"`
	Password         string              `bson:"password"` //only root user use this
	Projects         map[string][]string `bson:"projects"` //key projectid,value roles
}

type Role struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID string             `bson:"project_id"`
	RoleName  string             `bson:"role_name"`
	Comment   string             `bson:"comment"`
}
