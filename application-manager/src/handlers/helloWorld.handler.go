package handlers

import (
	"application-manager/src/services/auth/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HelloWorldHandler(c echo.Context) error {
	user := c.Get("user").(types.TokenValidationResponseData)
	return c.String(http.StatusOK, "Hello, "+user.Email)
}
