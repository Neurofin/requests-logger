package applicationHandlers

import (
	"application-manager/src/orchestrators"
	authTypes "application-manager/src/services/auth/store/types"
	types "application-manager/src/store/types"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateGetApplicationInput(appId string) (bool, error) {
	if strings.TrimSpace(appId) == "" {
		return false, errors.New("appId is missing or is not a string")
	}

	if _, err := primitive.ObjectIDFromHex(appId); err != nil {
		return false, errors.New("appId is not in valid ObjectID format")
	}

	return true, nil
}

func GetApplication(c echo.Context) error {

	responseData := types.ResponseBody{}
	traceId := c.Get("traceId").(string)

	id := c.Param("id")

	if valid, err := validateGetApplicationInput(id); !valid {
		responseData.TraceId = traceId
		responseData.Message = "Error fetching application"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	user, ok := c.Get("user").(authTypes.TokenValidationResponseData)
	if !ok {
		responseData.TraceId = traceId
		responseData.Message = "Unauthorized"
		responseData.Data = "User not found"
		return c.JSON(http.StatusUnauthorized, responseData)
	}

	result, err := orchestrators.GetApplication(id, user.Org)
	if err != nil {
		responseData.TraceId = traceId
		responseData.Message = "Error fetching application"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.TraceId = traceId
	responseData.Message = "Application details retrieved successfully"
	responseData.Data = result
	return c.JSON(http.StatusOK, responseData)
}
