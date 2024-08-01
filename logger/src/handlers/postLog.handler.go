package handlers

import (
	"logger/src/logics"
	loggerTypes "logger/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func PostLogHandler(c echo.Context) error {
	responseBody := loggerTypes.ResponseBody{}

	input := loggerTypes.PostLogInput{}
	if err := c.Bind(&input); err != nil {
		responseBody.Message = "Error parsing json, please check field types"
		responseBody.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseBody)
	}
	if err := input.Validate(); err != nil {
		responseBody.Message = "Error parsing json, please check field types"
		responseBody.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseBody)
	}

	// Insert the input to database
	newLog, err := logics.InsertLogLogic(input)
	if err != nil {
		responseBody.Message = "Error inserting log"
		responseBody.Data = err.Error()
		return c.JSON(http.StatusInternalServerError, responseBody)
	}

	responseBody.Message = "Log created successfully!"
	responseBody.Data = newLog
	return c.JSON(http.StatusCreated, responseBody)
}
