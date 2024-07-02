package orchestrators

import (
	"application-manager/src/logics"
	"application-manager/src/models"
	querierService "application-manager/src/services/querier"
	querierServiceTypes "application-manager/src/services/querier/store/types"
)

func ChecklistItemProcessOrchestrator(checklistItem models.ChecklistItemModel, application models.ApplicationModel) {

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

	// if isMasterQuery {
	// 	docsRequired = []string{}
	// }

	document := models.ApplicationDocumentModel{
		Application: application.Id,
	}

	uploadedDocumentsResults, err := document.GetDocsReadyToProcess(docsRequired)
	if err != nil {
		println("Error finding docs", err)
		return
	}

	uploadedDocuments := uploadedDocumentsResults.Data.([]models.ApplicationDocumentModel)
	documentTypeMapping := map[string][]interface{}{}
	for _, doc := range uploadedDocuments {
		documentTypeMapping[doc.Type] = append(documentTypeMapping[doc.Type], doc)
	}
	if len(documentTypeMapping) < len(docsRequired) && !checklistItem.MasterChecklistItem {
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

	// if isMasterQuery {
	// 	unclassifiedDocFound := false
	// 	for _, doc := range uploadedDocuments {
	// 		if doc.Status != "CLASSIFIED" {
	// 			unclassifiedDocFound = true
	// 			break
	// 		}
	// 	}
	// 	if unclassifiedDocFound {
	// 		return
	// 	}
	// }

	//Call Query Engine for result and update the result
	contextDocuments := []querierServiceTypes.ContextDocument{}
	for _, doc := range uploadedDocuments {
		docPath := doc.TextractLocation
		if flow.QueryFormat == "FILE" {
			docPath = doc.S3Location
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
		docFormat = "application/json"
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

	// applicationDocResult, err := application.GetApplication()
	// if err != nil {
	// 	println("application.GetApplication", err.Error())
	// 	return
	// }

	// application = applicationDocResult.Data.(models.ApplicationModel)

	if isMasterQuery {
		application.ApplicationDetails = queryResultData
		// application.Status = "PROCESSED"

		// document := models.ApplicationDocumentModel{
		// 	Application: application.Id,
		// }

		// uploadedDocumentsResults, err := document.GetDocsReadyToProcess(docsRequired)
		// if err != nil {
		// 	println("Error finding docs", err)
		// 	return
		// }

		// uploadedDocuments := uploadedDocumentsResults.Data.([]models.ApplicationDocumentModel)

		// for _, doc := range uploadedDocuments {
		// 	if doc.Status != "CLASSIFIED" {
		// 		println("Not classified")
		// 		application.Status = "PENDING"
		// 		break
		// 	}
		// }
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

		// if queryResultData["status"] == "Success" {
		// 	passedChecklistItems := application.PassedChecklistItems
		// 	itemFound := false
		// 	for i := range passedChecklistItems {
		// 		if passedChecklistItems[i]["goal"] == queryResultData["goal"] {
		// 			itemFound = true
		// 		}
		// 	}
		// 	if !itemFound {
		// 		passedChecklistItems = append(passedChecklistItems, queryResultData)
		// 		application.PassedChecklistItems = passedChecklistItems
		// 		application.UpdateApplication()
		// 	}
		// } else {
		// 	passedChecklistItems := application.PassedChecklistItems
		// 	itemFound := false
		// 	indexOfItem := -1

		// 	// Loop through the passed checklist items to find the item with the matching goal
		// 	for i := range passedChecklistItems {
		// 		if goal, ok := passedChecklistItems[i]["goal"]; ok && goal == queryResultData["goal"] {
		// 			itemFound = true
		// 			indexOfItem = i
		// 			break
		// 		}
		// 	}

		// 	// If the item is found, remove it from the list
		// 	if itemFound && indexOfItem >= 0 {
		// 		// Move the last item to the index of the found item
		// 		passedChecklistItems[indexOfItem] = passedChecklistItems[len(passedChecklistItems)-1]
		// 		// Truncate the slice to remove the last item
		// 		passedChecklistItems = passedChecklistItems[:len(passedChecklistItems)-1]
		// 		// Update the application with the modified checklist items
		// 		application.PassedChecklistItems = passedChecklistItems
		// 		application.UpdateApplication()
		// 	}
		// }
	}
}
