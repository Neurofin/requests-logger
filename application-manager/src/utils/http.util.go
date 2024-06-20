package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func HttpJsonPost(url string, body interface{}) (*http.Response, error) {

	response := &http.Response{}

	client := &http.Client{}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return response, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return response, err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err = client.Do(req)
	if err != nil {
		return response, err
	}

	return response, nil
}

func HttpGet(url string) (*http.Response, error) {

	response := &http.Response{}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return response, err
	}
	req.Header.Set("Content-Type", "application/json")

	response, err = client.Do(req)
	if err != nil {
		return response, err
	}

	return response, nil
}
