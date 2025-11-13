package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Инициализация логгера
var log = logrus.New()

func init() {
	// Устанавливаем JSON форматтер
	log.SetFormatter(&logrus.JSONFormatter{})
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, r)
		// log.Println(wrapper.StatusCode, r.Method, r.URL.Path, time.Since(start))
		log.WithFields(logrus.Fields{
			"status_code": wrapper.StatusCode,
			"method":      r.Method,
			"path":        r.URL.Path,
			"duration":    time.Since(start).String(),
			"duration_ms": time.Since(start).Milliseconds(),
		})

	})
}
