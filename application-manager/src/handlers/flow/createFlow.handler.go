package flowHandlers

import (
	"application-manager/src/logics"
	"application-manager/src/models"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateFlowInput struct {
	Org              string                     `json:"org"`
	Uid              string                     `json:"uid"`
	Name             string                     `json:"name"`
	Classifier       string                     `json:"classifier"`
	ClassifierPrompt string                     `json:"classifierPrompt,omitempty"`
	Engine           string                     `json:"engine"`
	DocumentTypes    []types.DocumentTypeInput  `json:"documentTypes"`
	ChecklistItems   []types.ChecklistItemInput `json:"checklistItems"`
}

func CreateFlow(c echo.Context) error {

	responseData := types.ResponseBody{}
	traceId := c.Get("traceId").(string)
	jsonInput := CreateFlowInput{}
	if err := c.Bind(&jsonInput); err != nil {
		println(err.Error())
		responseData.TraceId = traceId
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	org, err := primitive.ObjectIDFromHex(jsonInput.Org)
	if err != nil {
		println(err.Error())
		responseData.TraceId = traceId
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}
	newFlow := models.FlowModel{
		Org:  org,
		Name: jsonInput.Name,
		Uid:  jsonInput.Uid,
	}

	if _, err = newFlow.InsertFlow(); err != nil {
		println(err.Error())
		responseData.TraceId = traceId
		responseData.Message = "Error inserting flow document"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	flowDocResult, err := newFlow.GetFlow()
	if err != nil {
		println(err.Error())
		responseData.TraceId = traceId
		responseData.Message = "Error finding flow document"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	flowDoc := flowDocResult.Data.(models.FlowModel)

	if len(jsonInput.DocumentTypes) > 0 {
		if _, err := logics.BulkInsertFlowDocTypes(flowDoc.Id, jsonInput.DocumentTypes); err != nil {
			println(err.Error())
			responseData.TraceId = traceId
			responseData.Message = "Error inserting flow doc types"
			responseData.Data = err.Error()
			return c.JSON(http.StatusBadRequest, responseData)
		}
	}

	if len(jsonInput.ChecklistItems) > 0 {
		if _, err := logics.BulkInsertChecklistItems(flowDoc.Id, jsonInput.ChecklistItems); err != nil {
			println(err.Error())
			responseData.TraceId = traceId
			responseData.Message = "Error inserting checklist"
			responseData.Data = err.Error()
			return c.JSON(http.StatusBadRequest, responseData)
		}
	}
	responseData.TraceId = traceId
	responseData.Message = "Created Flow Successfully"
	responseData.Data = flowDoc
	return c.JSON(http.StatusCreated, responseData)
}
