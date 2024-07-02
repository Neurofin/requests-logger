package handlers

import (
	"auth/src/models"
	orchestrators "auth/src/orchestrator"
	"auth/src/store/enums"
	"auth/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Signup(c echo.Context) error {
	responseData := types.ResponseBody{}

	userDetails := c.Get("user").(models.UserModel)

	input := types.SignupInput{}
	if err := c.Bind(&input); err != nil {
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}
	isValid, err := input.Validate()
	if !isValid {
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	newUser := models.UserModel{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
		Password:  input.Password,
		Type:      enums.Member,
		Org:       userDetails.Org,
	}

	token, err := orchestrators.Signup(newUser)
	if err != nil {
		responseData.Message = "Error signing up"
		responseData.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.Message = "User signed up successfully"
	responseData.Data = token
	return c.JSON(http.StatusCreated, responseData)
}
