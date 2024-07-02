package handlers

import (
	orchestrators "auth/src/orchestrator"
	"auth/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	responseData := types.ResponseBody{}

	input := types.LoginInput{}
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

	token, err := orchestrators.Login(input)
	if err != nil {
		responseData.Message = "Error logging in"
		responseData.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.Message = "User loggedin successfully"
	responseData.Data = token
	return c.JSON(http.StatusOK, responseData)
}
