package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID  `bson:"_id"` //user's id
	UserName string              `bson:"user_name"`
	Password string              `bson:"password"`
	Projects map[string][]string `bson:"projects"` //key projectid,value roles
}

type Role struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID string             `bson:"project_id"`
	RoleName  string             `bson:"role_name"`
	Comment   string             `bson:"comment"`
}
