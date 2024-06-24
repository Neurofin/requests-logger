package logics

import (
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetApplicationDocumet(appId string, docId string) (models.ApplicationDocumentModel, error) {

	document := models.ApplicationDocumentModel{}

	applicationId, err := primitive.ObjectIDFromHex(appId)
	if err != nil {
		return document, err
	}

	documentId, err := primitive.ObjectIDFromHex(docId)
	if err != nil {
		return document, err
	}

	document.Application = applicationId
	document.Id = documentId

	operationResult, err := document.GetApplicationDocumentById()
	if err != nil {
		return document, err
	}

	document = operationResult.Data.(models.ApplicationDocumentModel)
	return document, nil
}
