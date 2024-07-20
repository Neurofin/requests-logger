package models

import (
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChecklistItemResultModel struct {
	Id            primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Application   primitive.ObjectID     `json:"application" bson:"application"`
	ChecklistItem primitive.ObjectID     `json:"checklistItem" bson:"checklistItem"`
	Result        map[string]interface{} `json:"result,omitempty" bson:"result,omitempty"`
	Overridden    bool                   `json:"overridden" bson:"overridden"`
	OverrideMeta  map[string]interface{} `json:"overrideMeta" bson:"overrideMeta"`
	types.Timestamps
}

func (queryResult *ChecklistItemResultModel) FindChecklistItemResult() (types.DbOperationResult, error) {
	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("checklistItemResult")

	data := ChecklistItemResultModel{}
	if err := collection.FindOne(context.Background(), bson.D{{
		Key:   "application",
		Value: queryResult.Application,
	}, {
		Key:   "checklistItem",
		Value: queryResult.ChecklistItem,
	}}).Decode(&data); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = data
	return result, nil
}

func (queryResult *ChecklistItemResultModel) InsertChecklistItemResult() (types.DbOperationResult, error) {
	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("checklistItemResult")

	queryResult.CreatedAt = time.Now()
	queryResult.UpdatedAt = time.Now()
	if _, err := collection.InsertOne(context.Background(), queryResult); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	return result, nil
}

func (queryResult *ChecklistItemResultModel) UpdateChecklistItemResult() (types.DbOperationResult, error) {
	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("checklistItemResult")

	queryResult.UpdatedAt = time.Now()
	if _, err := collection.UpdateOne(context.Background(), bson.D{{
		Key:   "_id",
		Value: queryResult.Id,
	}}, bson.D{{
		Key:   "$set",
		Value: queryResult,
	}},
	); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	return result, nil
}

func (queryResult *ChecklistItemResultModel) FindChecklistItemResultById() (types.DbOperationResult, error) {
	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("checklistItemResult")

	data := ChecklistItemResultModel{}
	if err := collection.FindOne(context.Background(), bson.D{{
		Key:   "_id",
		Value: queryResult.Id,
	}}).Decode(&data); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = data
	return result, nil
}
