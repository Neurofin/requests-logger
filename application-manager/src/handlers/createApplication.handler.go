package handlers

import (
	"application-manager/src/orchestrators"
	authStore "application-manager/src/services/auth/store/types"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateApplication(c echo.Context) error {

	responseData := types.ResponseBody{}

	user := c.Get("user").(authStore.TokenValidationResponseData)

	input := types.CreateApplicationInput{}
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

	newApplication, err := orchestrators.CreateApplication(user.Org, input.FlowId, input.NumberOfDocs)
	if err != nil {
		responseData.Message = "Error creating application"
		responseData.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.Message = "Created application successfully"
	responseData.Data = newApplication
	return c.JSON(http.StatusOK, responseData)

}
