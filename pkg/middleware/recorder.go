package middleware

import (
	"bytes"
	"net/http"
)

// responseRecorder is a custom http.ResponseWriter that captures
// the response body and status code for later inspection.
// Can be used to implement caching, logging, or other middleware
// that needs to analyze the response before sending it to the client.
type responseRecorder struct {
	http.ResponseWriter
	buffer bytes.Buffer
	status int
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.status = code
	rr.ResponseWriter.WriteHeader(code)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.buffer.Write(b)
	return rr.ResponseWriter.Write(b)
}

func (rr *responseRecorder) Status() int {
	return rr.status
}

func (rr *responseRecorder) Body() []byte {
	return rr.buffer.Bytes()
}
