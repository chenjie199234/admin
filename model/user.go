package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id"` //user's id
	UserName   string             `bson:"user_name"`
	Password   string             `bson:"password"`
	Department []string           `bson:"department"`
	Ctime      uint32             `bson:"ctime"`
	Roles      []string           `bson:"roles"`
}

type Role struct {
	RoleName string `bson:"role_name"`
	Comment  string `bson:"comment"`
	Ctime    uint32 `bson:"ctime"`
}
