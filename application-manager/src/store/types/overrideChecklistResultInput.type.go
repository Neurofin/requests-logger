package types

type OverrideChecklistResultInput struct {
	AppId             string `json:"applicationId"`
	ChecklistResultId string `json:"checklistResultId"`
	Status            string `json:"status"`
	Note              string `json:"note"`
}
