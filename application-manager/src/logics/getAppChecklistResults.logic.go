package logics

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAppChecklistResults(app primitive.ObjectID) ([]models.ChecklistItemResultModel, error) {

	checklistResults := []models.ChecklistItemResultModel{}

	operationResult, err := dbHelpers.GetAppChecklistResults(app)
	if err != nil {
		return checklistResults, err
	}

	checklistResults = operationResult.Data.([]models.ChecklistItemResultModel)
	return checklistResults, nil
}
