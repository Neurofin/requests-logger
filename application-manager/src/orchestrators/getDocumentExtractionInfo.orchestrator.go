package orchestrators

import (
	"application-manager/src/logics"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func validateGetDocumentExtractionInfoInput(appId string) (bool, error) {
	if strings.TrimSpace(appId) == "" {
		return false, errors.New("appId is missing or is not a string")
	}

	if _, err := primitive.ObjectIDFromHex(appId); err != nil {
		return false, errors.New("appId is not in valid ObjectID format")
	}

	return true, nil
}

func GetDocumentExtractionInfo(appId string, org primitive.ObjectID) ([]interface{}, error) {

	if valid, err := validateGetDocumentExtractionInfoInput(appId); !valid {
		return nil, err
	}

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
