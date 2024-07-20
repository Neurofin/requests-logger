package orchestrators

import (
	"application-manager/src/logics"
	authStore "application-manager/src/services/auth/store/types"
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func OverrideChecklistResult(input types.OverrideChecklistResultInput, user authStore.TokenValidationResponseData) error {

	application, err := primitive.ObjectIDFromHex(input.AppId)
	if err != nil {
		return err
	}

	checklistResult, err := primitive.ObjectIDFromHex(input.ChecklistResultId)
	if err != nil {
		return err
	}

	// Checklist Result Updation
	output, err := logics.OverrideChecklistResultStatus(checklistResult, input, user)
	if err != nil {
		return err
	}

	passedChecklistItems, err := logics.GetPassedChecklistResults(application)
	if err != nil {
		return err
	}

	applicationDoc, err := logics.GetApplication(application.Hex(), user.Org)
	if err != nil {
		return err
	}

	applicationDoc.PassedChecklistItems = passedChecklistItems
	if _, err := applicationDoc.UpdateApplication(); err != nil {
		return err
	}

	// Log insertion
	overrideMeta := map[string]interface{}{}

	overrideMeta["oldStatus"] = output["oldStatus"]
	overrideMeta["newStatus"] = input.Status

	if err := logics.InsertOverrideChecklistLog(overrideMeta, application, checklistResult, user.Id, input.Note); err != nil {
		return err
	}

	return nil
}
