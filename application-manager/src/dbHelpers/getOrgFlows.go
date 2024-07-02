package dbHelpers

import (
	"application-manager/src/models"
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetOrgFlows(org primitive.ObjectID) (types.DbOperationResult, error) {
	result := types.DbOperationResult{}

	var data []models.FlowModel

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.FlowCollection)

	filter := bson.D{
		{
			Key:   "org",
			Value: org,
		},
	}
	cursor, err := collection.Find(context.Background(), filter, &options.FindOptions{
		Sort: bson.D{{
			Key:   "timestamps.createdAt",
			Value: -1,
		}},
	})
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
