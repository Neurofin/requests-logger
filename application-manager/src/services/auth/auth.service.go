package authService

import (
	authStore "application-manager/src/services/auth/store/types"
	"application-manager/src/store/types"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func ValidateToken(token string) (authStore.TokenValidationResponseData, error) {
	responseData := authStore.TokenValidationResponseData{}

	authServiceUrl := os.Getenv("AUTH_SERVICE_URL") + "/user/validate"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, authServiceUrl, nil)
	if err != nil {
		return responseData, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	response, err := client.Do(req)
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
		errorMessage := responseBody.Message
		return responseData, errors.New(errorMessage)
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
