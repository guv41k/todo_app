package core_http_response

import "net/http"

var (
	StatusCodeUnitialized = -1
)

type ResponseWritter struct {
	http.ResponseWriter
	statusCode int
}

func NewResonseWriter(w http.ResponseWriter) *ResponseWritter {
	return &ResponseWritter{
		ResponseWriter: w,
		statusCode:     StatusCodeUnitialized,
	}
}

func (rw *ResponseWritter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWritter) GetStatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUnitialized {
		panic("no status code")
	}
	return rw.statusCode
}
