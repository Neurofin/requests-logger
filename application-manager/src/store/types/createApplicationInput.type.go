package types

import (
	"errors"
	"strings"
)

type DocsToBeUploaded struct {
	Name   string `json:"name"`
	Format string `json:"format"`
}

type CreateApplicationInput struct {
	FlowId           string             `json:"flowId"`
	DocsToBeUploaded []DocsToBeUploaded `json:"docsToBeUploaded"`
}

func (i *CreateApplicationInput) Validate() (bool, error) {
	trimmedFlowId := strings.TrimSpace(i.FlowId)
	if trimmedFlowId == "" {
		return false, errors.New("flow is missing or is not a string")
	}
	for _, doc := range i.DocsToBeUploaded {
		if strings.TrimSpace(doc.Name) == "" {
			return false, errors.New("document name is missing or is not a string")
		}
		if strings.TrimSpace(doc.Format) == "" {
			return false, errors.New("format name is missing or is not a string")
		}
	} 
	return true, nil
}
