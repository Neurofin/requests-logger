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

	mapDocumentNames := make(map[string][]interface{})
	for _, doc := range documents {
		mapDocumentNames[doc.Type] = append(mapDocumentNames[doc.Type], doc.Name)
	}

	mapDocuments := make(map[string][]interface{})
	for _, doc := range documents {
		mapDocuments[doc.Type] = append(mapDocuments[doc.Type], doc)
	}

	flowDocTypes, err := logics.GetFlowDocumentTypes(application.Flow)
	if err != nil {
		return result, err
	}

	for _, flowDocType := range flowDocTypes {
		docType := flowDocType.Uid
		element := map[string]interface{}{
			"name":     flowDocType.Name,
			"type":     flowDocType.Uid,
			"uploaded": false,
		}
		if _, ok := mapDocuments[docType]; ok {
			element["uploaded"] = true
			element["docNames"] = mapDocumentNames[docType]
			element["docs"] = mapDocuments[docType]
		}
		result = append(result, element)
	}

	return result, nil
}
