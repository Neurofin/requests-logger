package classifierService

import (
	classifierServiceTypes "application-manager/src/services/classifier/store/types"
	"application-manager/src/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func ClassifyDocument(text string) (classifierServiceTypes.ClassifierResponseData, error) {
	responseData := classifierServiceTypes.ClassifierResponseData{}

	classifierServiceUrl := os.Getenv("CLASSIFIER_SERVICE_URL") + "/classify"

	response, err := utils.HttpJsonPost(classifierServiceUrl, classifierServiceTypes.ClassifierInput{
		Text: text,
	})

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
		return responseData, errors.New("unknown error from classifier")
	}

	return responseData, nil
}
