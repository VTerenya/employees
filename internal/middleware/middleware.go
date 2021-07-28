package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		logrus.WithFields(logrus.Fields{
			"type":        "access log",
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"host":        r.Host,
		}).Info()
	})
}

func TimeLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		logrus.WithFields(logrus.Fields{
			"type":      "time log",
			"work_time": duration,
		}).Info()
	})
}

const CorrelationID = "ID"

func IDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New()
		ctx = context.WithValue(ctx, CorrelationID, id.String())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
