package loggerTypes

import (
	"time"

	"github.com/Neurofin/requests-logger/store/enum"
)

type PostLogInput struct {
	Type logTypeEnum.LogType 	`json:"type"`
	Stage logTypeEnum.StageType `json:"stage"`
	Data interface{} 			`json:"data"`
	TraceId string 				`json:"traceId"`
	Timestamp time.Time 		`json:"timestamp"`
}
