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

type ApplicationModel struct {
	Id                 primitive.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Org                primitive.ObjectID     `json:"org,omitempty" bson:"org,omitempty"`
	Flow               primitive.ObjectID     `json:"flow,omitempty" bson:"flow,omitempty"`
	ApplicationDetails map[string]interface{} `json:"applicationDetails,omitempty" bson:"applicationDetails,omitempty"`
	types.Timestamps
}

func (app *ApplicationModel) InsertApplication() (types.DbOperationResult, error) {

	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ApplicationCollection)

	app.CreatedAt = time.Now()
	app.UpdatedAt = time.Now()
	operationResult, err := collection.InsertOne(context.Background(), app)
	if err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = operationResult.InsertedID
	return result, nil
}

func (app *ApplicationModel) GetApplication() (types.DbOperationResult, error) {
	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ApplicationCollection)

	filter := bson.D{{
		Key:   "_id",
		Value: app.Id,
	}}

	if app.Org != primitive.NilObjectID {
		filter = append(filter, bson.E{
			Key:   "org",
			Value: app.Org,
		})
	}

	data := ApplicationModel{}
	if err := collection.FindOne(context.Background(), filter).Decode(&data); err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = data
	return result, nil
}

func (app *ApplicationModel) UpdateApplication() (types.DbOperationResult, error) {

	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ApplicationCollection)

	UpdateResult, err := collection.UpdateByID(context.Background(), app.Id, bson.D{{Key: "$set", Value: app}})

	if err != nil {
		return result, err
	}

	result.OperationSuccess = true
	result.Data = UpdateResult
	return result, nil
}
