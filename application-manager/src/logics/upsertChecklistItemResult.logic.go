package logics

import (
	"application-manager/src/models"
	"application-manager/src/store/types"
)

func UpsertChecklistItemResultLogic(queryResult models.ChecklistItemResultModel) (types.DbOperationResult, error) {
	result := types.DbOperationResult{
		OperationSuccess: false,
	}

	operationResult, err := queryResult.FindChecklistItemResult()
	if err != nil {
		if _, err := queryResult.InsertChecklistItemResult(); err != nil {
			return result, err
		}

		operationResult, err := queryResult.FindChecklistItemResult()
		if err != nil {
			return result, err
		}

		result.OperationSuccess = true
		result.Data = operationResult.Data
		return result, nil
	}

	data := operationResult.Data.(models.ChecklistItemResultModel)
	data.Result = queryResult.Result
	if _, err := data.UpdateChecklistItemResult(); err != nil {
		return result, nil
	}

	result.OperationSuccess = true
	result.Data = data
	return result, nil
}
