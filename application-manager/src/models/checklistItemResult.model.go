package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChecklistItemResultModel struct {
	Application primitive.ObjectID `json:"application" bson:"application"`
	ChecklistItem primitive.ObjectID `json:"checklistItem" bson:"checklistItem"`
	Result map[string]interface{} `json:"result" bson:"result"`
	Status string `json:"status" bson:"status"`
}
