package logics

import (
	"application-manager/src/models"
	"application-manager/src/store/types"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertApplication(input types.CreateApplicationLogicInput) (models.ApplicationModel, error) {
	newApplication := models.ApplicationModel{
		Org:                 input.Org,
		Flow:                input.Flow,
		Status:              "PENDING",
		TotalDocTypes:       input.TotalDocTypes,
		TotalChecklistItems: input.TotalChecklistItems,
		SignatureDocs:       input.SignatureDocs,
	}

	_, err := newApplication.GetApplication()
	if err == nil {
		return newApplication, errors.New("another app already exists with this name")
	}

	operationResult, err := newApplication.InsertApplication()
	if err != nil {
		return newApplication, err
	}

	insertedId := operationResult.Data.(primitive.ObjectID)
	newApplication.Id = insertedId

	output, err := newApplication.GetApplication()
	if err != nil {
		return newApplication, err
	}

	newApplication = output.Data.(models.ApplicationModel)
	return newApplication, nil
}
