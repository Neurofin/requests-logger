package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenValidationResponseData struct {
	Id              primitive.ObjectID              `json:"id,omitempty"`
	Org             primitive.ObjectID              `json:"org"`
	FirstName       string                          `json:"firstName"`
	LastName        string                          `json:"lastName"`
	Email           string                          `json:"email"`
	Phone           string                          `json:"phone,omitempty"`
	Type            string             				`json:"type"`
	Verified        bool                            `json:"verified"`
	Permissions     []string 						`json:"permissions"`
	CognitoUsername string                          `json:"cognitoUsername,omitempty"`
	Timestamps
}

type Timestamps struct {
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}