package models

import (
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"context"

	"go.mongodb.org/mongo-driver/bson"
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
	types.Timestamps
}

func (checklistItem *ChecklistItemModel) FetchFlowChecklist() (types.DbOperationResult, error) {

	result := types.DbOperationResult{}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("checklistItem")

	filter := bson.D{
		{
			Key:   "flow",
			Value: checklistItem.Flow,
		},
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return result, err
	}

	var data []ChecklistItemModel
	if err = cursor.All(context.TODO(), &data); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = data
	return result, nil
}
