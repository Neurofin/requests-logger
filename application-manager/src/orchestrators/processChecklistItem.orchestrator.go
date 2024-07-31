package orchestrators

import (
	"application-manager/src/logics"
	"application-manager/src/models"
	"application-manager/src/serverConfigs"
	querierService "application-manager/src/services/querier"
	querierServiceTypes "application-manager/src/services/querier/store/types"
	checklistResultStatusEnum "application-manager/src/store/enums"
	"fmt"
	"strings"
)

func ProcessChecklistItemOrchestrator(checklistItem models.ChecklistItemModel, application models.ApplicationModel) {
	flow := models.FlowModel{
		Id: application.Flow,
	}
	operationResult, err := flow.GetFlow()
	if err != nil {
		println("Error getting flow", err.Error())
		return
	}
	flow = operationResult.Data.(models.FlowModel)

	docsRequired := checklistItem.RequiredDocs

	isMasterQuery := checklistItem.MasterChecklistItem

	document := models.ApplicationDocumentModel{
		Application: application.Id,
	}

	uploadedDocumentsResults, err := document.GetDocsReadyToProcess(docsRequired)
	if err != nil {
		println("Error finding docs", err.Error())
		return
	}

	uploadedDocuments := uploadedDocumentsResults.Data.([]models.ApplicationDocumentModel)

	documentTypeMapping := make(map[string][]interface{})
	if len(uploadedDocuments) > 0 {
		for _, doc := range uploadedDocuments {
			documentTypeMapping[doc.Type] = append(documentTypeMapping[doc.Type], doc)
		}
	}
	if len(documentTypeMapping) < len(docsRequired) && !checklistItem.MasterChecklistItem {
		rule := ""
		if len(checklistItem.Rules) > 0 {
			rule = checklistItem.Rules[0]
		}
		// Create/Update query result
		queryResult := models.ChecklistItemResultModel{
			Application:   application.Id,
			ChecklistItem: checklistItem.Id,
			Result: map[string]interface{}{
				"goal":   checklistItem.Goal,
				"rule":   rule,
				"status": checklistResultStatusEnum.Failed,
				"reason": "Required files not uploaded",
			},
			RequiredDocs: checklistItem.RequiredDocs,
		}

		if _, err := logics.UpsertChecklistItemResultLogic(queryResult); err != nil {
			println("Error upserting query result", err.Error())
			return
		}
		return
	}

	//Call Query Engine for result and update the result
	contextDocuments := []querierServiceTypes.ContextDocument{}
	for _, doc := range uploadedDocuments {
		docPath := doc.TextractLocation
		if flow.QueryFormat == "FILE" {
			docPath = doc.S3Location
		}
		if checklistItem.QueryDocFormat == "TEXT" && doc.TextractLocation == "" {
			bucketName := serverConfigs.ApplicationDocumentsBucket
			objectKey := doc.Application.Hex() + "/" + doc.Id.Hex()
			output, err := ExtractTextAndUpdateDoc(bucketName, objectKey, document)
			if err != nil {
				fmt.Println("Error ", err.Error())
				return
			}
			docPath = output.S3Path
		}
		contextDoc := querierServiceTypes.ContextDocument{
			DocPath: docPath,
			DocType: doc.Type,
		}
		contextDocuments = append(contextDocuments, contextDoc)
	}

	engine := checklistItem.Engine
	if engine == "" {
		engine = flow.Engine
	}

	docFormat := ""
	if flow.QueryFormat == "FILE" {
		docFormat = "application/pdf"
	}

	if checklistItem.QueryDocFormat == "TEXT" {
		docFormat = ""
	}
	queryResultData, err := querierService.ResolveQuery(querierServiceTypes.ResolveQueryInput{
		ContextDocuments: contextDocuments,
		Prompt:           checklistItem.Prompt,
		Engine:           engine,
		DocFormat:        docFormat,
	})
	if err != nil {
		println("Error resolving query", err.Error())
		return
	}

	if isMasterQuery {
		application.ApplicationDetails = queryResultData
		application.UpdateApplication()

		return
	} else {

		rule := ""
		if len(checklistItem.Rules) > 0 {
			rule = strings.Join(checklistItem.Rules, ", ")
		}
		queryResult := models.ChecklistItemResultModel{
			Application:   application.Id,
			ChecklistItem: checklistItem.Id,
			Result: map[string]interface{}{
				"goal":   checklistItem.Goal,
				"rule":   rule,
				"status": queryResultData["status"],
				"reason": queryResultData["reason"],
			},
			RequiredDocs: checklistItem.RequiredDocs,
		}
		if _, err := logics.UpsertChecklistItemResultLogic(queryResult); err != nil {
			println("Error upserting query result", err.Error())
			return
		}
	}
}
