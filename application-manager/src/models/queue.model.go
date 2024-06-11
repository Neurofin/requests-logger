package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QueueModel struct {
	Id                 primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Org                primitive.ObjectID     `json:"org,omitempty" bson:"org,omitempty"`
	Name               string                 `json:"name" bson:"name"`
	ClassifierId       string                 `json:"classifier_id" bson:"classifier_id"`
	TextractAdapterIds map[string]interface{} `json:"textract_adapter_ids,omitempty" bson:"textract_adapter_ids,omitempty"`
	Checklist          map[string]interface{} `json:"checklist,omitempty" bson:"checklist,omitempty"`
	Queries            map[string]interface{} `json:"queries,omitempty" bson:"queries,omitempty"`
}
