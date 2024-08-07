package logUtils

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/Neurofin/requests-logger/store/types"
	"github.com/Neurofin/requests-logger/store/enum"
)

func LogRequestResponse(req *http.Request, requestBody []byte, res *echo.Response, responseBody []byte, responseHeaders http.Header, start, end time.Time, traceId string) {
	// Log request details
	requestLogData := map[string]interface{}{
		"time":           time.Now().Format(time.RFC3339Nano),
		"traceId":        traceId,
		"remote_ip":      req.RemoteAddr,
		"host":           req.Host,
		"method":         req.Method,
		"uri":            req.RequestURI,
		"user_agent":     req.UserAgent(),
		"requestHeaders": req.Header,
		"requestBody":    string(requestBody),
		"startTime":      start.Format(time.RFC3339Nano),
		"stage":          logTypeEnum.Start,
	}

	requestLogInput := loggerTypes.PostLogInput{
		Type:      logTypeEnum.API,
		Data:      requestLogData,
		TraceId:   traceId,
		Timestamp: time.Now(),
	}

	go PostLog(requestLogInput)

	// Log response details
	responseLogData := map[string]interface{}{
		"time":            time.Now().Format(time.RFC3339Nano),
		"traceId":         traceId,
		"status":          res.Status,
		"responseHeaders": responseHeaders,
		"responseBody":    string(responseBody),
		"startTime":       start.Format(time.RFC3339Nano),
		"endTime":         end.Format(time.RFC3339Nano),
		"latency":         end.Sub(start).String(),
		"stage":           logTypeEnum.End,
	}

	responseLogInput := loggerTypes.PostLogInput{
		Type:      logTypeEnum.API,
		Data:      responseLogData,
		TraceId:   traceId,
		Timestamp: time.Now(),
	}

	go PostLog(responseLogInput)
}
