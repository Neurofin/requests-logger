package types

type S3EventBody struct {
	Records []S3EventRecord
}

type S3EventRecord struct {
	EventVersion      string                 `json:"eventVersion"`
	EventSource       string                 `json:"eventSource"`
	AwsRegion         string                 `json:"awsRegion"`
	EventTime         string                 `json:"eventTime"`
	EventName         string                 `json:"eventName"`
	UserIdentity      map[string]interface{} `json:"userIdentity"`
	RequestParameters map[string]interface{} `json:"requestParameters"`
	ResponseElements  map[string]interface{} `json:"responseElements"`
	S3                S3Map                  `json:"s3"`
}

type S3Map struct {
	S3SchemaVersion string   `json:"s3SchemaVersion"`
	ConfigurationId string   `json:"configurationId"`
	Bucket          S3Bucket `json:"bucket"`
	Object          S3Object `json:"object"`
}

type S3Bucket struct {
	Name          string                 `json:"name"`
	OwnerIdentity map[string]interface{} `json:"ownerIdentity"`
	Arn           string                 `json:"arn"`
}

type S3Object struct {
	Key       string `json:"key"`
	Size      int    `json:"size"`
	Etag      string `json:"eTag"`
	Sequencer string `json:"sequencer"`
}
