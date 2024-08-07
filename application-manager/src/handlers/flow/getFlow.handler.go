package flowHandlers

import (
	"application-manager/src/orchestrators"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetFlow(c echo.Context) error {
	responseData := types.ResponseBody{}
	traceId := c.Get("traceId").(string)

	flowId := c.Param("flowId")

	data, err := orchestrators.GetFlow(flowId)
	if err != nil {
		println(err.Error())
		responseData.TraceId = traceId
		responseData.Message = "Error finding flow"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.TraceId = traceId
	responseData.Message = "Fetched flow successfully"
	responseData.Data = data
	return c.JSON(http.StatusOK, responseData)
}
