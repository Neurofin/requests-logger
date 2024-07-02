package models

import (
	"context"
	"time"

	"application-manager/src/serverConfigs"
	classifierServiceTypes "application-manager/src/services/classifier/store/types"
	"application-manager/src/store"
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicationDocumentModel struct {
	Id               primitive.ObjectID               `json:"id,omitempty" bson:"_id,omitempty"`
	Application      primitive.ObjectID               `json:"application,omitempty" bson:"application,omitempty"`
	Name             string                           `json:"name" bson:"name"`
	Format           string                           `json:"format" bson:"format"`
	Type             string                           `json:"type,omitempty" bson:"type,omitempty"`
	Status           string                           `json:"status" bson:"status"` // PENDING, UPLOADED, TEXTRACTED, CLASSIFIED, DELETED
	S3Location       string                           `json:"s3Location" bson:"s3Location"`
	TextractLocation string                           `json:"textractLocation,omitempty" bson:"textractLocation,omitempty"`
	ClassifierOutput classifierServiceTypes.ClassData `json:"classifierOutput,omitempty" bson:"classifierOutput,omitempty"` //TODO: Decouple classifier output type from db
	types.Timestamps
}

func (doc *ApplicationDocumentModel) InsertApplicationDocument() (types.DbOperationResult, error) {

	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ApplicationDocumentCollection)

	doc.CreatedAt = time.Now()
	doc.UpdatedAt = time.Now()
	operationResult, err := collection.InsertOne(context.Background(), doc)
	if err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = operationResult.InsertedID
	return result, nil
}

func (doc *ApplicationDocumentModel) GetApplicationDocumentById() (types.DbOperationResult, error) {

	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ApplicationDocumentCollection)

	filter := bson.D{{
		Key:   "_id",
		Value: doc.Id,
	}}

	document := ApplicationDocumentModel{}
	if err := collection.FindOne(context.Background(), filter).Decode(&document); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = document
	return result, nil
}

func (doc *ApplicationDocumentModel) UpdateDocument() (*types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ApplicationDocumentCollection)

	_, err := collection.UpdateByID(context.Background(), doc.Id, bson.D{{Key: "$set", Value: bson.D{
		{Key: "classifierOutput", Value: doc.ClassifierOutput},
		{Key: "s3Location", Value: doc.S3Location},
		{Key: "status", Value: doc.Status},
		{Key: "type", Value: doc.Type},
		{Key: "textractLocation", Value: doc.TextractLocation},
		{Key: "timestamps.updatedAt", Value: time.Now()},
	}}})

	if err != nil {
		return &types.DbOperationResult{OperationSuccess: false}, err
	}

	return &types.DbOperationResult{OperationSuccess: true}, nil
}

func (doc *ApplicationDocumentModel) GetDocsReadyToProcess(docTypes []string) (types.DbOperationResult, error) {
	result := types.DbOperationResult{}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ApplicationDocumentCollection)

	filter := bson.D{
		{
			Key:   "application",
			Value: doc.Application,
		},
	}
	if len(docTypes) > 0 {
		filter = append(filter, bson.E{Key: "type",
			Value: bson.D{
				{Key: "$in",
					Value: docTypes,
				},
			}})
	}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return result, err
	}

	var data []ApplicationDocumentModel
	if err = cursor.All(context.TODO(), &data); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = data
	return result, nil
}
