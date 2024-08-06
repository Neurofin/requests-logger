package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/Neurofin/requests-logger/logUtils"
	logTypeEnum "github.com/Neurofin/requests-logger/store/enum"
	loggerTypes "github.com/Neurofin/requests-logger/store/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		req := c.Request()
		res := c.Response()

		traceId := req.Header.Get("traceId")
		if traceId == "" {
			traceId = uuid.New().String() // Generate a new UUID for the traceId
			req.Header.Set("traceId", traceId)
			c.Set("traceId",traceId)
		}

		// Capture request body
		var requestBody bytes.Buffer
		if req.Body != nil {
			body, _ := io.ReadAll(req.Body)
			requestBody.Write(body)
			req.Body = io.NopCloser(&requestBody)
		}

		// Create a custom response writer to capture response body
		resBody := new(bytes.Buffer)
		crw := &loggerTypes.CustomResponseWriter{ResponseWriter: res.Writer, Body: resBody}
		res.Writer = crw

		// Proceed with the request
		err := next(c)

		end := time.Now()

		// Log the request and response details
		logRequest(req, requestBody.Bytes(), start, traceId)
		logResponse(res, crw.Body.Bytes(), crw.Header(), start, end, traceId)

		return err
	}
}

func logRequest(req *http.Request, requestBody []byte, start time.Time, traceId string) {
	logData := map[string]interface{}{
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
	}

	logInput := loggerTypes.PostLogInput{
		Type:      logTypeEnum.API,
		Stage:     logTypeEnum.Start,
		Data:      logData,
		TraceId:   traceId,
		Timestamp: time.Now(),
	}

	logUtils.PostLog(logInput)
}

func logResponse(res *echo.Response, responseBody []byte, responseHeaders http.Header, start, end time.Time, traceId string) {
	logData := map[string]interface{}{
		"time":            time.Now().Format(time.RFC3339Nano),
		"traceId":         traceId,
		"status":          res.Status,
		"responseHeaders": responseHeaders,
		"responseBody":    string(responseBody),
		"startTime":       start.Format(time.RFC3339Nano),
		"endTime":         end.Format(time.RFC3339Nano),
		"latency":         end.Sub(start).String(),
	}

	logInput := loggerTypes.PostLogInput{
		Type:      logTypeEnum.API,
		Stage:     logTypeEnum.End,
		Data:      logData,
		TraceId:   traceId,
		Timestamp: time.Now(),
	}

	logUtils.PostLog(logInput)
}

// func postLog(logInput loggerTypes.PostLogInput) {
// 	logInputJSON, err := json.Marshal(logInput)
// 	if err != nil {
// 		log.Printf("Error marshaling log data: %v", err)
// 		return
// 	}

// 	logServiceURL := os.Getenv("LOG_SERVICE_URL")
// 	if logServiceURL == "" {
// 		log.Printf("LOG_SERVICE_URL is not set")
// 		return
// 	}

// 	go func() {
// 		resp, err := http.Post(logServiceURL, "application/json", bytes.NewBuffer(logInputJSON))
// 		if err != nil {
// 			log.Printf("Error posting log data: %v", err)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		if resp.StatusCode != http.StatusCreated {
// 			body, _ := io.ReadAll(resp.Body)
// 			log.Printf("Unexpected response from log service: %s", body)
// 		}
// 	}()
// }
