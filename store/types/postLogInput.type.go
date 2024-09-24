package loggerTypes

import (
	"errors"
	"time"

	"github.com/Neurofin/requests-logger/store/enum"
)

type PostLogInput struct {
	Service   string 				`json:"service"`
	Stage 	  logTypeEnum.StageType `json:"stage"`
	Type 	  logTypeEnum.LogType 	`json:"type"`
	Data 	  interface{} 			`json:"data"`
	TraceId   string 				`json:"traceId"`
	UserName  string       			`json:"userName"`
	UserEmail string       			`json:"userEmail"`
	Timestamp time.Time 			`json:"timestamp"`
}

type RequestLogType struct {
	TraceId        string              `json:"traceId"`
	RemoteIP       string              `json:"remote_ip"`
	Host           string              `json:"host"`
	Method         string              `json:"method"`
	URI            string              `json:"uri"`
	UserAgent      string              `json:"user_agent"`
	RequestHeaders map[string][]string `json:"requestHeaders"`
	RequestBody    string              `json:"requestBody"`
	StartTime      string              `json:"startTime"`
	Timestamp      time.Time           `json:"timestamp"`
}

type ResponseLogType struct {
	TraceId         string              `json:"traceId"`
	Status          int                 `json:"status"`
	ResponseHeaders map[string][]string `json:"responseHeaders"`
	ResponseBody    string              `json:"responseBody"`
	StartTime       string              `json:"startTime"`
	EndTime         string              `json:"endTime"`
	Latency         string              `json:"latency"`
	Timestamp       time.Time           `json:"timestamp"`
}

func (i *PostLogInput) Validate() error {
	if err := i.Type.Validate(); err != nil {
		return err
	}

	if i.TraceId == "" {
		return errors.New("traceId not found")
	}

	return nil
}