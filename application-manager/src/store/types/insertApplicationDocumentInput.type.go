package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type InsertApplicationDocumentInput struct {
	DocId       primitive.ObjectID
	Application primitive.ObjectID
	Name        string
	Format      string
	Status      string
	S3Location  string
}
