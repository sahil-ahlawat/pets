package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	bytesWritten int
}

func (rw *responseWriterWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.New().String()
		}

		w.Header().Set("X-Request-ID", reqID)

		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapper, r)

		duration := time.Since(start)

		logEntry := map[string]interface{}{
			"level":       "info",
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
			"requestId":   reqID,
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      wrapper.statusCode,
			"durationMs":  duration.Milliseconds(),
			"bytes":       wrapper.bytesWritten,
			"remoteIp":    r.RemoteAddr,
			"userAgent":   r.UserAgent(),
		}

		if bytes, err := json.Marshal(logEntry); err == nil {
			println(string(bytes))
		}
	})
}
