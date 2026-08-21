package httpapi

import (
	"log/slog"
	"net/http"
	"time"
)

func Logger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			st := time.Now()
			next.ServeHTTP(w, r)
			log.Info("request", "method", r.Method, "path", r.URL.Path, "duration_ms", time.Since(st).Milliseconds())
		})
	}
}
func Recover(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if v := recover(); v != nil {
					log.Error("panic", slog.Any("value", v))
					http.Error(w, "internal error", 500)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = time.Now().UTC().Format("20060102150405.000000000")
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
