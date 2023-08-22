package test

import (
	"bytes"
	"github.com/cocktailrobots/openbar-server/pkg/util"
	"net/http"
)

var _ http.ResponseWriter = &ResponseWriter{}

// ResponseWriter is a mock implementation of http.ResponseWriter
type ResponseWriter struct {
	body       *bytes.Buffer
	header     http.Header
	returnCode *int
	writeError error
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		body:   bytes.NewBuffer(nil),
		header: make(http.Header),
	}
}

// Header returns the header map
func (wr *ResponseWriter) Header() http.Header {
	if wr.header == nil {
		wr.header = make(http.Header)
	}

	return wr.header
}

// Write writes the data to the buffer
func (wr *ResponseWriter) Write(data []byte) (int, error) {
	if wr.writeError != nil {
		return 0, wr.writeError
	}

	if wr.returnCode == nil {
		wr.returnCode = util.Ptr(http.StatusOK)
	}

	return wr.body.Write(data)
}

// WriteHeader sets the return code
func (wr *ResponseWriter) WriteHeader(statusCode int) {
	if wr.returnCode != nil {
		panic("Response code already set")
	} else if wr.returnCode == nil {
		wr.returnCode = &statusCode
	}
}

// Body returns the body as a byte array
func (wr *ResponseWriter) Body() []byte {
	return wr.body.Bytes()
}

// StatusCode returns the return code
func (wr *ResponseWriter) StatusCode() int {
	if wr.returnCode == nil {
		panic("Response code not set")
	}

	return *wr.returnCode
}

// SetWriteError sets the error that will be returned by all calls to Write until it is changed
func (wr *ResponseWriter) SetWriteError(err error) {
	wr.writeError = err
}
