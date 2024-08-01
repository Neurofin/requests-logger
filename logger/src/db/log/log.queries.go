package logDb

import (
	"context"
	"logger/src/serverConfigs"
	"logger/src/store"
	loggerTypes "logger/src/store/types"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (i *Model) Insert() (loggerTypes.DbOperationResult, error) {

	operationResult := loggerTypes.DbOperationResult{}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.LogCollection)

	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()
	result, err := collection.InsertOne(context.Background(), i)
	if err != nil {
		return operationResult, err
	}

	operationResult.OperationSuccess = true
	operationResult.Data = result.InsertedID
	return operationResult, nil
}

func FetchById(id primitive.ObjectID) (loggerTypes.DbOperationResult, error) {
	operationResult := loggerTypes.DbOperationResult{}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.LogCollection)

	log := Model{}
	if err := collection.FindOne(context.Background(), bson.D{{Key: "_id", Value: id}}).Decode(&log); err != nil {
		return operationResult, err
	}

	operationResult.OperationSuccess = true
	operationResult.Data = log
	return operationResult, nil
}
