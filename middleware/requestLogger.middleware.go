package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/Neurofin/requests-logger/logUtils"
	"github.com/Neurofin/requests-logger/store/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func LoggingMiddleware(service string) echo.MiddlewareFunc {
	return func (next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			traceId := req.Header.Get("traceId")
			if traceId == "" {
				traceId = uuid.New().String()
				req.Header.Set("traceId", traceId)
				c.Set("traceId", traceId)
			}

			var userDetails types.TokenValidationResponseData
			
			user := c.Get("user")
			if user != nil {
				userDetails = c.Get("user").(types.TokenValidationResponseData)
			}

			// Capture request headers
			// Capture request body
			var requestBody bytes.Buffer
			if req.Body != nil {
				body, _ := io.ReadAll(req.Body)
				requestBody.Write(body)
				req.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			resBody := new(bytes.Buffer)
			crw := &types.CustomResponseWriter{ResponseWriter: res.Writer, Body: resBody}
			res.Writer = crw

			err := next(c)

			end := time.Now()

			// Log the request and response details asynchronously
			go logUtils.LogRequestResponse(req, requestBody.Bytes(), res, crw.Body.Bytes(), crw.Header(), start, end, traceId, service, userDetails)

			return err
		}
	}
}