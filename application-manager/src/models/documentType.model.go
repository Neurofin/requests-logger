package models

import (
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentTypeModel struct {
	Uid         string             `json:"uid" bson:"uid"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description,omitempty" bson:"descriptioin,omitempty"`
	Flow        primitive.ObjectID `json:"flow" bson:"flow"`
	types.Timestamps
}
