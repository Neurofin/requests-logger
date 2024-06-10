package handlers

import (
	"auth/src/models"
	"auth/src/store/enums"
	"auth/src/store/types"
	"auth/src/utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AdminSignup(c echo.Context) error {

	jsonBody := types.AdminSignupInput{}
	c.Bind(&jsonBody)
	isValid, inputErr := jsonBody.Validate()

	if !isValid {

		println(inputErr.Error())
		responseData := types.ResponseBody{
			Message: "Error parsing json, please check type of each parameter",
			Data: inputErr.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	encryptedPassword, encryptionError := utils.EncrytPassword(jsonBody.Password)

	if encryptionError != nil {
		println(encryptionError.Error())
		responseData := types.ResponseBody{
			Message: "Error encryting password",
			Data: encryptionError.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	orgObjectId, err := primitive.ObjectIDFromHex(jsonBody.Org)
	if err != nil {
		println(err.Error())
		responseData := types.ResponseBody{
			Message: "Error signing up",
			Data: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	newUser := models.UserModel {
		FirstName: jsonBody.FirstName,
		LastName: jsonBody.LastName,
		Email: jsonBody.Email,
		Phone: jsonBody.Phone,
		Password: encryptedPassword,
		Type: enums.Member,
		Org: orgObjectId,
	}

	output, err := newUser.InsertUser()

	if err != nil {

		println(err.Error())
		responseData := types.ResponseBody{
			Message: "Error inserting doc to database",
			Data: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData := types.ResponseBody{
		Message: "User signed up successfully",
		Data: output,
	}
	return c.JSON(http.StatusOK, responseData)
}
