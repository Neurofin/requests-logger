package models

import (
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChecklistItemModel struct {
	Id           primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Name         string              `json:"name" bson:"name"`
	Goal         string              `json:"goal" bson:"goal"`
	Rules        []string            `json:"rules" bson:"rules"`
	Taxonomy     []string            `json:"taxonomy" bson:"taxonomy"`
	Prompt       string              `json:"prompt" bson:"prompt"`
	Group        ChecklistGroupModel `json:"group,omitempty" bson:"group,omitempty"`
	RequiredDocs []string            `json:"requiredDocs" bson:"requiredDocs"`
	Flow         primitive.ObjectID  `json:"flow" bson:"flow"`
	types.Timestamps
}
