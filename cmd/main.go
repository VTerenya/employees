package main

import (
	"log"
	"net/http"

	"github.com/VTerenya/employees/internal/handler"
	"github.com/VTerenya/employees/internal/middleware"
	"github.com/VTerenya/employees/internal/repository"
	"github.com/VTerenya/employees/internal/service"
	"github.com/sirupsen/logrus"
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
	logrus.SetFormatter(&logrus.JSONFormatter{})
	mux := http.NewServeMux()
	logger := middleware.NewLogger()
	myData := repository.NewDataBase()
	myRepo := repository.NewRepo(myData)
	myServ := service.NewServ(myRepo)
	myH := handler.NewHandler(myServ, logger)
	mux.Handle("/positions", myH)
	mux.Handle("/position/", myH)
	mux.Handle("/position", myH)
	mux.Handle("/employees", myH)
	mux.Handle("/employee/", myH)
	mux.Handle("/employee", myH)
	idMiddleware := logger.IDMiddleware(mux)
	timerMiddleware := logger.TimeLogMiddleware(idMiddleware)
	accessLogMiddleware := logger.AccessLogMiddleware(timerMiddleware)
	err := http.ListenAndServe("localhost:8080", accessLogMiddleware)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Run()
}
