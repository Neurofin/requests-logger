package models

import (
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OverrideChecklistResultLog struct {
	Id              primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Application     primitive.ObjectID     `json:"application" bson:"application"`
	ChecklistResult primitive.ObjectID     `json:"checklistResult" bson:"checklistResult"`
	OverrideMeta    map[string]interface{} `json:"overrideMeta" bson:"overrideMeta"`
	Note            string                 `json:"note" bson:"note"`
	User            primitive.ObjectID     `json:"user" bson:"user"`
	types.Timestamps
}

func (log *OverrideChecklistResultLog) InsertLog() (types.DbOperationResult, error) {

	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.OverrideChecklistResultLog)

	log.CreatedAt = time.Now()
	log.UpdatedAt = time.Now()
	operationResult, err := collection.InsertOne(context.Background(), log)
	if err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = operationResult.InsertedID
	return result, nil
}
