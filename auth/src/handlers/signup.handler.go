package handlers

import (
	"auth/src/logics"
	"auth/src/models"
	"auth/src/store/enums"
	"auth/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Signup(c echo.Context) error {

	userDetails := c.Get("user").(models.UserModel)

	jsonBody := types.SignupInput{}
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

	newUser := models.UserModel {
		FirstName: jsonBody.FirstName,
		LastName: jsonBody.LastName,
		Email: jsonBody.Email,
		Phone: jsonBody.Phone,
		Password: jsonBody.Password,
		Type: enums.Member,
		Org: userDetails.Org,
	}

	token, err := logics.SignupLogic(newUser)
	if err != nil {

		println(err.Error())
		responseData := types.ResponseBody{
			Message: "Error signing up",
			Data: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData := types.ResponseBody{
		Message: "User signed up successfully",
		Data: token,
	}
	return c.JSON(http.StatusCreated, responseData)
}
