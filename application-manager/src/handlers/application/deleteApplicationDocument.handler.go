package applicationHandlers

import (
	"application-manager/src/orchestrators"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func DeleteApplicationDocument(c echo.Context) error {
	responseData := types.ResponseBody{}

	appId := c.Param("appId")
	docId := c.Param("docId")

	if err := orchestrators.DeleteApplicationDocument(appId, docId); err != nil {
		responseData.Message = "Error deleting document"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.Message = "Document is deleted successfully"
	responseData.Data = "Success"
	return c.JSON(http.StatusOK, responseData)
}
