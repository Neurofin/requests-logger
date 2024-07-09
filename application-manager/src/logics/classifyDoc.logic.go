package logics

import (
	classifierService "application-manager/src/services/classifier"
	classifierServiceTypes "application-manager/src/services/classifier/store/types"
	querierService "application-manager/src/services/querier"
	querierServiceTypes "application-manager/src/services/querier/store/types"
	"errors"
)

func ClassifyDoc(text string, isLLMBased bool, docPath string, prompt string, isFileBased bool) (classifierServiceTypes.ClassData, error) {

	classData := classifierServiceTypes.ClassData{}

	if !isLLMBased {
		response, err := classifierService.ClassifyDocument(text)
		if err != nil {
			return classData, err
		}

		classData = response.Data[0]
	} else {
		docFormat := ""
		if isFileBased {
			docFormat = "application/json"
		}
		response, err := querierService.Classify(querierServiceTypes.ClassifyInput{
			DocPath:   docPath,
			Prompt:    prompt,
			DocFormat: docFormat,
		})
		if err != nil {
			return classData, err
		}

		if response == nil {
			return classData, errors.New("error, classification sent null")
		}

		name, ok := response["category"].(string)
		if !ok {
			return classData, errors.New("error from classification")
		}
		classData.Name = name
		classData.Score = response["confidence"].(float64)
	}

	return classData, nil
}
