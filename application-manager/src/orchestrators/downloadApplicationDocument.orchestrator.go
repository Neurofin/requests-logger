package orchestrators

import (
	"application-manager/src/logics"
	"net/url"
)

func DownloadApplicationDocument(appId string, docId string) (map[string]string, error) {

	data := map[string]string{}

	document, err := logics.GetApplicationDocumet(appId, docId)
	if err != nil {
		return data, err
	}

	s3Location := document.S3Location
	parsed, err := url.Parse(s3Location)
	if err != nil {
		return data, err
	}

	bucket := parsed.Host
	key := parsed.Path

	urlResult, err := logics.GetPresignedDownloadUrl(bucket, key)
	if err != nil {
		return data, err
	}

	data["presignUrl"] = urlResult
	return data, nil
}
