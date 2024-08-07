package types

type ResponseBody struct {
	TraceId string `json:"traceId"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
