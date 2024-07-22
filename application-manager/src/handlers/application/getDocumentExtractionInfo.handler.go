package applicationHandlers

import (
	"application-manager/src/orchestrators"
	authTypes "application-manager/src/services/auth/store/types"
	"application-manager/src/store/types"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateGetDocumentExtractionInfoInput(appId string) (bool, error) {
	if strings.TrimSpace(appId) == "" {
		return false, errors.New("appId is missing or is not a string")
	}

	if _, err := primitive.ObjectIDFromHex(appId); err != nil {
		return false, errors.New("appId is not in valid ObjectID format")
	}

	return true, nil
}

func GetDocumentExtractionInfo(c echo.Context) error {

	responseData := types.ResponseBody{}

	id := c.Param("id")

	if valid, err := validateGetDocumentExtractionInfoInput(id); !valid {
		return err
	}

	user, ok := c.Get("user").(authTypes.TokenValidationResponseData)
	if !ok {
		responseData.Message = "Unauthorized"
		responseData.Data = "User not found"
		return c.JSON(http.StatusUnauthorized, responseData)
	}

	result, err := orchestrators.GetDocumentExtractionInfo(id, user.Org)
	if err != nil {
		responseData.Message = "Error fetching document extraction info"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.Message = "Document extraction info retrieved successfully"
	responseData.Data = result
	return c.JSON(http.StatusOK, responseData)
}
