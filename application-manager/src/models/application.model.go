package models

import (
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApplicationModel struct {
	Id      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name" bson:"name"`
	Created time.Time          `json:"created" bson:"created"`
	User    primitive.ObjectID `json:"user" bson:"user"`
	//Flow    primitive.ObjectID `json:"flow,omitempty" bson:"flow,omitempty"`
	Results map[string]interface{} `json:"results,omitempty" bson:"results,omitempty"`
}

func (app *ApplicationModel) InsertApplication() (*types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("application")

	// Check for existing application with the same name
	filter := bson.D{{Key: "name", Value: app.Name}}

	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return &types.DbOperationResult{OperationSuccess: false}, err
	}

	if count > 0 {
		return &types.DbOperationResult{OperationSuccess: false}, errors.New("application with this name already exists")
	}

	app.Created = time.Now()
	_, err = collection.InsertOne(context.Background(), app)
	if err != nil {
		return &types.DbOperationResult{OperationSuccess: false}, err
	}

	return &types.DbOperationResult{OperationSuccess: true}, nil
}
