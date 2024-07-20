package types

import (
	"errors"
	"strings"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddApplicationDocumentsInput struct {
	ApplicationId    string             `json:"applicationId"`
	DocsToBeUploaded []DocsToBeUploaded `json:"docsToBeUploaded"`
}

func (i *AddApplicationDocumentsInput) Validate() (bool, error) {
	if i.ApplicationId == "" {
		return false, errors.New("ApplicationId is missing or is not a String")
	}
	
	if _, err:= primitive.ObjectIDFromHex(i.ApplicationId); err!=nil{
		return false, errors.New("ApplicationId is not in valid ObjectID format")
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