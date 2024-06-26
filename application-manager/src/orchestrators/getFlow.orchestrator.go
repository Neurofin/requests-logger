package orchestrators

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/models"
)

func GetFlow(flowId string) (map[string]interface{}, error) {
	result := map[string]interface{}{}

	flow := models.FlowModel{
		Uid: flowId,
	}

	getResult, err := flow.GetFlow()
	if err != nil {
		return result, err
	}

	flow = getResult.Data.(models.FlowModel)

	getOperationResult, err := dbHelpers.GetFlowDocumentTypes(flow.Id)
	if err != nil {
		return result, err
	}

	docTypes := getOperationResult.Data.([]models.DocumentTypeModel)

	checklistItemsResult, err := dbHelpers.GetFlowChecklist(flow.Id)
	if err != nil {
		return result, err
	}

	checklistItems := checklistItemsResult.Data.([]models.ChecklistItemModel)

	result["flow"] = flow
	result["docTypes"] = docTypes
	result["checklistItems"] = checklistItems
	return result, nil
}
