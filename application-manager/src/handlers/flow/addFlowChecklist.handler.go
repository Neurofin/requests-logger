package flowHandlers

import (
	"application-manager/src/logics"
	"application-manager/src/models"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type input struct {
	Checklist []types.ChecklistItemInput `json:"checklist"`
}

func AddFlowChecklist(c echo.Context) error {
	responseData := types.ResponseBody{}

	flowId := c.Param("flowId")
	jsonInput := input{}
	if err := c.Bind(&jsonInput); err != nil {
		println(err.Error())
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	flow := models.FlowModel{
		Uid: flowId,
	}

	getResult, err := flow.GetFlow()
	if err != nil {
		println(err.Error())
		responseData.Message = "Error finding flow"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	flow = getResult.Data.(models.FlowModel)

	if _, err := logics.BulkInsertChecklistItems(flow.Id, jsonInput.Checklist); err != nil {
		println(err.Error())
		responseData.Message = "Error inserting checklist"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.Message = "Fetched flow successfully"
	responseData.Data = flow
	return c.JSON(http.StatusOK, responseData)
}
