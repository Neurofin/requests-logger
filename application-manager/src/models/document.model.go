package models

import (

	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentModel struct {
	Id               primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string                 `json:"name" bson:"name"`
	Application      primitive.ObjectID     `json:"application" bson:"application"`
	Status           string                 `json:"status" bson:"status"`
	S3Location       string                 `json:"s3_location" bson:"s3_location"`
	TextractLocation string                 `json:"textract_location" bson:"textract_location"`
	ClassifierOutput map[string]interface{} `json:"classifier_output" bson:"classifier_output"`
	QueriesOutput    map[string]interface{} `json:"queries_output" bson:"queries_output"`
}
