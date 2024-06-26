package orchestrators

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ApplicationDocChecklistOrchestrator(docType string, application primitive.ObjectID) {
	applicationDoc := models.ApplicationModel{
		Id: application,
	}
	applicationDocResult, err := applicationDoc.GetApplication()
	if err != nil {
		println("application.GetApplication", err.Error())
		return
	}

	applicationDoc = applicationDocResult.Data.(models.ApplicationModel)

	checklistResults, err := dbHelpers.GetFlowChecklist(applicationDoc.Flow)
	if err != nil {
		println("checklist.FetchChecklistItemsForDocType", err.Error())
		return
	}

	flowChecklist := checklistResults.Data.([]models.ChecklistItemModel)

	// Loop through checklists
	for _, checklistItem := range flowChecklist {
		ChecklistItemProcessOrchestrator(checklistItem, applicationDoc)
	}
}
