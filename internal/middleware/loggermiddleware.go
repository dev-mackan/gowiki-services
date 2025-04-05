package middleware

import (
	"log"
	"net/http"
)

type WrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *WrappedResponseWriter) WriteHeader(status int) {
	w.statusCode = status
	w.ResponseWriter.WriteHeader(status)
}

func NewLoggerMiddleware(entityName string) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s <-- %s %s %+v %s", entityName, r.Method, r.RemoteAddr, r.Header.Values("referer"), r.URL.Path)
			writer := &WrappedResponseWriter{ResponseWriter: w, statusCode: 200}
			next.ServeHTTP(writer, r)
			log.Printf("%s --> %d %s %s", entityName, writer.statusCode, r.RemoteAddr, r.URL.Path)
		})
	}
}
