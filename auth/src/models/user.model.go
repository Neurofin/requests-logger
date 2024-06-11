package models

import (
	"auth/src/serverConfigs"
	"auth/src/store"
	"auth/src/store/enums"
	"auth/src/store/types"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	Id primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Org primitive.ObjectID `json:"org" bson:"org"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName string `json:"lastName" bson:"lastName"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Phone string `json:"phone,omitempty" bson:"phone,omitempty"`
	Type enums.UserTypeEnum `json:"type" bson:"type"`
	Password string `json:"password" bson:"password"`
	Verified bool `json:"verified" bson:"verified"`
	types.Timestamps
}

func (user *UserModel) InsertUser() (*types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("user")

	userType := user.Type

	userTypeIsValid, userTypevalidationError := userType.Validate()

	if !userTypeIsValid || userTypevalidationError != nil {
		result := &types.DbOperationResult{
			OperationSuccess: false,
		}
		return result, errors.New("invalid user type")
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	_, err := collection.InsertOne(context.Background(), user)

	if err != nil {
		result := &types.DbOperationResult{
			OperationSuccess: false,
		}
		return result, err
	}

	result := &types.DbOperationResult{
		OperationSuccess: true,
	}
	return result, err
}

func (user *UserModel) GetUser() (*types.DbOperationResult, error) {
	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection("user")

	var userDoc UserModel

	filter := bson.D{{
		Key: "$or",
		Value: bson.A{
			bson.D{{
				Key: "email",
				Value: user.Email,
			}},
			bson.D{{
				Key: "phone",
				Value: user.Phone,
			}},
			bson.D{{
				Key: "_id",
				Value: user.Id,
			}},
		},
	}}
	err := collection.FindOne(context.Background(), filter).Decode(&userDoc)
	if err !=nil {
		result := &types.DbOperationResult{
			OperationSuccess: false,
		}
		return result, err
	}

	result := &types.DbOperationResult{
		OperationSuccess: true,
		Data: userDoc,
	}
	return result, err
}
