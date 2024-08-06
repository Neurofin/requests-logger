package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/Neurofin/requests-logger/store/enum"
	"github.com/Neurofin/requests-logger/store/types"
)

// CustomResponseWriter captures the response body
type CustomResponseWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggingMiddleware logs the entire request and response
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

		// Log the request and response details
		logRequest(c, traceId, req, requestBody.Bytes())
		logResponse(c, traceId, res, resBody.Bytes())

		return err
	}
}

// serializeObject serializes an object using reflection
func serializeObject(obj interface{}) (string, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	
	// Create a map to hold fields
	result := make(map[string]interface{})
	
	// Iterate through the fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i)

		if value.Kind() == reflect.Ptr && !value.IsNil() {
			value = value.Elem()
		}

		// Only include exported fields
		if field.PkgPath == "" {
			result[field.Name] = value.Interface()
		}
	}

	// Convert map to JSON
	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// logRequest logs the entire request details
func logRequest(c echo.Context, traceID string, req *http.Request, requestBody []byte) {
	reqDetailsJSON, err := serializeObject(req)
	if err != nil {
		log.Printf("Error serializing request details: %v", err)
		return
	}

	logInput := loggerTypes.PostLogInput{
		Type:      logTypeEnum.API,
		Data:      json.RawMessage(reqDetailsJSON),
		Stage:     logTypeEnum.Start,
		TraceId:   traceID,
		Timestamp: time.Now(),
	}

	postLog(logInput)
}

// logResponse logs the entire response details
func logResponse(c echo.Context, traceID string, res *echo.Response, responseBody []byte) {
	resDetailsJSON, err := serializeObject(res)
	if err != nil {
		log.Printf("Error serializing response details: %v", err)
		return
	}

	logInput := loggerTypes.PostLogInput{
		Type:      logTypeEnum.API,
		Data:      json.RawMessage(resDetailsJSON),
		Stage:     logTypeEnum.End,
		TraceId:   traceID,
		Timestamp: time.Now(),
	}

	postLog(logInput)
}

// postLog sends the log data to the log service
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
