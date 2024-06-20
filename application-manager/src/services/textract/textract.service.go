package textractService

import (
	textractServiceTypes "application-manager/src/services/textract/store/types"
	"application-manager/src/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func ExtractText(input textractServiceTypes.ExtractTextInput) (textractServiceTypes.ExtractTextResponseData, error) {
	responseData := textractServiceTypes.ExtractTextResponseData{}

	textractServiceUrl := os.Getenv("TEXTRACT_SERVICE_URL") + "/textract" + "?file_source=" + input.SourceUrl + "&file_output=" + input.OutputS3Path

	response, err := utils.HttpJsonPost(textractServiceUrl, nil)

	if err != nil {
		return responseData, err
	}

	defer response.Body.Close()

	responseBodyInBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return responseData, err
	}

	err = json.Unmarshal(responseBodyInBytes, &responseData)
	if err != nil {
		return responseData, err
	}

	if response.StatusCode != http.StatusOK {
		return responseData, errors.New("unknown error from textract service")
	}

	return responseData, nil
}
