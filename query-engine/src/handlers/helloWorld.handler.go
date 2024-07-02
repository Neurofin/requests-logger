package handlers

import (
	"net/http"
	"query-engine/src/services/auth/store/types"

	"github.com/labstack/echo/v4"
)

func HelloWorldHandler(c echo.Context) error {
	user := c.Get("user").(types.TokenValidationResponseData)
	return c.String(http.StatusOK, "Hello, " + user.Email)
}
