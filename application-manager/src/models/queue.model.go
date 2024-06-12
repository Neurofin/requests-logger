package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueModel struct {
	Id                 primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Org                primitive.ObjectID     `json:"org,omitempty" bson:"org,omitempty"`
	Name               string                 `json:"name" bson:"name"`
	ClassifierId       string                 `json:"classifierId" bson:"classifierId"`
	TextractAdapterIds map[string]interface{} `json:"textractAdapterIds,omitempty" bson:"textractAdapterIds,omitempty"`
	Checklist          map[string]interface{} `json:"checklist,omitempty" bson:"checklist,omitempty"`
	Queries            map[string]interface{} `json:"queries,omitempty" bson:"queries,omitempty"`
}
