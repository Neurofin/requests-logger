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

func GetOrgApplications(org primitive.ObjectID, page int, pageSize int) (types.DbOperationResult, int64, error) {
    result := types.DbOperationResult{}

    var data []models.ApplicationModel

    collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ApplicationCollection)

    filter := bson.D{
        {
            Key:   "org",
            Value: org,
        },
    }

    totalCount, err := collection.CountDocuments(context.Background(), filter)
    if err != nil {
        return result, 0, err
    }

    // The offset for pagination
    skip := int64((page - 1) * pageSize)
    limit := int64(pageSize)

    // Adjust skip if it exceeds the total count
    if skip >= totalCount {
        skip = totalCount - limit
        if skip < 0 {
            skip = 0
        }
    }

    cursor, err := collection.Find(context.Background(), filter, &options.FindOptions{
        Sort: bson.D{{
            Key:   "timestamps.createdAt",
            Value: -1,
        }},
        Skip:  &skip,
        Limit: &limit,
    })
    if err != nil {
        return result, 0, err
    }

    if err = cursor.All(context.TODO(), &data); err != nil {
        return result, 0, err
    }

    result.OperationSuccess = true
    result.Data = data
    return result, totalCount, nil
}
