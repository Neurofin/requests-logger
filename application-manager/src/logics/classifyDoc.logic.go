package logics

import (
	classifierService "application-manager/src/services/classifier"
	classifierServiceTypes "application-manager/src/services/classifier/store/types"
)

func ClassifyDoc(text string) (classifierServiceTypes.ClassData, error) {

	classData := classifierServiceTypes.ClassData{}

	response, err := classifierService.ClassifyDocument(text)
	if err != nil {
		return classData, err
	}
	
	classData = response.Data[0]

	return classData, nil
}
