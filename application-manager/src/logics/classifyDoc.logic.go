package logics

import (
	classifierService "application-manager/src/services/classifier"
	classifierServiceTypes "application-manager/src/services/classifier/store/types"
	querierService "application-manager/src/services/querier"
	querierServiceTypes "application-manager/src/services/querier/store/types"
)

func ClassifyDoc(text string, isLLMBased bool, docPath string, prompt string) (classifierServiceTypes.ClassData, error) {

	classData := classifierServiceTypes.ClassData{}

	if !isLLMBased {
		response, err := classifierService.ClassifyDocument(text)
		if err != nil {
			return classData, err
		}

		classData = response.Data[0]
	} else {
		response, err := querierService.Classify(querierServiceTypes.ClassifyInput{
			DocPath:   docPath,
			Prompt:    prompt,
			DocFormat: "application/json",
		})
		if err != nil {
			return classData, err
		}

		classData.Name = response["category"].(string)
		classData.Score = response["confidence"].(float64)
	}

	return classData, nil
}
