package logUtils

import (
	"time"

	"github.com/Neurofin/requests-logger/store/enum"
	"github.com/Neurofin/requests-logger/store/types"
	"github.com/google/uuid"
)

func postErrorLog(service string, method string, err error) {

	traceId := uuid.New().String()

	errorLogData := map[string]interface{}{
		"traceId": traceId,
		"error":   err.Error(),
		"method":  method,
	}

	errorLogInput := loggerTypes.PostLogInput{
		Service:   service,
		Stage:     logTypeEnum.End,
		Type:      logTypeEnum.Error,
		Data:      errorLogData,
		TraceId:   traceId,
		Timestamp: time.Now(),
	}

	go PostLog(errorLogInput)

}