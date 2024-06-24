package handlers

import (
	"application-manager/src/models"
	"application-manager/src/store/types"
	"application-manager/src/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Org            string               `json:"org"`
	Uid            string               `json:"uid"`
	Name           string               `json:"name"`
	DocumentTypes  []DocumentTypeInput  `json:"documentTypes"`
	ChecklistItems []ChecklistItemInput `json:"checklistItems"`
}

func CreateFlow(c echo.Context) error {

	responseData := types.ResponseBody{}
	jsonInput := CreateFlowInput{}
	if err := c.Bind(&jsonInput); err != nil {
		println(err.Error())
		responseData.Message = "Error parsing json, please check type of each parameter"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	org, err := primitive.ObjectIDFromHex(jsonInput.Org)
	if err != nil {
		println(err.Error())
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
			Uid:         docType.Uid,
		}
		docTypeDoc.CreatedAt = time.Now()
		docTypeDoc.UpdatedAt = time.Now()
		documentTypesToBeInserted = append(documentTypesToBeInserted, docTypeDoc)
	}

	if _, err := utils.BulkInsertToDb(documentTypesToBeInserted, "documentType"); err != nil {
		println(err.Error())
		responseData.Message = "Error inserting document types"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	var checklistItemsToBeInserted []interface{}
	for _, checklistItem := range jsonInput.ChecklistItems {
		checklistItemDoc := models.ChecklistItemModel{
			Name:         checklistItem.Name,
			Goal:         checklistItem.Goal,
			Rules:        checklistItem.Rules,
			Taxonomy:     checklistItem.Taxonomy,
			Prompt:       checklistItem.Prompt,
			GroupUid:     checklistItem.GroupUid,
			RequiredDocs: checklistItem.RequiredDocs, //TODO: Add validation with doctype
			Flow:         flowDoc.Id,
		}
		checklistItemDoc.CreatedAt = time.Now()
		checklistItemDoc.UpdatedAt = time.Now()
		checklistItemsToBeInserted = append(checklistItemsToBeInserted, checklistItemDoc)
	}
	if _, err := utils.BulkInsertToDb(checklistItemsToBeInserted, "checklistItem"); err != nil {
		println(err.Error())
		responseData.Message = "Error inserting checklist items"
		responseData.Data = err.Error()
		return c.JSON(http.StatusBadRequest, responseData)
	}

	responseData.Message = "Created Flow Successfully"
	responseData.Data = flowDoc
	return c.JSON(http.StatusCreated, responseData)
}
