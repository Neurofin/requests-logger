package logUtils

import (
	"net/http"
	"time"

	"github.com/Neurofin/requests-logger/store/enum"
	"github.com/Neurofin/requests-logger/store/types"
	"github.com/labstack/echo/v4"
)

func LogRequestResponse(req *http.Request, requestBody []byte, res *echo.Response, responseBody []byte, responseHeaders http.Header, start, end time.Time, traceId string, service string, user types.TokenValidationResponseData) {
	// Log request details
	requestLogData := types.RequestLogType{
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

	requestLogInput := types.PostLogInput{
		Service:   service,
		Stage:     logTypeEnum.Start,
		Type:      logTypeEnum.API,
		Data:      requestLogData,
		TraceId:   traceId,
		FirstName: user.FirstName,
		Email: 	   user.Email,
		Timestamp: time.Now(),
	}

	go PostLog(requestLogInput)

	// Log response details
	responseLogData := types.ResponseLogType{
		TraceId:         traceId,
		Status:          res.Status,
		ResponseHeaders: responseHeaders,
		ResponseBody:    string(responseBody),
		StartTime:       start.Format(time.RFC3339Nano),
		EndTime:         end.Format(time.RFC3339Nano),
		Latency:         end.Sub(start).String(),
		Timestamp:       time.Now(),
	}

	responseLogInput := types.PostLogInput{
		Service:   service,
		Stage: 	   logTypeEnum.End,	
		Type:      logTypeEnum.API,
		Data:      responseLogData,
		TraceId:   traceId,
		FirstName: user.FirstName,
		Email: 	   user.Email,
		Timestamp: time.Now(),
	}

	go PostLog(responseLogInput)
}
