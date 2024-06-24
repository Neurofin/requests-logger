package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateApplicationLogicInput struct {
	Org  primitive.ObjectID
	Flow primitive.ObjectID
}
