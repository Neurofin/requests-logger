package utils

import (
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"context"
)

func BulkInsertToDb(items []interface{}, collectionName string) (types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(collectionName)

	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	if _, err := collection.InsertMany(context.Background(), items); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	return result, nil
}
