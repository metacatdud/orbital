package orbital

import (
	"net/http"
	"orbital/pkg/logger"
)

func PanicRecoverMiddleware() Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					http.Error(w, "internal error", 500)
				}
			}()
			next(w, r)
		}
	}
}

func LoggerMiddleware(log *logger.Logger) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Info("Route", "method", r.Method, "path", r.URL.Path)
			next(w, r)
		}
	}
}
