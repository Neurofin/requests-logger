package orchestrators

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"
	"application-manager/src/serverConfigs"
	"application-manager/src/store"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func DocumentClassificationEventListener(application primitive.ObjectID, documentIds []primitive.ObjectID, limitProcessToDocs bool) error {

	timeout := 7 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	collection := serverConfigs.MongoDBClient.Database(store.DbName).Collection(store.ApplicationDocumentCollection)

	matchPipeline := bson.D{{
		Key: "$match", Value: bson.D{
			{Key: "operationType", Value: "update"},
			{Key: "documentKey._id", Value: bson.D{{
				Key:   "$in",
				Value: documentIds,
			}}},
			{Key: "updateDescription.updatedFields.status", Value: "CLASSIFIED"},
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

			operationResult, err := dbHelpers.GetDocuments(documentIds)
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

			if limitProcessToDocs {
				docTypes := []string{}

				for _, doc := range documents {
					docTypes = append(docTypes, doc.Type)
				}

				err = ProcessLimitedChecklist(application, docTypes)
				if err != nil {
					return err
				}
			} else {
				err = ProcessAllChecklist(application)
				if err != nil {
					return err
				}
			}
			return nil
		default:
			if changeStream.Next(ctx) {
				var event bson.M
				if err := changeStream.Decode(&event); err != nil {
					return err
				}
				fmt.Printf("Received insert event: %v\n", event)

				operationResult, err := dbHelpers.GetDocuments(documentIds)
				if err != nil {
					return err
				}

				//Check documents status
				documents := operationResult.Data.([]models.ApplicationDocumentModel)
				allDocsCalssified := true
				for _, doc := range documents {
					if doc.Status != "CLASSIFIED" {
						cancel()
						ctx, cancel = context.WithTimeout(context.Background(), timeout)
						defer cancel()
						allDocsCalssified = false
					}
				}
				if allDocsCalssified {
					cancel()
				}
			}
		}
	}
}
