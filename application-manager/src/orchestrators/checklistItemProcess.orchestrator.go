package orchestrators

import (
	"application-manager/src/logics"
	"application-manager/src/models"
	querierService "application-manager/src/services/querier"
	querierServiceTypes "application-manager/src/services/querier/store/types"
	"fmt"
)

func ChecklistItemProcessOrchestrator(checklistItem models.ChecklistItemModel, application models.ApplicationModel) {
	docsRequired := checklistItem.RequiredDocs

	isMasterQuery := checklistItem.MasterChecklistItem

	if isMasterQuery {
		docsRequired = []string{}
	}

	document := models.ApplicationDocumentModel{
		Application: application.Id,
	}

	uploadedDocumentsResults, err := document.GetDocsReadyToProcess(docsRequired)
	if err != nil {
		println("Error finding docs", err)
		return
	}

	uploadedDocuments := uploadedDocumentsResults.Data.([]models.ApplicationDocumentModel)
	if len(uploadedDocuments) < len(docsRequired) && !checklistItem.MasterChecklistItem {
		// Create/Update query result
		queryResult := models.ChecklistItemResultModel{
			Application:   application.Id,
			ChecklistItem: checklistItem.Id,
			Result: map[string]interface{}{
				"goal":   checklistItem.Goal,
				"rule":   checklistItem.Rules[0],
				"status": "Failed",
				"reason": "Required files not uploaded",
			},
		}

		if _, err := logics.UpsertChecklistItemResultLogic(queryResult); err != nil {
			println("Error upserting query result", err)
			return
		}
		return
	}

	if isMasterQuery {
		for _, doc := range uploadedDocuments {
			if doc.Status != "CLASSIFIED" {
				return
			}
		}
	}

	//Call Query Engine for result and update the result
	contextDocuments := []querierServiceTypes.ContextDocument{}
	for _, doc := range uploadedDocuments {
		contextDoc := querierServiceTypes.ContextDocument{
			DocPath: doc.TextractLocation,
			DocType: doc.Type,
		}
		contextDocuments = append(contextDocuments, contextDoc)
	}
	queryResultData, err := querierService.ResolveQuery(querierServiceTypes.ResolveQueryInput{ContextDocuments: contextDocuments, Prompt: checklistItem.Prompt})
	if err != nil {
		println("Error resolving query", err.Error())
		return
	}

	applicationDocResult, err := application.GetApplication()
	if err != nil {
		println("application.GetApplication", err.Error())
		return
	}

	application = applicationDocResult.Data.(models.ApplicationModel)

	if isMasterQuery {
		application.ApplicationDetails = queryResultData
		application.Status = "PROCESSED"

		uniqueDocTypesMap := make(map[string]bool)
		for _, doc := range uploadedDocuments {
			uniqueDocTypesMap[doc.Type] = true
		}

		uniqueDocTypes := []string{}
		for docType := range uniqueDocTypesMap {
			uniqueDocTypes = append(uniqueDocTypes, docType)
		}
		application.UploadedDocTypes = uniqueDocTypes
		application.UpdateApplication()
		return
	} else {

		queryResult := models.ChecklistItemResultModel{
			Application:   application.Id,
			ChecklistItem: checklistItem.Id,
			Result:        queryResultData,
		}
		if _, err := logics.UpsertChecklistItemResultLogic(queryResult); err != nil {
			println("Error upserting query result", err)
			return
		}

		if queryResultData["status"] == "Success" {
			passedChecklistItems := application.PassedChecklistItems
			itemFound := false
			for i := range passedChecklistItems {
				if passedChecklistItems[i]["goal"] == queryResultData["goal"] {
					itemFound = true
				}
			}
			if !itemFound {
				passedChecklistItems = append(passedChecklistItems, queryResultData)
				application.PassedChecklistItems = passedChecklistItems
				application.UpdateApplication()
			}
		} else {
			passedChecklistItems := application.PassedChecklistItems
			itemFound := false
			indexOfItem := -1
			for i := range passedChecklistItems {
				if passedChecklistItems[i]["goal"] == queryResultData["goal"] {
					itemFound = true
					indexOfItem = i
				}
			}
			if itemFound && indexOfItem >= 0 {
				fmt.Println("Item Found", itemFound)
				fmt.Println("Index", indexOfItem)
				fmt.Println("In ", passedChecklistItems)
				passedChecklistItems[indexOfItem] = passedChecklistItems[len(passedChecklistItems)-1]
				passedChecklistItems = passedChecklistItems[:len(passedChecklistItems)-1]
				fmt.Println("Out ", passedChecklistItems)
				application.PassedChecklistItems = passedChecklistItems
				application.UpdateApplication()
			}
		}
	}
}
