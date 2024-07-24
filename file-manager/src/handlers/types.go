package handlers

import (
	"errors"
	"strings"
)

type ResponseBody struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type createUploadUrlInput struct {
	Bucket      string `json:"bucket"`
	Key         string `json:"key"`
	ContentType string `json:"contentType"`
}

type getDownloadUrlInput struct {
	Bucket      string `query:"bucket"`
	Key         string `query:"key"`
	ContentType string `query:"contentType"`
}

type deleteFileInput struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

func validateCreateUploadUrlInput(input createUploadUrlInput) error {
	if strings.TrimSpace(input.Bucket) == "" {
		return errors.New("bucket is required and cannot be empty")
	}
	if strings.TrimSpace(input.Key) == "" {
		return errors.New("key is required and cannot be empty")
	}
	return nil
}

func validateGetDownloadUploadUrlInput(input getDownloadUrlInput) error {
	if strings.TrimSpace(input.Bucket) == "" {
		return errors.New("bucket is required and cannot be empty")
	}
	if strings.TrimSpace(input.Key) == "" {
		return errors.New("key is required and cannot be empty")
	}
	return nil
}

func validateDeleteFileInput(input deleteFileInput) error {
	if strings.TrimSpace(input.Bucket) == "" {
		return errors.New("bucket is required and cannot be empty")
	}
	if strings.TrimSpace(input.Key) == "" {
		return errors.New("key is required and cannot be empty")
	}
	return nil
}