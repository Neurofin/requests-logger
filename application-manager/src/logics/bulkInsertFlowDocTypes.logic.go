package logics

import (
	"application-manager/src/models"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"application-manager/src/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BulkInsertFlowDocTypes(flow primitive.ObjectID, docTypes []types.DocumentTypeInput) ([]interface{}, error) {
	var documentTypesToBeInserted []interface{}
	for _, docType := range docTypes {
		docTypeDoc := models.DocumentTypeModel{
			Name:        docType.Name,
			Description: docType.Description,
			Flow:        flow,
			Uid:         docType.Uid,
		}
		docTypeDoc.CreatedAt = time.Now()
		docTypeDoc.UpdatedAt = time.Now()
		documentTypesToBeInserted = append(documentTypesToBeInserted, docTypeDoc)
	}

	if _, err := utils.BulkInsertToDb(documentTypesToBeInserted, store.DocumentTypeCollection); err != nil {
		return documentTypesToBeInserted, err
	}

	return documentTypesToBeInserted, nil
}
