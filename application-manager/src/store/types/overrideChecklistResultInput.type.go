package types

import (
	"errors"
	"regexp"
)

type OverrideChecklistResultInput struct {
	AppId             string `json:"applicationId"`
	ChecklistResultId string `json:"checklistResultId"`
	Status            string `json:"status"`
	Note              string `json:"note"`
}

// Regex pattern for MongoDB ObjectID
var objectIdPattern = regexp.MustCompile("^[a-fA-F0-9]{24}$")

func (i *OverrideChecklistResultInput) Validate() (bool, error) {
	if i.AppId == "" {
		return false, errors.New("ApplicationId is missing or is not a String")
	}

	if !objectIdPattern.MatchString(i.AppId) {
		return false, errors.New("ApplicationId is not in valid ObjectID format")
	}

	if i.ChecklistResultId == "" {
		return false, errors.New("ChecklistResultId is missing or is not a String")
	}

	if !objectIdPattern.MatchString(i.ChecklistResultId) {
		return false, errors.New("ChecklistResultId is not in valid ObjectID format")
	}

	if i.Status != "Success" && i.Status != "Failed" {
    	return false, errors.New("status must be either 'Success' or 'Failed'")
	}

	if i.Note == "" {
		return false, errors.New("override Note is missing or is not a String")
	}

	return true, nil
}