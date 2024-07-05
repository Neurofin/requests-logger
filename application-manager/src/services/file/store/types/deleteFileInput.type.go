package fileServiceTypes

type DeleteFileInput struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}
