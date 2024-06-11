package logics

import (
	"auth/src/models"
	"auth/src/store/types"
)

func SignupLogic(user models.UserModel) (types.SignedJwtToken, error) {

	token := types.SignedJwtToken{}

	newUser := models.UserModel {
		FirstName: user.FirstName,
		LastName: user.LastName,
		Email: user.Email,
		Phone: user.Phone,
		Password: user.Password,
		Type: user.Type,
		Org: user.Org,
	}

	newUserDoc, err := CreateUserLogic(newUser)
	if err != nil {
		return token, err
	}

	token, err = GenerateUserJwt(newUserDoc.Id.Hex(), newUserDoc.Org.Hex())
	if err != nil {
		return token, err
	}

	return token, nil
}
