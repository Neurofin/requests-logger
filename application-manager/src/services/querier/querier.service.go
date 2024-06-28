package querierService

import (
	querierServiceTypes "application-manager/src/services/querier/store/types"
	"application-manager/src/store/types"
	"application-manager/src/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func ResolveQuery(input querierServiceTypes.ResolveQueryInput) (map[string]interface{}, error) {
	responseData := make(map[string]interface{})

	serviceUrl := os.Getenv("QUERIER_SERVICE_URL") + "/resolve"

	response, err := utils.HttpJsonPost(serviceUrl, input)

	if err != nil {
		return responseData, err
	}

	defer response.Body.Close()

	responseBodyInBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return responseData, err
	}

	responseBody := types.ResponseBody{}
	err = json.Unmarshal(responseBodyInBytes, &responseBody)
	if err != nil {
		return responseData, err
	}

	if response.StatusCode != http.StatusOK {
		return responseData, errors.New("unknown error from querier service")
	}

	data := responseBody.Data.(map[string]interface{})

	marshelledData, err := json.Marshal(data)
	if err != nil {
		return responseData, err
	}

	err = json.Unmarshal(marshelledData, &responseData)
	if err != nil {
		return responseData, err
	}

	return responseData, nil
}

func Classify(input querierServiceTypes.ClassifyInput) (map[string]interface{}, error) {
	responseData := make(map[string]interface{})

	serviceUrl := os.Getenv("QUERIER_SERVICE_URL") + "/classify"

	response, err := utils.HttpJsonPost(serviceUrl, input)

	if err != nil {
		return responseData, err
	}

	defer response.Body.Close()

	responseBodyInBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return responseData, err
	}

	responseBody := types.ResponseBody{}
	err = json.Unmarshal(responseBodyInBytes, &responseBody)
	if err != nil {
		return responseData, err
	}

	if response.StatusCode != http.StatusOK {
		return responseData, errors.New("unknown error from querier service")
	}

	data := responseBody.Data.(map[string]interface{})

	marshelledData, err := json.Marshal(data)
	if err != nil {
		return responseData, err
	}

	err = json.Unmarshal(marshelledData, &responseData)
	if err != nil {
		return responseData, err
	}

	return responseData, nil
}
