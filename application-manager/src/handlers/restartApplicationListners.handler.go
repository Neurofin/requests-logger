package handlers

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"
	"application-manager/src/orchestrators"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RestartApplicationListners(c echo.Context) error {
	operationResult, err := dbHelpers.GetPendingApplications()
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, err)
	}

	applications, ok := operationResult.Data.([]models.ApplicationModel)
	if !ok {
		fmt.Println("Issue fetching applications")
		return c.JSON(http.StatusBadRequest, err)
	}

	const batchSize = 100000

	var wg sync.WaitGroup
	for i := 0; i < len(applications); i += batchSize {
		end := i + batchSize
		if end > len(applications) {
			end = len(applications)
		}

		// Process a batch of applications
		batch := applications[i:end]

		for _, app := range batch {
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
				wg.Add(1)
				go func(appId primitive.ObjectID, docs []primitive.ObjectID) {
					defer wg.Done()
					orchestrators.DocumentClassificationEventListener(appId, docs, false)
				}(app.Id, docsToBeListenedTo)
			} else {
				orchestrators.ProcessAllChecklist(app.Id)
			}
		}

		// Wait for the current batch of goroutines to finish
		wg.Wait()
	}

	output := map[string]interface{}{}
	output["message"] = "Success"
	return c.JSON(http.StatusOK, output)
}
