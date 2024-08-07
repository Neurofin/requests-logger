package applicationHandlers

import (
	"application-manager/src/orchestrators"
	authStore "application-manager/src/services/auth/store/types"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func OverrideChecklistResult(c echo.Context) error {

	user := c.Get("user").(authStore.TokenValidationResponseData)

	responseData := types.ResponseBody{}
	traceId := c.Get("traceId").(string)

	input := types.OverrideChecklistResultInput{}
	if err := c.Bind(&input); err != nil {
		responseData.TraceId = traceId
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}
	// Added Validation
	isValid, err := input.Validate()
	if !isValid {
		responseData.TraceId = traceId
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	if err := orchestrators.OverrideChecklistResult(input, user); err != nil {
		responseData.TraceId = traceId
		responseData.Message = "Error overriding checklist result"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.TraceId = traceId
	responseData.Message = "Override Successful"
	responseData.Data = input
	return c.JSON(http.StatusOK, responseData)
}
