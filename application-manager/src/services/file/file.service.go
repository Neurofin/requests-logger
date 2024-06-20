package fileService

import (
	fileStore "application-manager/src/services/file/store/types"
	"application-manager/src/store/types"
	"application-manager/src/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func GetPresignedUploadUrl(input fileStore.GetPresignedUrlInput) (fileStore.PresignedUrlResponseData, error) {
	responseData := fileStore.PresignedUrlResponseData{}

	authServiceUrl := os.Getenv("FILE_SERVICE_URL") + "/presign"

	response, err := utils.HttpJsonPost(authServiceUrl, input)

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
		message := responseBody.Message
		if message == "" {
			message = "Error from file-manager"
		}
		return responseData, errors.New(message)
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

func GetPresignedDownloadUrl(input fileStore.GetPresignedUrlInput) (fileStore.PresignedUrlResponseData, error) {
	responseData := fileStore.PresignedUrlResponseData{}

	authServiceUrl := os.Getenv("FILE_SERVICE_URL") + "/presign" + "?bucket=" + input.Bucket + "&key=" + input.Key

	response, err := utils.HttpGet(authServiceUrl)

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
		message := responseBody.Message
		if message == "" {
			message = "Error from file-manager"
		}
		return responseData, errors.New(message)
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
