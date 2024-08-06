package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	logTypeEnum "github.com/Neurofin/requests-logger/store/enum"
	loggerTypes "github.com/Neurofin/requests-logger/store/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()
		req := c.Request()
		res := c.Response()

		traceId := req.Header.Get("traceId")
		if traceId == "" {
			traceId = uuid.New().String() // Generate a new UUID for the traceId
			req.Header.Set("traceI", traceId)
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
		crw := &CustomResponseWriter{ResponseWriter: res.Writer, body: resBody}
		res.Writer = crw

		// Proceed with the request
		err := next(c)

		end := time.Now()

		// Log the request and response details
		logRequest(req, requestBody.Bytes(), start, traceId)
		logResponse(res, crw.body.Bytes(), crw.Header(), start, end, traceId)

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

	postLog(logInput)
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

	postLog(logInput)
}

func postLog(logInput loggerTypes.PostLogInput) {
	logInputJSON, err := json.Marshal(logInput)
	if err != nil {
		log.Printf("Error marshaling log data: %v", err)
		return
	}

	logServiceURL := os.Getenv("LOG_SERVICE_URL")
	if logServiceURL == "" {
		log.Printf("LOG_SERVICE_URL is not set")
		return
	}

	go func() {
		resp, err := http.Post(logServiceURL, "application/json", bytes.NewBuffer(logInputJSON))
		if err != nil {
			log.Printf("Error posting log data: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Unexpected response from log service: %s", body)
		}
	}()
}
