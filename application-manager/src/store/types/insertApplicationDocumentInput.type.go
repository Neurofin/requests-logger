package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type InsertApplicationDocumentInput struct {
	DocId       primitive.ObjectID
	Application primitive.ObjectID
	Status      string
	S3Location  string
}
