package handlers

import (
	"application-manager/src/logics"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetApplicationChecklistResults(c echo.Context) error {
	responseData := types.ResponseBody{}

	id := c.Param("id")

	appId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		responseData.Message = "Error fetching application checklist results"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	checklistResults, err := logics.GetAppChecklistResults(appId)
	if err != nil {
		responseData.Message = "Error fetching application checklist results"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.Message = "Application checklist results fetched successfully"
	responseData.Data = checklistResults
	return c.JSON(http.StatusOK, responseData)
}
