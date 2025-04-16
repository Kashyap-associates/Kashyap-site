package server

import (
	"compress/gzip"
	"log/slog"
	"net/http"
	"strings"
)

// a middleware type for all middlewares
type Middleware func(http.Handler) http.Handler

// logger type 
type loggerResponse struct {
	http.ResponseWriter
	statusCode int
}

// gzip type
type gzipResponse struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

// function to couple all middlewares
func createStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, x := range xs {
			next = x(next)
		}
		return next
	}
}

// logging middleware function to log each request
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &loggerResponse{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.statusCode),
		)
	})
}

// function below will write the response code 
func (rec *loggerResponse) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// gzip middleware function to compress data and send it
func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzrw := gzipResponse{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzrw, r)
	})
}

// function below will return the compressed data
func (w gzipResponse) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// function to handle errors
func handleError(w http.ResponseWriter, r *http.Request, err error, toggle bool) {
	slog.Error(err.Error())
	if toggle {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/error", http.StatusFound)
}
