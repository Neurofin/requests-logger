package models

import (
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentTypeModel struct {
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"descriptioin"`
	Flow        primitive.ObjectID `json:"flow" bson:"flow"`
	types.Timestamps
}

func BulkInsertDocumentTypes(docTypes []interface{}) (types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("documentType")

	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	if _, err := collection.InsertMany(context.Background(), docTypes); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	return result, nil
}
