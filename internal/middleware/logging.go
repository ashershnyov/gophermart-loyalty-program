package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type responseData struct {
	status int
	size   int
	resp   []byte
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	r.responseData.resp = b
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// Logging is a middleware to log uri, HTTP method, handling duration, reponse status and response size.
func Logging(logger *slog.Logger) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			uri := r.RequestURI
			method := r.Method

			lrw := &loggingResponseWriter{
				w,
				&responseData{},
			}
			h.ServeHTTP(lrw, r)

			dur := time.Since(start)

			errorMsg := ""

			var logLevel slog.Level
			switch lrw.responseData.status {
			case http.StatusOK, http.StatusCreated, http.StatusAccepted,
				http.StatusNonAuthoritativeInfo, http.StatusNoContent, http.StatusResetContent,
				http.StatusPartialContent, http.StatusMultiStatus, http.StatusAlreadyReported,
				http.StatusIMUsed:

				logLevel = slog.LevelInfo
			default:
				errorMsg = strconv.QuoteToASCII(string(lrw.responseData.resp))
				logLevel = slog.LevelWarn
			}

			logger.LogAttrs(
				context.Background(),
				logLevel,
				"incoming request",
				slog.String("uri", uri),
				slog.String("method", method),
				slog.Duration("duration", dur),
				slog.Int("status_code", lrw.responseData.status),
				slog.Int("response_size", lrw.responseData.size),
				slog.String("error_msg", errorMsg),
			)
		})
	}
}
