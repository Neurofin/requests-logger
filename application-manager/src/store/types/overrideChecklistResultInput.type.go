package types

import (
	checklistResultStatusEnum "application-manager/src/store/enums"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OverrideChecklistResultInput struct {
	AppId             string `json:"applicationId"`
	ChecklistResultId string `json:"checklistResultId"`
	Status            checklistResultStatusEnum.Status `json:"status"`
	Note              string `json:"note"`
}


func (i *OverrideChecklistResultInput) Validate() (bool, error) {
	if i.AppId == "" {
		return false, errors.New("ApplicationId is missing or is not a String")
	}

	if _, err := primitive.ObjectIDFromHex(i.AppId); err != nil {
		return false, errors.New("ApplicationId is not in valid ObjectID format")
	}

	if i.ChecklistResultId == "" {
		return false, errors.New("ChecklistResultId is missing or is not a String")
	}

	if _, err := primitive.ObjectIDFromHex(i.ChecklistResultId); err != nil {
		return false, errors.New("ChecklistResultId is not in valid ObjectID format")
	}

	if i.Status!=checklistResultStatusEnum.Success && i.Status!=checklistResultStatusEnum.Failed {
    	return false, errors.New("status must be either 'Success' or 'Failed'")
	}

	if i.Note == "" {
		return false, errors.New("override Note is missing or is not a String")
	}

	return true, nil
}