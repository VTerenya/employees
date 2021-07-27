package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Entry *logrus.Logger
}

func NewLogger() *Logger {
	return &Logger{Entry: logrus.New()}
}

func (l Logger) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.Entry.WithFields(logrus.Fields{
			"type":        "access log",
			"method":      r.Method,
			"remote_addr": r.RemoteAddr,
			"host":        r.Host,
			"path":        r.URL.Path,
		}).Info()
		next.ServeHTTP(w, r)
	})
}

func (l Logger) TimeLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		l.Entry.WithFields(logrus.Fields{
			"type":      "time log",
			"work_time": duration,
		}).Info()
	})
}

type correlationKey string

var CorrelationID correlationKey = "ID"

func (l Logger) IDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id := uuid.New()
		ctx = context.WithValue(ctx, CorrelationID, id.String())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (l Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Entry.WithFields(fields)
}
