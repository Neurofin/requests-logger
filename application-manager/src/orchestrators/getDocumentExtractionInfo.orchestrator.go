package orchestrators

import (
	"application-manager/src/logics"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetDocumentExtractionInfo(appId string, org primitive.ObjectID) ([]interface{}, error) {

	result := []interface{}{} //TODO: Add type

	application, err := logics.GetApplication(appId, org)
	if err != nil {
		return result, err
	}

	documents, err := logics.GetApplicationDocuments(application.Id)
	if err != nil {
		return result, err
	}

	mapDocuments := make(map[string]interface{})
	for _, doc := range documents {
		mapDocuments[doc.Type] = doc
	}

	flowDocTypes, err := logics.GetFlowDocumentTypes(application.Flow)
	if err != nil {
		return result, err
	}

	for _, flowDocType := range flowDocTypes {
		Uid := flowDocType.Uid
		element := map[string]interface{}{
			"name":     flowDocType.Name,
			"type":     flowDocType.Uid,
			"uploaded": false,
		}
		if _, ok := mapDocuments[Uid]; ok {
			element["uploaded"] = true
		}
		result = append(result, element)
	}

	return result, nil
}
