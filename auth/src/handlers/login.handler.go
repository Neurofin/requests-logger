package handlers

import (
	"auth/src/logics"
	"auth/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	jsonBody := types.LoginInput{}
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

	token, err := logics.LoginUserLogic(jsonBody)

	if err != nil {
		println(err.Error())
		responseData := types.ResponseBody{
			Message: "Error logging in",
			Data: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData := types.ResponseBody{
		Message: "User loggedin successfully",
		Data: token,
	}
	return c.JSON(http.StatusOK, responseData)
}
