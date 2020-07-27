package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func Loging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			logrus.WithFields(logrus.Fields{
				"at": time.Now().Format("2006-01-02 15:04:05"),
			}).Info("incoming request\n")

		}

		logrus.WithFields(logrus.Fields{
			"at":     time.Now().Format("2006-01-02 15:04:05"),
			"method": r.Method,
			"uri":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Info("incoming request\n")

		next.ServeHTTP(w, r)
	})
}
