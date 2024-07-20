package logics

import (
	"application-manager/src/models"
	authStore "application-manager/src/services/auth/store/types"
	"application-manager/src/store/types"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OverrideChecklistResultStatus(checklistResult primitive.ObjectID, input types.OverrideChecklistResultInput, user authStore.TokenValidationResponseData) (map[string]interface{}, error) {

	output := map[string]interface{}{}

	checklistResultObj := models.ChecklistItemResultModel{Id: checklistResult}

	operationResult, err := checklistResultObj.FindChecklistItemResultById()
	if err != nil {
		return output, err
	}

	checklistResultObj = operationResult.Data.(models.ChecklistItemResultModel)

	oldStatus := checklistResultObj.Result["status"]
	checklistResultObj.Result["status"] = input.Status

	overrideStatus := checklistResultObj.Overridden
	if overrideStatus {
		overrideStatus = false
	} else {
		overrideStatus = true
	}
	checklistResultObj.Overridden = overrideStatus

	checklistResultObj.OverrideMeta = map[string]interface{}{}
	checklistResultObj.OverrideMeta["note"] = input.Note
	checklistResultObj.OverrideMeta["user"] = strings.TrimSpace(user.FirstName + " " + user.LastName)
	if checklistResultObj.OverrideMeta["user"] == "" {
		checklistResultObj.OverrideMeta["user"] = user.Email
	}
	if _, err := checklistResultObj.UpdateChecklistItemResult(); err != nil {
		return output, err
	}

	output["oldStatus"] = oldStatus
	output["newStatus"] = input.Status
	return output, nil
}
