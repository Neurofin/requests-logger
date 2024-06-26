package logics

import (
	"application-manager/src/models"
	"application-manager/src/store"
	"application-manager/src/store/types"
	"application-manager/src/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func BulkInsertChecklistItems(flow primitive.ObjectID, checklistItems []types.ChecklistItemInput) ([]interface{}, error) {
	var checklistItemsToBeInserted []interface{}
	for _, checklistItem := range checklistItems {
		checklistItemDoc := models.ChecklistItemModel{
			Name:         checklistItem.Name,
			Goal:         checklistItem.Goal,
			Rules:        checklistItem.Rules,
			Taxonomy:     checklistItem.Taxonomy,
			Prompt:       checklistItem.Prompt,
			GroupUid:     checklistItem.GroupUid,
			RequiredDocs: checklistItem.RequiredDocs, //TODO: Add validation with doctype
			Flow:         flow,
		}
		checklistItemDoc.CreatedAt = time.Now()
		checklistItemDoc.UpdatedAt = time.Now()
		checklistItemsToBeInserted = append(checklistItemsToBeInserted, checklistItemDoc)
	}
	if _, err := utils.BulkInsertToDb(checklistItemsToBeInserted, store.ChecklistItemCollection); err != nil {
		return checklistItemsToBeInserted, err
	}

	return checklistItemsToBeInserted, nil
}
