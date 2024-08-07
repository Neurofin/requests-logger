package applicationHandlers

import (
	"application-manager/src/orchestrators"
	"application-manager/src/store/types"
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateDeleteApplicationDocumentInput(appId string, docId string) (bool, error) {
	if strings.TrimSpace(appId) == "" {
		return false, errors.New("appId is missing or is not a string")
	}

	if _, err := primitive.ObjectIDFromHex(appId); err != nil {
		return false, errors.New("appId is not in valid ObjectID format")
	}

	if strings.TrimSpace(docId) == "" {
		return false, errors.New("docId is missing or is not a string")
	}

	if _, err := primitive.ObjectIDFromHex(docId); err != nil {
		return false, errors.New("docId is not in valid ObjectID format")
	}

	return true, nil
}

func DeleteApplicationDocument(c echo.Context) error {
	responseData := types.ResponseBody{}

	appId := c.Param("appId")
	docId := c.Param("docId")
	traceId := c.Get("traceId").(string)

	if valid, err := validateDeleteApplicationDocumentInput(appId, docId); !valid {
		responseData.TraceId = traceId
		responseData.Message = "Error deleting document"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	if err := orchestrators.DeleteApplicationDocument(appId, docId); err != nil {
		responseData.TraceId = traceId
		responseData.Message = "Error deleting document"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.TraceId = traceId
	responseData.Message = "Document is deleted successfully"
	responseData.Data = "Success"
	return c.JSON(http.StatusOK, responseData)
}
