package applicationHandlers

import (
	"application-manager/src/orchestrators"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func DownloadApplicationDocument(c echo.Context) error {
	responseData := types.ResponseBody{}

	appId := c.Param("appId")
	docId := c.Param("docId")
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
