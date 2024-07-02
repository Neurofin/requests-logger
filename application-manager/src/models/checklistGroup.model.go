package models

import (
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChecklistGroupModel struct {
	Uid  string             `json:"uid" bson:"uid"`
	Name string             `json:"name" bson:"name"`
	Goal string             `json:"goal" bson:"goal"`
	Flow primitive.ObjectID `json:"flow" bson:"flow"`
	types.Timestamps
}
