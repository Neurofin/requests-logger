package orchestrators

import (
	"application-manager/src/logics"
	"application-manager/src/models"
	querierService "application-manager/src/services/querier"
	querierServiceTypes "application-manager/src/services/querier/store/types"
)

func ChecklistItemProcessOrchestrator(checklistItem models.ChecklistItemModel, application models.ApplicationModel) {
	docsRequired := checklistItem.RequiredDocs

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

	isMasterQuery := checklistItem.MasterChecklistItem

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

	if isMasterQuery {
		application.ApplicationDetails = queryResultData
		application.UpdateApplication()
		return
	}
	queryResult := models.ChecklistItemResultModel{
		Application:   application.Id,
		ChecklistItem: checklistItem.Id,
		Result:        queryResultData,
	}
	if _, err := logics.UpsertChecklistItemResultLogic(queryResult); err != nil {
		println("Error upserting query result", err)
		return
	}
}
