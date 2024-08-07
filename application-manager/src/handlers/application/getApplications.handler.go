package applicationHandlers

import (
	"application-manager/src/logics"
	authType "application-manager/src/services/auth/store/types"
	"application-manager/src/store/types"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetApplications(c echo.Context) error {
    responseBody := types.ResponseBody{}
    user := c.Get("user").(authType.TokenValidationResponseData)
    traceId := c.Get("traceId").(string)

    pageStr := c.QueryParam("page")
    pageSizeStr := c.QueryParam("pageSize")

    page, err := strconv.Atoi(pageStr)
    if err != nil || page <= 0 {
        page = 1 // default to page 1
    }

    pageSize, err := strconv.Atoi(pageSizeStr)
    if err != nil || pageSize <= 0 {
        pageSize = 10 // default to 10 items
    }

    data, totalPages, err := logics.GetOrgApplications(user.Org, page, pageSize)
    if err != nil {
        responseBody.TraceId = traceId
        responseBody.Message = "Error retrieving applications"
        responseBody.Data = err.Error()
        return c.JSON(http.StatusInternalServerError, responseBody)
    }

    responseBody.TraceId = traceId
    responseBody.Message = "Applications retrieved successfully"
    responseBody.Data = map[string]interface{}{
        "applications": data,
        "totalPages":   totalPages,
    }
    return c.JSON(http.StatusOK, responseBody)
}
