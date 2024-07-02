package applicationHandlers

import (
	"application-manager/src/logics"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetApplicationDocuments(c echo.Context) error {

	responseData := types.ResponseBody{}

	id := c.Param("id")

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
