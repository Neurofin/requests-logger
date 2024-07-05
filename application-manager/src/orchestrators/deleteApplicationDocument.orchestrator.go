package orchestrators

import (
	"application-manager/src/models"
	fileService "application-manager/src/services/file"
	fileServiceTypes "application-manager/src/services/file/store/types"
	"net/url"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DeleteApplicationDocument(appId string, docId string) error {

	docObjectId, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		return err
	}

	appObjectId, err := primitive.ObjectIDFromHex(appId)
	if err != nil {
		return err
	}

	document := models.ApplicationDocumentModel{
		Id: docObjectId,
	}

	operationResult, err := document.GetApplicationDocumentById()
	if err != nil {
		return err
	}

	document = operationResult.Data.(models.ApplicationDocumentModel)

	//TODO: Sperate logic
	s3Location := document.S3Location
	if s3Location != "" {
		s3Path := s3Location
		parsed, err := url.Parse(s3Path)
		if err != nil {
			return err
		}

		bucket := parsed.Host
		key := parsed.Path[1:]

		fileService.DeleteFile(fileServiceTypes.DeleteFileInput{
			Bucket: bucket,
			Key:    key,
		})
	}
	//TODO: Sperate logic
	textractLocation := document.TextractLocation
	if textractLocation != "" {
		s3Path := textractLocation
		parsed, err := url.Parse(s3Path)
		if err != nil {
			return err
		}

		bucket := parsed.Host
		key := parsed.Path[1:]

		fileService.DeleteFile(fileServiceTypes.DeleteFileInput{
			Bucket: bucket,
			Key:    key,
		})
	}

	docType := document.Type

	_, err = document.DeleteDocumentById()
	if err != nil {
		return err
	}

	go ProcessLimitedChecklist(appObjectId, []string{docType})

	return nil
}
