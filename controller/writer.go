package controller

import (
	"bytes"
	"io"
	"net/http"
	"sync/atomic"
)

// ResponseWriter - write a response
type ResponseWriter struct {
	statusCode int32
	header     http.Header
	body       *bytes.Buffer
	written    int64
}

// NewResponseWriter - create a new response writer
func newResponseWriter() *ResponseWriter {
	w := new(ResponseWriter)
	w.header = make(http.Header)
	w.body = new(bytes.Buffer)
	return w
}

// SetStatusCode - return the response status code
func (w *ResponseWriter) SetStatusCode(code int) {
	atomic.StoreInt32(&w.statusCode, int32(code))
}

// StatusCode - return the response status code
func (w *ResponseWriter) StatusCode() int {
	return int(atomic.LoadInt32(&w.statusCode))
}

// Header - return the response http.Header
func (w *ResponseWriter) Header() http.Header {
	return w.header
}

// Body - return the response body
func (w *ResponseWriter) Body() []byte {
	return w.body.Bytes()
}

// Written - return bytes written
func (w *ResponseWriter) Written() int64 {
	return w.written
}

// Write - write the response body
func (w *ResponseWriter) Write(p []byte) (int, error) {
	w.written += int64(len(p))
	return w.body.Write(p)
}

// WriteHeader - write the response status code
func (w *ResponseWriter) WriteHeader(statusCode int) {
	atomic.CompareAndSwapInt32(&w.statusCode, 0, int32(statusCode))
}

// Response - return the response
func (w *ResponseWriter) Response() *http.Response {
	r := new(http.Response)
	if w.statusCode == 0 {
		r.StatusCode = http.StatusOK
	} else {
		r.StatusCode = int(w.statusCode)
	}
	if r.StatusCode >= http.StatusOK && r.StatusCode <= http.StatusMultipleChoices {
		r.Header = w.header
		r.Body = io.NopCloser(bytes.NewReader(w.body.Bytes()))
	}
	return r
}
