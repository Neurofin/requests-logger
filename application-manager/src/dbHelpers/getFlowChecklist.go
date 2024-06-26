package dbHelpers

import (
	"application-manager/src/models"
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetFlowChecklist(flow primitive.ObjectID) (types.DbOperationResult, error) {

	result := types.DbOperationResult{}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ChecklistItemCollection)

	filter := bson.D{
		{
			Key:   "flow",
			Value: flow,
		},
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return result, err
	}

	var data []models.ChecklistItemModel
	if err = cursor.All(context.TODO(), &data); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = data
	return result, nil
}
