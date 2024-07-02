package logics

import (
	"auth/src/store/types"
	"auth/src/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateUserJwt(userId string, orgId string) (types.SignedJwtToken, error) {

	tokenObject := types.SignedJwtToken{}

	userDetails := jwt.MapClaims{
		"userId": userId,
		"orgId":  orgId,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}
	signedJwt, err := utils.SignJwtToken(&userDetails)
	if err != nil {
		return tokenObject, err
	}

	tokenObject.Token = signedJwt
	return tokenObject, nil
}
