package orchestrators

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/logics"
	"application-manager/src/models"
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddApplicationDocument(input types.AddApplicationDocumentInput) (map[string]interface{}, error) {

	output := make(map[string]interface{}) // TODO: Create type for this output

	bucket := store.ApplicationDocumentsBucket

	docId := primitive.NewObjectID()
	documentKey := input.ApplicationId + "/" + docId.Hex()

	presignUrl, err := logics.GetPresignedUploadUrl(bucket, documentKey, input.Format)
	if err != nil {
		return output, err
	}

	appId, err := primitive.ObjectIDFromHex(input.ApplicationId)
	if err != nil {
		return output, err
	}

	insertAppDocInput := types.InsertApplicationDocumentInput{
		DocId:       docId,
		Application: appId,
		Name:        input.Name,
		Format:      input.Format,
		Status:      "PENDING", //TODO: Create enum
		S3Location:  "s3://" + bucket + "/" + documentKey,
	}

	applicationDoc, err := logics.InsertApplicationDocument(insertAppDocInput)
	if err != nil {
		return output, err
	}

	// Listener for updates to the documents
	go ApplicationDocumentClassificationEventListener(appId)

	output["document"] = applicationDoc
	output["presignUrl"] = presignUrl
	return output, nil
}

func ApplicationDocumentClassificationEventListener(application primitive.ObjectID) error {

	timeout := 5 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ApplicationDocumentCollection)

	matchPipeline := bson.D{{
		Key: "$match", Value: bson.D{
			{Key: "operationType", Value: "insert"},
			{Key: "fullDocument.status", Value: "CLASSIFIED"},
			{Key: "fullDocument.application", Value: application},
		},
	}}
	changeStream, err := collection.Watch(ctx, mongo.Pipeline{matchPipeline})
	if err != nil {
		return err
	}
	defer changeStream.Close(ctx)

	fmt.Println("Listening for insert events...")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("No events for 5 minutes.")
			document := models.ApplicationDocumentModel{
				Application: application,
			}
			operationResult, err := document.GetDocsReadyToProcess([]string{})
			if err != nil {
				return err
			}

			//Check documents status
			documents := operationResult.Data.([]models.ApplicationDocumentModel)
			for _, doc := range documents {
				if doc.Status != "CLASSIFIED" {
					fmt.Println("All docs are not classified, error in pipeline ", application)
					application := models.ApplicationModel{
						Id: application,
					}
					operationResult, err := application.GetApplication()
					if err != nil {
						return err
					}
					application = operationResult.Data.(models.ApplicationModel)
					application.Status = "ERROR"
					application.UpdateApplication()
					return nil
				}
			}

			//If all docs are classified, start checklist process
			applicationDoc := models.ApplicationModel{
				Id: application,
			}
			applicationDocResult, err := applicationDoc.GetApplication()
			if err != nil {
				println("application.GetApplication", err.Error())
				return err
			}

			applicationDoc = applicationDocResult.Data.(models.ApplicationModel)

			flowOperationResult, err := dbHelpers.GetFlowChecklist(applicationDoc.Flow)
			if err != nil {
				return err
			}

			checklist := flowOperationResult.Data.([]models.ChecklistItemModel)

			var wg sync.WaitGroup
			for _, item := range checklist {
				wg.Add(1)
				go func() {
					defer wg.Done()
					ChecklistItemProcessOrchestrator(item, applicationDoc)
				}()
			}

			wg.Wait()

			//Once all checklists are done, set appropriate application status
			// Update uploaded Doc types
			uniqueDocTypesMap := make(map[string]bool)
			for _, doc := range documents {
				uniqueDocTypesMap[doc.Type] = true
			}

			uniqueDocTypes := []string{}
			for docType := range uniqueDocTypesMap {
				uniqueDocTypes = append(uniqueDocTypes, docType)
			}
			applicationDoc.UploadedDocTypes = uniqueDocTypes

			// Update Succesful checklist items
			operationResult, err = dbHelpers.GetAppChecklistResults(application)
			if err != nil {
				return err
			}

			checklistResults := operationResult.Data.([]models.ChecklistItemResultModel)

			passedChecklistItems := applicationDoc.PassedChecklistItems
			for _, result := range checklistResults {
				if result.Result["status"] == "Success" {
					passedChecklistItems = append(passedChecklistItems, result.Result)
				}
			}
			applicationDoc.PassedChecklistItems = passedChecklistItems

			applicationDoc.Status = "PROCESSED"
			applicationDoc.UpdateApplication()
			return nil
		default:
			if changeStream.Next(ctx) {
				var event bson.M
				if err := changeStream.Decode(&event); err != nil {
					return err
				}
				fmt.Printf("Received insert event: %v\n", event)
				cancel()
				ctx, cancel = context.WithTimeout(context.Background(), timeout)
				defer cancel()
			}
		}
	}
}
