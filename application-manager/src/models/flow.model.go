package models

import (
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FlowModel struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Uid        string             `json:"uid" bson:"uid"`
	Org        primitive.ObjectID `json:"org,omitempty" bson:"org,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Classifier string             `json:"classifier,omitempty" bson:"classifier,omitempty"`
	types.Timestamps
}

func (flow *FlowModel) InsertFlow() (types.DbOperationResult, error) {
	result := types.DbOperationResult{OperationSuccess: false}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("flow")

	flow.CreatedAt = time.Now()
	flow.UpdatedAt = time.Now()
	_, err := collection.InsertOne(context.Background(), flow)
	if err != nil {
		return result, err
	}

	result.OperationSuccess = true
	return result, nil
}

func (flow *FlowModel) GetFlow() (types.DbOperationResult, error) {
	result := types.DbOperationResult{OperationSuccess: false}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("flow")

	flowDoc := FlowModel{}
	filter := bson.D{{
		Key: "$or",
		Value: bson.A{
			bson.D{{
				Key:   "uid",
				Value: flow.Uid,
			}},
			bson.D{{
				Key:   "_id",
				Value: flow.Id,
			}},
		},
	}}

	if flow.Org != primitive.NilObjectID {
		filter = append(filter, bson.E{
			Key:   "org",
			Value: flow.Org,
		})
	}

	if err := collection.FindOne(context.Background(), filter).Decode(&flowDoc); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = flowDoc
	return result, nil
}
