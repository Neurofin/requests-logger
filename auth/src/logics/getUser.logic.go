package logics

import "auth/src/models"

// TODO: Make input specific to logic requirement
func GetUser(user models.UserModel) (models.UserModel, error) {

	requiredUser := models.UserModel{
		Email:    user.Email,
		Phone:    user.Phone,
		Password: user.Password,
	}

	userFetchResult, err := requiredUser.GetUser()
	if err != nil {
		return requiredUser, err
	}

	requiredUser = userFetchResult.Data.(models.UserModel)
	return requiredUser, nil
}
