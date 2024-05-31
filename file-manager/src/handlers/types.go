package handler

type ResponseBody struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type createUploadUrlInput struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}

type getDownloadUrlInput struct {
	Bucket string `query:"bucket"`
	Key    string `query:"key"`
}
