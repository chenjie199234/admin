package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id"` //user's id
	UserName   string             `bson:"user_name"`
	Password   string             `bson:"password"`
	Department []string           `bson:"department"`
	ProjectIDs []string           `bson:"project_ids"` //element is project
	Roles      []string           `bson:"roles"`       //element is project:rolename
}

type Role struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProjectID string             `bson:"project_id"`
	RoleName  string             `bson:"role_name"`
	Comment   string             `bson:"comment"`
}
