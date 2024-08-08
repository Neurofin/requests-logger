package logUtils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/Neurofin/requests-logger/store/types"
	"github.com/Neurofin/requests-logger/store/enum"
)

func LogRequestResponse(req *http.Request, requestBody []byte, res *echo.Response, responseBody []byte, responseHeaders http.Header, start, end time.Time, traceId string, service string) {
	// Log request details
	requestLogData := loggerTypes.RequestLogType{
		TraceId:        traceId,
		RemoteIP:       req.RemoteAddr,
		Host:           req.Host,
		Method:         req.Method,
		URI:            req.RequestURI,
		UserAgent:      req.UserAgent(),
		RequestHeaders: req.Header,
		RequestBody:    string(requestBody),
		StartTime:      start.Format(time.RFC3339Nano),
		Timestamp:      time.Now(),
	}

	requestLogInput := loggerTypes.PostLogInput{
		Service:   service,
		Stage:     logTypeEnum.Start,
		Type:      logTypeEnum.API,
		Data:      requestLogData,
		TraceId:   traceId,
		Timestamp: time.Now(),
	}

	go PostLog(requestLogInput)

	// Log response details
	responseLogData := loggerTypes.ResponseLogType{
		TraceId:         traceId,
		Status:          res.Status,
		ResponseHeaders: responseHeaders,
		ResponseBody:    string(responseBody),
		StartTime:       start.Format(time.RFC3339Nano),
		EndTime:         end.Format(time.RFC3339Nano),
		Latency:         end.Sub(start).String(),
		Timestamp:       time.Now(),
	}

	responseLogInput := loggerTypes.PostLogInput{
		Service:   service,
		Stage: 	   logTypeEnum.End,	
		Type:      logTypeEnum.API,
		Data:      responseLogData,
		TraceId:   traceId,
		Timestamp: time.Now(),
	}

	go PostLog(responseLogInput)
}
