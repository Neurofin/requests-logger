package applicationHandlers

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"
	"application-manager/src/orchestrators"
	authStore "application-manager/src/services/auth/store/types"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateApplication(c echo.Context) error {

	responseData := types.ResponseBody{}

	user := c.Get("user").(authStore.TokenValidationResponseData)
	traceId := c.Get("traceId").(string)

	input := types.CreateApplicationInput{}
	if err := c.Bind(&input); err != nil {
		responseData.TraceId = traceId
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}
	isValid, err := input.Validate()
	if !isValid {
		responseData.TraceId = traceId
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	if input.FlowId == "" {
		operationResult, err := dbHelpers.GetOrgFlows(user.Org)
		if err != nil {
			responseData.TraceId = traceId
			responseData.Message = "Error creating application"
			responseData.Data = err.Error()
			return c.JSON(http.StatusInternalServerError, responseData)
		}
		flows := operationResult.Data.([]models.FlowModel)
		flow := flows[0]
		input.FlowId = flow.Uid
	}
	newApplication, err := orchestrators.CreateApplication(user.Org, input.FlowId, input.DocsToBeUploaded)
	if err != nil {
		responseData.TraceId = traceId
		responseData.Message = "Error creating application"
		responseData.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.TraceId = traceId
	responseData.Message = "Created application successfully"
	responseData.Data = newApplication
	return c.JSON(http.StatusOK, responseData)

}
