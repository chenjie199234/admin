package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID         primitive.ObjectID `bson:"_id"` //user's id
	Name       string             `bson:"name"`
	Department []string           `bson:"department"`
	Ctime      uint32             `bson:"ctime"`
}
