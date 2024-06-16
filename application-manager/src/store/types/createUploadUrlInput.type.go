package types

type CreateUploadUrlInput struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}
