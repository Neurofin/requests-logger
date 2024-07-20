package logics

import (
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OverrideChecklistResultStatus(checklistResult primitive.ObjectID, newStatus string) (map[string]interface{}, error) {

	output := map[string]interface{}{}

	checklistResultObj := models.ChecklistItemResultModel{Id: checklistResult}

	operationResult, err := checklistResultObj.FindChecklistItemResultById()
	if err != nil {
		return output, err
	}

	checklistResultObj = operationResult.Data.(models.ChecklistItemResultModel)

	oldStatus := checklistResultObj.Result["status"]
	checklistResultObj.Result["status"] = newStatus

	overrideStatus := checklistResultObj.Overridden
	if overrideStatus {
		overrideStatus = false
	} else {
		overrideStatus = true
	}
	checklistResultObj.Overridden = overrideStatus
	if _, err := checklistResultObj.UpdateChecklistItemResult(); err != nil {
		return output, err
	}

	output["oldStatus"] = oldStatus
	output["newStatus"] = newStatus
	return output, nil
}
