package types

type CreateDocumentInput struct {
	DocumentName  string `json:"documentName" validate:"required"`
	ApplicationId string `json:"applicationId" validate:"required"`
}
