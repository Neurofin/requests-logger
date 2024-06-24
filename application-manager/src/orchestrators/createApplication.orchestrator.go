package orchestrators

import (
	"application-manager/src/logics"
	"application-manager/src/store/types"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateApplication(org primitive.ObjectID, flowUId string, numberOfDocs int) (map[string]interface{}, error) {

	output := map[string]interface{}{} //TODO: Update type

	flow, err := logics.GetFlow(flowUId)
	if err != nil {
		return output, err
	}

	logicInput := types.CreateApplicationLogicInput{
		Org:  org,
		Flow: flow.Id,
	}

	newApplication, err := logics.CreateApplication(logicInput)
	if err != nil {
		return output, err
	}

	documents := []map[string]interface{}{} //TODO: Update type
	for i := 1; i <= numberOfDocs; i++ {
		appId := newApplication.Id
		input := types.AddApplicationDocumentInput{
			ApplicationId: appId.Hex(),
		}
		result, err := AddApplicationDocument(input)
		if err != nil {
			return output, err
		}
		documents = append(documents, result)
	}

	output["application"] = newApplication
	output["documents"] = documents
	return output, nil
}
