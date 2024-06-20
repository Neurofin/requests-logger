package handlers

import (
	"application-manager/src/models"
	authServiceTypes "application-manager/src/services/auth/store/types"
	"application-manager/src/store/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DocumentTypeInput struct {
	Uid         string `json:"uid"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ChecklistItemInput struct {
	Name         string   `json:"name"`
	Goal         string   `json:"goal"`
	Rules        []string `json:"rules"`
	Taxonomy     []string `json:"taxonomy"`
	Prompt       string   `json:"prompt"`
	GroupUid     string   `json:"groupUid,omitempty"`
	RequiredDocs []string `json:"requiredDocs"`
}

type CreateFlowInput struct {
	Uid            string               `json:"uid"`
	Name           string               `json:"name"`
	DocumentTypes  []DocumentTypeInput  `json:"documentTypes"`
	ChecklistItems []ChecklistItemInput `json:"checklistItems"`
}

func CreateFlow(c echo.Context) error {
	user := c.Get("user").(authServiceTypes.TokenValidationResponseData)

	responseData := types.ResponseBody{}
	jsonInput := CreateFlowInput{}
	if err := c.Bind(&jsonInput); err != nil {
		println(err.Error())
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	newFlow := models.FlowModel{
		Org:  user.Org,
		Name: jsonInput.Name,
		Uid:  jsonInput.Uid,
	}

	if _, err := newFlow.InsertFlow(); err != nil {
		println(err.Error())
		responseData.Message = "Error inserting flow document"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	flowDocResult, err := newFlow.GetFlow()
	if err != nil {
		println(err.Error())
		responseData.Message = "Error finding flow document"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	flowDoc := flowDocResult.Data.(models.FlowModel)

	var documentTypesToBeInserted []interface{}
	for _, docType := range jsonInput.DocumentTypes {
		docTypeDoc := models.DocumentTypeModel{
			Name:        docType.Name,
			Description: docType.Description,
			Flow:        flowDoc.Id,
		}
		documentTypesToBeInserted = append(documentTypesToBeInserted, docTypeDoc)
	}

	if _, err := models.BulkInsertDocumentTypes(documentTypesToBeInserted); err != nil {
		println(err.Error())
		responseData.Message = "Error inserting document types"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.Message = "Created Flow Successfully"
	responseData.Data = newFlow
	return c.JSON(http.StatusCreated, responseData)
}
