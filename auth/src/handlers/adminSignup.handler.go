package handlers

import (
	orchestrators "auth/src/orchestrator"
	"auth/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AdminSignup(c echo.Context) error {
	responseData := types.ResponseBody{}

	input := types.AdminSignupInput{}
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

	token, err := orchestrators.AdminSignup(input)
	if err != nil {
		responseData.Message = "Error while signing up"
		responseData.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.Message = "User signed up successfully"
	responseData.Data = token
	return c.JSON(http.StatusCreated, responseData)
}
