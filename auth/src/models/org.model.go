package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"auth/src/serverConfigs"
	"auth/src/store"
	"auth/src/store/types"
)


type OrgModel struct {
	Id   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	types.Timestamps
}

func (org *OrgModel) InsertOrg() (*types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("org")

	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()
	_, err := collection.InsertOne(context.Background(), org)

	if err !=nil {
		result := &types.DbOperationResult{
			OperationSuccess: false,
		}
		return result, err
	}

	result := &types.DbOperationResult{
		OperationSuccess: true,
	}
	return result, err
}

func (org *OrgModel) UpdateOrg() (*types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("org")

	org.UpdatedAt = time.Now()
	_, err := collection.UpdateByID(context.Background(), org.Id, org)
	if err !=nil {
		result := &types.DbOperationResult{
			OperationSuccess: false,
		}
		return result, err
	}

	result := &types.DbOperationResult{
		OperationSuccess: true,
	}
	return result, err
}

func (org *OrgModel) GetOrg() (*types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("org")

	var orgDoc OrgModel
	err := collection.FindOne(context.Background(), org).Decode(&orgDoc)

	if err !=nil {
		result := &types.DbOperationResult{
			OperationSuccess: false,
		}
		return result, err
	}

	result := &types.DbOperationResult{
		OperationSuccess: true,
		Data: orgDoc,
	}
	return result, err
}
