package logics

import (
	"auth/src/models"
	"auth/src/store/types"
	"auth/src/utils"
)

func LoginUserLogic(loginInput types.LoginInput) (types.SignedJwtToken, error) {

	token := types.SignedJwtToken{}

	inputUser := models.UserModel{
		Email: loginInput.Email,
		Phone: loginInput.Phone,
		Password: loginInput.Password,
	}

	userFetchResult, findError := inputUser.GetUser()
	if findError != nil {
		return token, findError
	}

	user := userFetchResult.Data.(models.UserModel)

	passwordValidationError := utils.VerifyPassword(inputUser.Password, user.Password)

	if passwordValidationError != nil {
		return token, passwordValidationError
	}

	token, err := GenerateUserJwt(user.Id.Hex(), user.Org.Hex())
	if err != nil {
		return token, err
	}

	return token, nil
}
