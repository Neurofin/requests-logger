package models

import (
	"application-manager/src/store/types"
	//"go.mongodb.org/mongo-driver/bson"
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
