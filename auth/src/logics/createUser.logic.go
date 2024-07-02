package logics

import (
	"auth/src/models"
	"auth/src/utils"
	"errors"
)

// TODO: Make input specific to logic requirement
func CreateUser(user models.UserModel) (models.UserModel, error) {

	encryptedPassword, encryptionError := utils.EncrytPassword(user.Password)

	if encryptionError != nil {
		user.Password = ""
		return user, encryptionError
	}

	newUser := models.UserModel{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone,
		Password:  encryptedPassword,
		Type:      user.Type,
		Org:       user.Org,
	}

	existingUserResult, _ := newUser.GetUser()
	operationStatus := existingUserResult.OperationSuccess
	if operationStatus {
		return newUser, errors.New("email/phone already exists")
	}

	if _, err := newUser.InsertUser(); err != nil {
		return newUser, err
	}

	getUserResult, err := newUser.GetUser()
	if err != nil {
		return newUser, err
	}

	userDoc := getUserResult.Data.(models.UserModel)
	userDoc.Password = ""
	return userDoc, nil
}
