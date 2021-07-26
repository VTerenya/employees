package middleware

import (
	"context"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var Entry *logrus.Entry

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		Entry.WithFields(logrus.Fields{
			"type":        "access log",
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"host":        r.Host,
		}).Info(r.URL.Path)
	})
}

func TimeLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		start := time.Now()
		Entry.WithFields(logrus.Fields{
			"type":      "time log",
			"work_time": time.Since(start),
		}).Info(r.URL.Path)
	})
}

type contextKey string

var contextID contextKey = "ID"

func IdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New()
		ctx := context.WithValue(r.Context(), contextID, id.String())
		r = r.WithContext(ctx)
		Entry.WithFields(logrus.Fields{
			"type":           "id log",
			"correlation_id": id.String(),
		}).Info(r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
