package handlers

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"
	"application-manager/src/orchestrators"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RestartApplicationListners() {
	operationResult, err := dbHelpers.GetPendingApplications()
	if err != nil {
		fmt.Println(err)
		return
	}

	applications, ok := operationResult.Data.([]models.ApplicationModel)
	if !ok {
		fmt.Println("Issue fetching applications")
		return
	}

	for _, app := range applications {
		operationResult, err := dbHelpers.GetApplicationDocuments(app.Id)
		if err != nil {
			fmt.Println(err)
			continue
		}

		docsToBeListenedTo := []primitive.ObjectID{}
		//Check documents status
		documents := operationResult.Data.([]models.ApplicationDocumentModel)
		for _, doc := range documents {
			if doc.Status != "CLASSIFIED" {
				docsToBeListenedTo = append(docsToBeListenedTo, doc.Id)
			}
		}

		if len(docsToBeListenedTo) > 0 {
			orchestrators.DocumentClassificationEventListener(app.Id, docsToBeListenedTo, false)
		} else {
			orchestrators.ProcessAllChecklist(app.Id)
		}
	}
}
