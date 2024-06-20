package middleware

import (
	authService "application-manager/src/services/auth"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ValidateToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jsonData := types.ValidateTokenInput{}
		jsonData.Token = c.Request().Header.Get("authorization")

		token, err := jsonData.Validate()
		if err != nil {
			println(err.Error())
			responseBody := types.ResponseBody{
				Message: "Invalid Token",
				Data:    err.Error(),
			}
			return c.JSON(http.StatusUnauthorized, responseBody)
		}

		userDetails, err := authService.ValidateToken(token)
		if err != nil {
			println(err.Error())
			responseBody := types.ResponseBody{
				Message: "Invalid Token",
				Data:    err.Error(),
			}
			return c.JSON(http.StatusUnauthorized, responseBody)
		}

		c.Set("user", userDetails)
		return next(c)
	}
}
