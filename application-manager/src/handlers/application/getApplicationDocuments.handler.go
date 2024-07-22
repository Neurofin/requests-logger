package applicationHandlers

import (
	"application-manager/src/logics"
	"application-manager/src/store/types"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateGetApplicationDocumentsInput(appId string)(bool, error){
	if strings.TrimSpace(appId) == "" {
		return false, errors.New("appId is missing or is not a string")
	}
	
	if _, err := primitive.ObjectIDFromHex(appId); err!= nil {
		return false, errors.New("appId is not in valid ObjectID format")
	}
	
	return true, nil;
}

func GetApplicationDocuments(c echo.Context) error {

	responseData := types.ResponseBody{}

	id := c.Param("id")

	if valid, err := validateGetApplicationDocumentsInput(id); !valid {
		responseData.Message = "Error fetching documents info"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	appId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		responseData.Message = "Error fetching documents info"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	documents, err := logics.GetApplicationDocuments(appId)
	if err != nil {
		responseData.Message = "Error fetching documents info"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.Message = "Documents info retrieved successfully"
	responseData.Data = documents
	return c.JSON(http.StatusOK, responseData)
}
