package handlers

import (
	"application-manager/src/logics"
	"application-manager/src/models"
	"application-manager/src/orchestrators"
	"application-manager/src/store/types"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func S3ObjectCreatedEventHandler(eventBody types.S3EventBody) {

	eventRecords := eventBody.Records

	for _, eventRecord := range eventRecords {
		s3 := eventRecord.S3
		bucket := s3.Bucket
		bucketName := bucket.Name

		object := s3.Object
		objectKey := object.Key

		splitArray := strings.Split(objectKey, "/")
		docId := splitArray[1]

		docObjectId, err := primitive.ObjectIDFromHex(docId)
		if err != nil {
			print("Error in doc object id", err)
			continue
		}
		document := models.ApplicationDocumentModel{
			Id: docObjectId,
		}

		documentFetchResult, err := document.GetApplicationDocumentById()
		if err != nil {
			println("models.Document.GetDocument", err.Error())
			//TODO: Log error
			continue
		}

		document = documentFetchResult.Data.(models.ApplicationDocumentModel)
		document.Status = "UPLOADED"
		_, err = document.UpdateDocument()
		if err != nil {
			println("models.Document.UpdateDocument", err.Error())
			//TODO: Log error
			continue
		}

		applicationDoc := models.ApplicationModel{
			Id: document.Application,
		}
		applicationDocResult, err := applicationDoc.GetApplication()
		if err != nil {
			println("application.GetApplication", err.Error())
			return
		}

		applicationDoc = applicationDocResult.Data.(models.ApplicationModel)

		flow := models.FlowModel{
			Id: applicationDoc.Flow,
		}
		operationResult, err := flow.GetFlow()
		if err != nil {
			println("Error getting flow", err.Error())
			return
		}
		flow = operationResult.Data.(models.FlowModel)

		isFileBased := flow.QueryFormat == "FILE"

		text := ""
		fileForClassification := "s3://" + bucketName + "/" + objectKey
		if !isFileBased {
			output, err := orchestrators.ExtractTextAndUpdateDoc(bucketName, objectKey, document)
			if err != nil {
				fmt.Println("Error ", err)
			}
			text = output.Text
			fileForClassification = output.S3Path
		}

		isLLMBased := flow.Classifier == "LLM_BASED"

		// Classify the file and update database with classifier output
		classificationOutput, err := logics.ClassifyDoc(text, isLLMBased, fileForClassification, flow.ClassifierPrompt, isFileBased)
		if err != nil {
			println("logics.ClassifyDoc", err.Error())
			//TODO: Log error
			continue
		}

		document.ClassifierOutput = classificationOutput
		document.Status = "CLASSIFIED"
		document.Type = classificationOutput.Name
		if classificationOutput.Score < 80 {
			document.Type = "OTHER"
		}
		_, err = document.UpdateDocument()
		if err != nil {
			println("models.Document.UpdateDocument", err.Error())
			//TODO: Log error
			continue
		}

		println("S3 Object Created Event Handled")
	}

}
