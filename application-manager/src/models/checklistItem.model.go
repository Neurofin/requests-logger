package models

import (
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChecklistItemModel struct {
	Id                  primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name                string             `json:"name" bson:"name"`
	Goal                string             `json:"goal" bson:"goal"`
	Rules               []string           `json:"rules" bson:"rules"`
	Taxonomy            []string           `json:"taxonomy" bson:"taxonomy"`
	Prompt              string             `json:"prompt" bson:"prompt"`
	GroupUid            string             `json:"groupUid,omitempty" bson:"groupUid,omitempty"`
	RequiredDocs        []string           `json:"requiredDocs" bson:"requiredDocs"`
	Flow                primitive.ObjectID `json:"flow" bson:"flow"`
	MasterChecklistItem bool               `json:"masterChecklistItem,omitempty" bson:"masterChecklistItem,omitempty"`
	Engine              string             `json:"engine" bson:"engine"`
	types.Timestamps
}
