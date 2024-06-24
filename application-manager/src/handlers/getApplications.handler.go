package handlers

import (
	"application-manager/src/logics"
	authType "application-manager/src/services/auth/store/types"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetApplications(c echo.Context) error {
	responseBody := types.ResponseBody{}
	user := c.Get("user").(authType.TokenValidationResponseData)

	//access db collection
	data, err := logics.GetOrgApplications(user.Org)
	if err != nil {
		responseBody.Message = "Error retrieving applications"
		responseBody.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseBody)
	}

	responseBody.Message = "Applications retrieved successfully"
	responseBody.Data = data
	return c.JSON(http.StatusOK, responseBody)
}
