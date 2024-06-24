package handlers

import (
	"application-manager/src/orchestrators"
	authTypes "application-manager/src/services/auth/store/types"
	types "application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetApplication(c echo.Context) error {

	responseData := types.ResponseBody{}

	id := c.Param("id")

	user, ok := c.Get("user").(authTypes.TokenValidationResponseData)
	if !ok {
		responseData.Message = "Unauthorized"
		responseData.Data = "User not found"
		return c.JSON(http.StatusUnauthorized, responseData)
	}

	result, err := orchestrators.GetApplication(id, user.Org)
	if err != nil {
		responseData.Message = "Error fetching application"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.Message = "Application details retrieved successfully"
	responseData.Data = result
	return c.JSON(http.StatusOK, responseData)
}
