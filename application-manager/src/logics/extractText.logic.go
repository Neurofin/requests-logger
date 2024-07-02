package logics

import (
	textractService "application-manager/src/services/textract"
	textractServiceTypes "application-manager/src/services/textract/store/types"
)

func ExtractTextLogic(sourceS3Path string, outputS3Path string) (string, error) {

	extractedText := ""

	response, err := textractService.ExtractText(textractServiceTypes.ExtractTextInput{
		SourceUrl:   sourceS3Path,
		OutputS3Path: outputS3Path,
	})
	if err != nil {
		return extractedText, err
	}
	
	extractedText = response.Text

	return extractedText, nil
}
