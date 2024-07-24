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

func validateDownloadApplicationDocumentInput(appId string, docId string) (bool, error) {
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

func DownloadApplicationDocument(c echo.Context) error {
	responseData := types.ResponseBody{}

	appId := c.Param("appId")
	docId := c.Param("docId")

	if valid, err := validateDownloadApplicationDocumentInput(appId, docId); !valid {
		responseData.Message = "Error fetching application"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	result, err := orchestrators.DownloadApplicationDocument(appId, docId)
	if err != nil {
		responseData.Message = "Error fetching application"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.Message = "Document presigned url retrieved successfully"
	responseData.Data = result
	return c.JSON(http.StatusOK, responseData)
}
