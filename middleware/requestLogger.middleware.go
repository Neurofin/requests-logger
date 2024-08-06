package middleware

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/Neurofin/requests-logger/logUtils"
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

		fmt.Println("request: ",req)
		fmt.Println("response: ", res)

		// Capture request headers
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
		logUtils.LogRequestResponse(req, requestBody.Bytes(), res, crw.Body.Bytes(), crw.Header(), start, end, traceId)

		return err
	}
}
