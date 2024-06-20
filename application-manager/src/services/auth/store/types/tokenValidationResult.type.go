package types

import (
	"application-manager/src/services/auth/store/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TokenValidationResponseData struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Org       primitive.ObjectID `json:"org"`
	FirstName string             `json:"firstName"`
	LastName  string             `json:"lastName"`
	Email     string             `json:"email"`
	Phone     string             `json:"phone"`
	Type      enums.UserTypeEnum `json:"type"`
	Verified  bool               `json:"verified"`
	Password  string             `json:"password"`
}
