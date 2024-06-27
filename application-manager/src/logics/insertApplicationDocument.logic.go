package logics

import (
	"application-manager/src/models"
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertApplicationDocument(input types.InsertApplicationDocumentInput) (models.ApplicationDocumentModel, error) {

	applicationDocument := models.ApplicationDocumentModel{
		Id:          input.DocId,
		Name:        input.Name,
		Format:      input.Format,
		Application: input.Application,
		Status:      input.Status,
		S3Location:  input.S3Location,
	}

	dbResult, err := applicationDocument.InsertApplicationDocument()
	if err != nil {
		return applicationDocument, err
	}

	applicationDocument.Id = dbResult.Data.(primitive.ObjectID)
	documentResult, err := applicationDocument.GetApplicationDocumentById()
	if err != nil {
		return applicationDocument, err
	}

	applicationDocument = documentResult.Data.(models.ApplicationDocumentModel)
	return applicationDocument, nil
}
