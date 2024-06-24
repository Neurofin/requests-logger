package handlers

import (
	"net/http"

	"application-manager/src/orchestrators"
	"application-manager/src/store/types"

	"github.com/labstack/echo/v4"
)

func AddApplicationDocument(c echo.Context) error {

	responseData := types.ResponseBody{}

	input := types.AddApplicationDocumentInput{}
	if err := c.Bind(&input); err != nil {
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}
	// TODO: Add Validation for Input and uncomment below code
	// isValid, err := input.Validate()
	// if !isValid {
	// 	responseData.Message = "Error parsing json, please check type of each parameter"
	// 	responseData.Data = err.Error()
	// 	return c.JSON(http.StatusBadRequest, responseData)
	// }

	data, err := orchestrators.AddApplicationDocument(input)
	if err != nil {
		responseData.Message = "Error adding document"
		responseData.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseData)
	}

	responseData.Message = "Created document and presigned URL successfully"
	responseData.Data = data
	return c.JSON(http.StatusOK, responseData)
}
