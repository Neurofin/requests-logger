package logics

import (
	"auth/src/models"
	"auth/src/utils"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateToken(token string) (models.UserModel, error) {

	user := models.UserModel{}

	jwtClaims, err := utils.VerifyToken(token)

	if err != nil {
		return user, err
	}

	expiry := jwtClaims["exp"].(float64)
	currentTime := float64(time.Now().Unix())

	if currentTime > expiry {
		return user, errors.New("token expired")
	}

	orgId := jwtClaims["orgId"].(string)
	userId := jwtClaims["userId"].(string)

	user.Org, err = primitive.ObjectIDFromHex(orgId)
	if err != nil {
		return user, err
	}

	user.Id, err = primitive.ObjectIDFromHex(userId)
	if err != nil {
		return user, err
	}

	getUserData, err := user.GetUser()
	if err != nil {
		return user, err
	}

	userDetails := getUserData.Data.(models.UserModel)
	return userDetails, nil
}
