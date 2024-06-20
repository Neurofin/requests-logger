package fileServiceTypes

type GetPresignedUrlInput struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}
