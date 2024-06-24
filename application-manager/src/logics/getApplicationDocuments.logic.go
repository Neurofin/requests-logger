package logics

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetApplicationDocuments(app primitive.ObjectID) ([]models.ApplicationDocumentModel, error) {
	documents := []models.ApplicationDocumentModel{}

	fetchResult, err := dbHelpers.GetApplicationDocuments(app)
	if err != nil {
		return documents, err
	}

	documents = fetchResult.Data.([]models.ApplicationDocumentModel)
	return documents, nil
}
