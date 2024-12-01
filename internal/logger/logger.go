package logger

import (
	"log"
	"net/http"
	"time"
)

// logResponseWriter is a custom wrapper around http.ResponseWriter to capture status code and response size
type logResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

// WriteHeader captures the status code
func (lrw *logResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the number of bytes written in the response
func (lrw *logResponseWriter) Write(p []byte) (n int, err error) {
	n, err = lrw.ResponseWriter.Write(p)
	lrw.bytesWritten += n
	return n, err
}

// LogRequest logs the details of each request and response.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &logResponseWriter{ResponseWriter: w}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		log.Printf("[Request] %s %s %d %s", r.Method, r.URL.Path, lrw.statusCode, duration)
	})
}
