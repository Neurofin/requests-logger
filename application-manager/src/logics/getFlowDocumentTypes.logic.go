package logics

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetFlowDocumentTypes(flow primitive.ObjectID) ([]models.DocumentTypeModel, error) {

	documentTypes := []models.DocumentTypeModel{}

	operationResult, err := dbHelpers.GetFlowDocumentTypes(flow)
	if err != nil {
		return documentTypes, err
	}

	documentTypes = operationResult.Data.([]models.DocumentTypeModel)
	return documentTypes, nil
}
