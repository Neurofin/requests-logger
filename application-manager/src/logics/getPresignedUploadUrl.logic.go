package logics

import (
	fileService "application-manager/src/services/file"
	fileServiceTypes "application-manager/src/services/file/store/types"
)

func GetPresignedUploadUrl(bucket string, key string) (string, error) {

	presignUrl := ""

	presignInput := fileServiceTypes.GetPresignedUrlInput{
		Bucket: bucket,
		Key:    key,
	}

	presignUrlResult, err := fileService.GetPresignedUploadUrl(presignInput)
	if err != nil {
		return presignUrl, err
	}

	return presignUrlResult.URL, nil
}
