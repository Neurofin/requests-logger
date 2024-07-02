package orchestrators

import (
	"auth/src/logics"
	"auth/src/models"
	"auth/src/store/types"
)

func Signup(user models.UserModel) (types.SignedJwtToken, error) {

	token := types.SignedJwtToken{}

	newUser := models.UserModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  user.Password,
		Type:      user.Type,
		Org:       user.Org,
	}

	newUserDoc, err := logics.CreateUser(newUser)
	if err != nil {
		return token, err
	}

	token, err = logics.GenerateUserJwt(newUserDoc.Id.Hex(), newUserDoc.Org.Hex())
	if err != nil {
		return token, err
	}

	return token, nil
}
