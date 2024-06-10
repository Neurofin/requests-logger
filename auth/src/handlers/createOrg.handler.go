package handlers

import (
	"auth/src/models"
	"auth/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateOrg(c echo.Context) error {
	jsonBody := types.CreateOrgInput{}
	c.Bind(&jsonBody)
	isValid, inputErr := jsonBody.Validate()

	if !isValid {

		println(inputErr.Error())
		responseData := types.ResponseBody{
			Message: "Error parsing json, please check type of each parameter",
			Data: inputErr.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	newOrg := models.OrgModel {
		Name: jsonBody.Name,
	}

	output, err := newOrg.InsertOrg()

	if err != nil {

		println(err.Error())
		responseData := types.ResponseBody{
			Message: "Error inserting doc to database",
			Data: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData := types.ResponseBody{
		Message: "Created Org successfully",
		Data: output,
	}
	return c.JSON(http.StatusOK, responseData)
}
