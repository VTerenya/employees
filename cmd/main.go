package main

import (
	"log"
	"net/http"

	"github.com/VTerenya/employees/internal/handler"
	"github.com/VTerenya/employees/internal/middleware"
	"github.com/VTerenya/employees/internal/repository"
	"github.com/VTerenya/employees/internal/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler interface { // nolint: deadcode
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
	r := mux.NewRouter()
	myData := repository.NewDataBase()
	myRepo := repository.NewRepo(myData)
	myServ := service.NewServ(myRepo)
	myH := handler.NewHandler(myServ)
	r.HandleFunc("/positions", myH.GetPositions).Queries("limit", "{limit:\\S+}", "offset", "{offset:\\S+}")
	r.HandleFunc("/employees", myH.GetEmployees).Queries("limit", "{limit:\\S+}", "offset", "{offset:\\S+}")
	r.HandleFunc("/position/{id:\\S+}", myH.GetPosition).Methods("GET")
	r.HandleFunc("/employee/{id:\\S+}", myH.GetEmployee).Methods("GET")
	r.HandleFunc("/position/{id:\\S+}", myH.DeletePosition).Methods("DELETE")
	r.HandleFunc("/employee/{id:\\S+}", myH.DeleteEmployee).Methods("DELETE")
	r.HandleFunc("/position", myH.UpdatePosition).Methods("PUT")
	r.HandleFunc("/employee", myH.UpdateEmployee).Methods("PUT")
	r.HandleFunc("/position", myH.CreatePosition).Methods("POST")
	r.HandleFunc("/employee", myH.CreateEmployee).Methods("POST")
	r.Use(middleware.IDMiddleware, middleware.TimeLogMiddleware, middleware.AccessLogMiddleware)
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Run()
}
