package orchestrators

import (
	"application-manager/src/dbHelpers"
	"application-manager/src/logics"
	"application-manager/src/models"
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateApplication(org primitive.ObjectID, flowUId string, docsToBeUploaded []types.DocsToBeUploaded) (map[string]interface{}, error) {

	output := map[string]interface{}{} //TODO: Update type

	flow, err := logics.GetFlow(flowUId)
	if err != nil {
		return output, err
	}

	flowDocTypesResults, err := dbHelpers.GetFlowDocumentTypes(flow.Id)
	if err != nil {
		return output, err
	}
	flowDocTypes := flowDocTypesResults.Data.([]models.DocumentTypeModel)
	totalDocTypes := len(flowDocTypes)

	flowChecklistItemsResult, err := dbHelpers.GetFlowChecklist(flow.Id)
	if err != nil {
		return output, err
	}
	flowChecklistItems := flowChecklistItemsResult.Data.([]models.ChecklistItemModel)
	totalChecklistItems := len(flowChecklistItems) - 1

	logicInput := types.CreateApplicationLogicInput{
		Org:                 org,
		Flow:                flow.Id,
		TotalDocTypes:       totalDocTypes,
		TotalChecklistItems: totalChecklistItems,
	}

	newApplication, err := logics.InsertApplication(logicInput)
	if err != nil {
		return output, err
	}

	documents := []map[string]interface{}{} //TODO: Update type
	for _, doc := range docsToBeUploaded {
		appId := newApplication.Id
		input := types.AddApplicationDocumentInput{
			ApplicationId: appId.Hex(),
			Name:          doc.Name,
			Format:        doc.Format,
		}
		result, err := AddApplicationDocument(input)
		if err != nil {
			return output, err
		}
		documents = append(documents, result)
	}

	// Listener for updates to the documents
	documentIds := []primitive.ObjectID{}
	for _, doc := range documents {
		document := doc["document"].(models.ApplicationDocumentModel)
		documentIds = append(documentIds, document.Id)
	}
	go ApplicationDocumentClassificationEventListener(newApplication.Id, documentIds)

	output["application"] = newApplication
	output["documents"] = documents
	return output, nil
}
