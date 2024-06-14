package handlers

import (
	"application-manager/src/models"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateApplication(c echo.Context) error {
	jsonBody := types.CreateApplicationInput{}

	if err := c.Bind(&jsonBody); err != nil {
		responseData := types.ResponseBody{
			Message: "Error parsing json, please check type of each parameter",
			Data:    err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	isValid, inputErr := jsonBody.Validate()
	if !isValid {
		responseData := types.ResponseBody{
			Message: "Error validating input",
			Data:    inputErr.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	newApplication := models.ApplicationModel{
		Id:   primitive.NewObjectID(),
		Name: jsonBody.Name,
		User: jsonBody.User,
	}

	output, err := newApplication.InsertApplication()
	if err != nil {
		responseData := types.ResponseBody{
			Message: "Error inserting document into database",
			Data:    err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData := types.ResponseBody{
		Message: "Created application successfully",
		Data:    output,
	}
	return c.JSON(http.StatusOK, responseData)

}
