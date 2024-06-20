package fileServiceTypes

type PresignedUrlResponseData struct {
	URL          string                 `json:"URL"`
	Method       string                 `json:"Method"`
	SignedHeader map[string]interface{} `json:"SignedHeader"`
}
