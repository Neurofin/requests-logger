package orchestrators

import (
	"auth/src/logics"
	"auth/src/models"
	"auth/src/store/types"
	"auth/src/utils"
)

func Login(loginInput types.LoginInput) (types.SignedJwtToken, error) {

	token := types.SignedJwtToken{}

	inputUser := models.UserModel{
		Email:    loginInput.Email,
		Phone:    loginInput.Phone,
		Password: loginInput.Password,
	}

	user, err := logics.GetUser(inputUser)
	if err != nil {
		return token, err
	}

	if err := utils.VerifyPassword(inputUser.Password, user.Password); err != nil {
		return token, err
	}

	token, err = logics.GenerateUserJwt(user.Id.Hex(), user.Org.Hex())
	if err != nil {
		return token, err
	}

	return token, nil
}
