package handlers

import (
	"application-manager/src/logics"
	"application-manager/src/models"
	fileService "application-manager/src/services/file"
	fileStore "application-manager/src/services/file/store/types"
	"application-manager/src/store/types"
	"application-manager/src/utils"
)

func S3ObjectCreatedEventHandler(eventBody  types.S3EventBody) {

	eventRecords := eventBody.Records

	for _, eventRecord := range eventRecords {
		s3 := eventRecord.S3
		bucket := s3.Bucket
		bucketName := bucket.Name

		object := s3.Object
		objectKey := object.Key

		sourceS3Path := "s3://"+bucketName+"/"+objectKey
		outputS3Path := "s3://"+"textractor-dump"+"/"+objectKey

		document := models.DocumentModel{
			S3Location: sourceS3Path,
		}

		documentFetchResult, err := document.GetDocument()
		if err != nil {
			println("models.Document.GetDocument", err.Error())
			//TODO: Log error
			continue
		}

		document = documentFetchResult.Data.(models.DocumentModel)

		// Extract text from file
		text, err := logics.ExtractTextLogic(sourceS3Path, outputS3Path)
		if err != nil {
			println("logics.ExtractTextLogic", err.Error())
			//TODO: Log error
			continue
		}

		//TODO: Add extract to s3
		textractBucketName := "extracted-application-docs"
		textractOutputPath := "s3://"+textractBucketName+"/"+objectKey

		presignInput := fileStore.GetPresignedUrlInput{
			Bucket: textractBucketName,
			Key:    objectKey,
		}
		presignUrl, err := fileService.GetPresignedUploadUrl(presignInput)
		if err != nil {
			println("fileService.GetPresignedUploadUrl", err.Error())
			//TODO: Log error
			continue
		}

		utils.PutObjectToS3(presignUrl.URL, text)

		// Update db with extracted file location
		document.TextractLocation = textractOutputPath
		document.Status = "TEXTRACTED"
		_, err = document.UpdateDocument()
		if err != nil {
			println("models.Document.UpdateDocument", err.Error())
			//TODO: Log error
			continue
		}

		// Classify the file and update database with classifier output
		classificationOutput, err := logics.ClassifyDoc(text)
		if err != nil {
			println("logics.ClassifyDoc", err.Error())
			//TODO: Log error
			continue
		}
						
		document.ClassifierOutput = classificationOutput
		document.Status = "CLASSIFIED"
		_, err = document.UpdateDocument()
		if err != nil {
			println("models.Document.UpdateDocument", err.Error())
			//TODO: Log error
			continue
		}

		println("S3 Object Created Event Handled")
	}

	// TODO: Run checklist based on uploaded files
}