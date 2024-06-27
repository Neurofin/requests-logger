package types

type AddApplicationDocumentsInput struct {
	ApplicationId    string             `json:"applicationId"`
	DocsToBeUploaded []DocsToBeUploaded `json:"docsToBeUploaded"`
}
