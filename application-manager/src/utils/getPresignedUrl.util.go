package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"application-manager/src/store/types"
)

func GetPresignedUrl(input types.CreateUploadUrlInput) (string, error) {
	client := &http.Client{}

	//change the URL to the PORT on which file-manager is running (URL : "http://localhost:<PORT>/presign")
	url := "http://localhost:3000/presign"

	reqBody, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get presigned URL, status code: %d", resp.StatusCode)
	}

	var responseBody types.ResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", err
	}

	return responseBody.Data.(string), nil
}
