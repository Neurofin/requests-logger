package logics

import (
	"errors"
	logDb "logger/src/db/log"
	loggerTypes "logger/src/store/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertLogLogic(input loggerTypes.PostLogInput) (logDb.Model, error) {

	newlog := logDb.Model{
		Type:      input.Type,
		Data:      input.Data,
		TraceId:   input.TraceId,
		Timestamp: input.Timestamp,
	}

	operationResult, err := newlog.Insert()
	if err != nil {
		return newlog, err
	}

	logId, ok := operationResult.Data.(primitive.ObjectID)
	if !ok {
		return newlog, errors.New("invalid object id")
	}

	operationResult, err = logDb.FetchById(logId)
	if err != nil {
		return newlog, err
	}

	newlog = operationResult.Data.(logDb.Model)
	return newlog, nil
}
