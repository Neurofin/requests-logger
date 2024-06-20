package models

import (
	"context"
	"errors"
	"time"

	"application-manager/src/serverConfigs"
	classifierServiceTypes "application-manager/src/services/classifier/store/types"
	"application-manager/src/store"
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentModel struct {
	Id               primitive.ObjectID               `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string                           `json:"name" bson:"name"`
	Application      primitive.ObjectID               `json:"application" bson:"application"`
	Status           string                           `json:"status" bson:"status"`
	S3Location       string                           `json:"s3Location" bson:"s3Location"`
	TextractLocation string                           `json:"textractLocation" bson:"textractLocation"`
	ClassifierOutput classifierServiceTypes.ClassData `json:"classifierOutput" bson:"classifierOutput"`
	QueriesOutput    map[string]interface{}           `json:"queriesOutput" bson:"queriesOutput"`
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

func (doc *DocumentModel) GetDocument() (*types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("document")

	filter := bson.D{{
		Key: "$or",
		Value: bson.A{
			bson.D{{
				Key:   "s3Location",
				Value: doc.S3Location,
			}},
			bson.D{{
				Key:   "name",
				Value: doc.Name,
			}},
			bson.D{{
				Key:   "_id",
				Value: doc.Id,
			}},
		},
	}}

	document := DocumentModel{}
	err := collection.FindOne(context.Background(), filter).Decode(&document)
	if err != nil {
		result := &types.DbOperationResult{
			OperationSuccess: false,
		}
		return result, err
	}

	result := &types.DbOperationResult{
		OperationSuccess: true,
		Data:             document,
	}
	return result, err
}

func (doc *DocumentModel) UpdateDocument() (*types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("document")

	_, err := collection.UpdateByID(context.Background(), doc.Id, bson.D{{Key: "$set", Value: bson.D{
		{Key: "classifierOutput", Value: doc.ClassifierOutput},
		{Key: "queriesOutput", Value: doc.QueriesOutput},
		{Key: "s3Location", Value: doc.S3Location},
		{Key: "status", Value: doc.Status},
		{Key: "textractLocation", Value: doc.TextractLocation},
		{Key: "timestamps.updatedAt", Value: time.Now()},
	}}})

	if err != nil {
		return &types.DbOperationResult{OperationSuccess: false}, err
	}

	return &types.DbOperationResult{OperationSuccess: true}, nil
}
