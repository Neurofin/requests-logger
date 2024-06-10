package handlers

import (
	"auth/src/models"
	"auth/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetOrg(c echo.Context) error {
	jsonInput := types.CreateOrgInput{}
	c.Bind(&jsonInput)
	isValid, inputErr := jsonInput.Validate()

	if !isValid {

		println(inputErr.Error())
		responseData := types.ResponseBody{
			Message: "Error parsing json, please check type of each parameter",
			Data:    inputErr.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	org := models.OrgModel{
		Name: jsonInput.Name,
	}

	output, err := org.GetOrg()

	if err != nil {

		println(err.Error())
		responseData := types.ResponseBody{
			Message: "Error inserting doc to database",
			Data:    err.Error(),
		}
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData := types.ResponseBody{
		Message: "Created Org successfully",
		Data:    output.Data,
	}
	return c.JSON(http.StatusOK, responseData)
}
