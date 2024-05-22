package logging

import (
	"CRM/go/apiGateway/internal/logger"
	"fmt"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	//lrw.ResponseWriter.WriteHeader(code)
}

func colorStatusCode(code int) string {
	var colorStart string
	var colorEnd string = "\033[0m" // Resets the color

	if code >= 200 && code < 300 {
		colorStart = "\033[32m" // Green for 2xx success codes
	} else if code >= 400 && code < 500 {
		colorStart = "\033[33m" // Yellow for 4xx client errors
	} else if code >= 500 {
		colorStart = "\033[31m" // Red for 5xx server errors
	} else {
		colorStart = "\033[0m" // Default terminal color
	}

	return fmt.Sprintf("%s%d%s", colorStart, code, colorEnd)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := NewLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		logger.CreateLog("info", fmt.Sprintf("| %v | %v | %v | %v |", colorStatusCode(lrw.statusCode), time.Since(start), r.Method, r.RequestURI))
	})
}
