package textractServiceTypes

type ExtractTextInput struct {
	SourceUrl    string `json:"sourceUrl"`
	OutputS3Path string `json:"outputS3Path"`
}
