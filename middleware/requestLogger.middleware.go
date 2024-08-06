package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/Neurofin/requests-logger/store/enum"
	"github.com/Neurofin/requests-logger/store/types"
	"github.com/google/uuid"
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
		req := c.Request()
		res := c.Response()

		// Generate a traceID for the entire request-response cycle
		traceId, ok := c.Get("traceId").(string)
		if !ok || traceId == "" {
			traceId = uuid.New().String() // Generate a new UUID for the traceID
			c.Set("traceID", traceId)
			req.Header.Set("traceId", traceId) // Add traceID to request header
		}
		res.Header().Set("traceId", traceId) // Add traceID to response header

		var requestBody []byte
		if req.Body != nil {
			requestBody, _ = io.ReadAll(req.Body)
			req.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		logRequest(c, traceId)

		resBody := new(bytes.Buffer)
		crw := &CustomResponseWriter{ResponseWriter: res.Writer, body: resBody}
		res.Writer = crw

		err := next(c)

		logResponse(c, traceId)

		return err
	}
}

func logRequest(c echo.Context, traceID string) {
	req := c.Request()
	logInput := loggerTypes.PostLogInput{
		Type:      logTypeEnum.API,
		Data:      req,
		Stage:     logTypeEnum.Start,
		TraceId:   traceID,
		Timestamp: time.Now(),
	}

	postLog(logInput)
}

func logResponse(c echo.Context, traceID string) {
	res := c.Response()
	logInput := loggerTypes.PostLogInput{
		Type:      logTypeEnum.API,
		Data:      res,
		Stage:     logTypeEnum.End,
		TraceId:   traceID,
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
