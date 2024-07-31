package logics

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"
	checklistResultStatusEnum "application-manager/src/store/enums"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetPassedChecklistResults(application primitive.ObjectID) ([]map[string]interface{}, error) {

	passedChecklistItems := []map[string]interface{}{}

	operationResult, err := dbHelpers.GetAppChecklistResults(application)
	if err != nil {
		return passedChecklistItems, err
	}

	checklistResults := operationResult.Data.([]models.ChecklistItemResultModel)

	for _, result := range checklistResults {
		if result.Result["status"] == checklistResultStatusEnum.Success {
			passedChecklistItems = append(passedChecklistItems, result.Result)
		}
	}

	return passedChecklistItems, nil
}
