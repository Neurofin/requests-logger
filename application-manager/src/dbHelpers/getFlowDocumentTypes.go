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

func GetFlowDocumentTypes(flow primitive.ObjectID) (types.DbOperationResult, error) {

	result := types.DbOperationResult{}

	var data []models.DocumentTypeModel

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.DocumentTypeCollection)

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

	if err = cursor.All(context.TODO(), &data); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = data
	return result, nil
}
