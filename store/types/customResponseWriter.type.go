package types

import (
	"bytes"
	"net/http"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	Body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}