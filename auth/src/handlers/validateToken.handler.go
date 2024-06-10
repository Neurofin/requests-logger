package handlers

import (
	"auth/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ValidateToken(c echo.Context) error {
	userDetails := c.Get("user")

	responseBody := types.ResponseBody{
		Message: "Token verified",
		Data: userDetails,
	}
	return c.JSON(http.StatusOK, responseBody)
}
