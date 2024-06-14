package types

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateApplicationInput struct {
	Name string             `json:"name,omitempty" query:"name"`
	User primitive.ObjectID `json:"user,omitempty" query:"user"`
}

func (i *CreateApplicationInput) Validate() (bool, error) {
	if i.Name == "" {
		return false, errors.New("name is missing or is not a string")
	}
	if i.User.IsZero() {
		return false, errors.New("user is missing or is not a valid ObjectID")
	}
	return true, nil
}
