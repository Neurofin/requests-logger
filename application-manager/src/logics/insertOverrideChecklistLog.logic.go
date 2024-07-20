package logics

import (
	"application-manager/src/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertOverrideChecklistLog(overrideMeta map[string]interface{}, application primitive.ObjectID, checklistResult primitive.ObjectID, user primitive.ObjectID, note string) error {

	log := models.OverrideChecklistResultLog{
		Application:     application,
		ChecklistResult: checklistResult,
		OverrideMeta:    overrideMeta,
		User:            user,
		Note:            note,
	}

	if _, err := log.InsertLog(); err != nil {
		return err
	}

	return nil
}
