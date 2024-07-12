package signatureService

import (
	signatureServiceTypes "application-manager/src/services/signature/store/types"
	"application-manager/src/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func ExtractSignatures(input signatureServiceTypes.SignatureInput) (map[string][]string, error) {
	responseData := map[string][]string{}

	signatureServiceUrl := os.Getenv("SIGNATURE_SERVICE_URL") + "/extract"

	response, err := utils.HttpJsonPost(signatureServiceUrl, input)

	if err != nil {
		return responseData, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return responseData, errors.New("unknown error from signature service")
	}

	responseBodyInBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return responseData, err
	}

	err = json.Unmarshal(responseBodyInBytes, &responseData)
	if err != nil {
		return responseData, err
	}

	return responseData, nil
}
