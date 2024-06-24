package types

import (
	"errors"
)

type CreateApplicationInput struct {
	FlowId       string `json:"flowId"`
	NumberOfDocs int    `json:"numberOfDocs"`
}

func (i *CreateApplicationInput) Validate() (bool, error) {
	if i.FlowId == "" {
		return false, errors.New("flow is missing or is not a string")
	}
	return true, nil
}
