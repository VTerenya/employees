package main

import (
	"github.com/VTerenya/employees/internal/middleware"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"

	"github.com/VTerenya/employees/internal/handler"
	"github.com/VTerenya/employees/internal/repository"
	"github.com/VTerenya/employees/internal/service"
)

type Handler interface {
	GetPositions(w http.ResponseWriter, r *http.Request)
	GetEmployees(w http.ResponseWriter, r *http.Request)
	GetPosition(w http.ResponseWriter, r *http.Request)
	GetEmployee(w http.ResponseWriter, r *http.Request)
	CreatePosition(w http.ResponseWriter, r *http.Request)
	CreateEmployee(w http.ResponseWriter, r *http.Request)
	DeletePosition(w http.ResponseWriter, r *http.Request)
	DeleteEmployee(w http.ResponseWriter, r *http.Request)
	UpdatePosition(w http.ResponseWriter, r *http.Request)
	UpdateEmployee(w http.ResponseWriter, r *http.Request)
}

func Run() {
	contextLogger := logrus.WithFields(logrus.Fields{
		"logger": "LOGRUS",
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	middleware.Entry = contextLogger

	mux := http.NewServeMux()
	myData := repository.NewDataBase()
	myRepo := repository.NewRepo(myData)
	myServ := service.NewServ(myRepo)
	myH := handler.NewHandler(myServ)
	mux.Handle("/positions", myH)
	mux.Handle("/position/", myH)
	mux.Handle("/position", myH)
	mux.Handle("/employees", myH)
	mux.Handle("/employee/", myH)
	mux.Handle("/employee", myH)
	timeHandler := middleware.TimeLogMiddleware(mux)
	accessHandler := middleware.AccessLogMiddleware(timeHandler)
	idHandler := middleware.IdMiddleware(accessHandler)
	err := http.ListenAndServe("localhost:8080", idHandler)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Run()
}