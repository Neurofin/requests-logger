package textractService

import (
	signatureServiceTypes "application-manager/src/services/signature/store/types"
	"application-manager/src/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func GetSignatures(input signatureServiceTypes.SignatureInput) ([]string, error) {
	responseData := []string{}

	textractServiceUrl := os.Getenv("SIGNATURE_SERVICE_URL") + "/sign"

	response, err := utils.HttpJsonPost(textractServiceUrl, input)

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
