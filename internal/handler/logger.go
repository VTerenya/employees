package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type ILogger interface {
	AccessLogMiddleware(next http.Handler) http.Handler
	TimeLogMiddleware(next http.Handler) http.Handler
	IDMiddleware(next http.Handler) http.Handler
	WithFields(fields logrus.Fields) *logrus.Entry
}
