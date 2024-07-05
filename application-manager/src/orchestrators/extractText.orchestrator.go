package orchestrators

import (
	"application-manager/src/logics"
	"application-manager/src/models"
	fileService "application-manager/src/services/file"
	fileServiceTypes "application-manager/src/services/file/store/types"
	"application-manager/src/utils"
)

type ExtractTextOrchestratorOutput struct {
	Text   string
	S3Path string
}

func ExtractTextAndUpdateDoc(bucketName string, objectKey string, document models.ApplicationDocumentModel) (ExtractTextOrchestratorOutput, error) {

	output := ExtractTextOrchestratorOutput{}

	sourceS3Path := "s3://" + bucketName + "/" + objectKey
	outputS3Path := "s3://" + "textractor-dump" + "/" + objectKey
	// Extract text from file
	text, err := logics.ExtractTextLogic(sourceS3Path, outputS3Path)
	if err != nil {
		return output, err
	}

	//Add extract to s3
	textractBucketName := "extracted-application-docs"
	textractLocation := "s3://" + textractBucketName + "/" + objectKey

	presignInput := fileServiceTypes.GetPresignedUrlInput{
		Bucket: textractBucketName,
		Key:    objectKey,
	}
	presignUrl, err := fileService.GetPresignedUploadUrl(presignInput)
	if err != nil {
		return output, err
	}

	utils.PutObjectToS3(presignUrl.URL, text)

	// Update db with extracted file location
	document.TextractLocation = textractLocation
	document.Status = "TEXTRACTED"
	_, err = document.UpdateDocument()
	if err != nil {
		return output, err
	}

	output.Text = text
	output.S3Path = textractLocation
	return output, nil
}
