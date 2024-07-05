package handlers

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
