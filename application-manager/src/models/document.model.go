package models

import (
	"context"
	"errors"
	"time"

	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentModel struct {
	Id               primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string                 `json:"name" bson:"name"`
	Application      primitive.ObjectID     `json:"application" bson:"application"`
	Status           string                 `json:"status" bson:"status"`
	S3Location       string                 `json:"s3Location" bson:"s3Location"`
	TextractLocation string                 `json:"textractLocation" bson:"textractLocation"`
	ClassifierOutput map[string]interface{} `json:"classifierOutput" bson:"classifierOutput"`
	QueriesOutput    map[string]interface{} `json:"queriesOutput" bson:"queriesOutput"`
	types.Timestamps
}

func (doc *DocumentModel) InsertDocument() (*types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("document")

	// Check for existing document with the same name under the same application
	filter := bson.D{
		{Key: "name", Value: doc.Name},
		{Key: "application", Value: doc.Application},
	}

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return &types.DbOperationResult{OperationSuccess: false}, err
	}
	//logic to check duplicacy
	if count > 0 {
		return &types.DbOperationResult{OperationSuccess: false}, errors.New("document with this name already exists for this application")
	}

	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	_, err = collection.InsertOne(context.Background(), doc)
	if err != nil {
		return &types.DbOperationResult{OperationSuccess: false}, err
	}

	return &types.DbOperationResult{OperationSuccess: true}, nil
}
