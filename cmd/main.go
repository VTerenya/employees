package main

import (
	"log"
	"net/http"

	"github.com/VTerenya/employees/internal/handler"
	"github.com/VTerenya/employees/internal/middleware"
	"github.com/VTerenya/employees/internal/repository"
	"github.com/VTerenya/employees/internal/service"
	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()
	logger := middleware.NewLogger()
	myData := repository.NewDataBase()
	myRepo := repository.NewRepo(myData)
	myServ := service.NewServ(myRepo)
	myH := handler.NewHandler(myServ, logger)
	r.HandleFunc("/positions", myH.GetPositions).Queries("limit", "{limit:\\S+}", "offset", "{offset:\\S+}")
	r.HandleFunc("/employees", myH.GetEmployees).Queries("limit", "{limit:\\S+}", "offset", "{offset:\\S+}")
	r.HandleFunc("/position/{id:\\S+}", myH.GetPosition).Methods("GET")
	r.HandleFunc("/employee/{id:\\S+}", myH.GetEmployee).Methods("GET")
	r.HandleFunc("/position/{id:\\S+}", myH.DeletePosition).Methods("DELETE")
	r.HandleFunc("/employee/{id:\\S+}", myH.DeleteEmployee).Methods("DELETE")
	r.HandleFunc("/position/{id:\\S+}", myH.UpdatePosition).Methods("UPDATE")
	r.HandleFunc("/employee/{id:\\S+}", myH.UpdateEmployee).Methods("UPDATE")
	r.HandleFunc("/position", myH.CreatePosition).Methods("POST")
	r.HandleFunc("/employee", myH.CreateEmployee).Methods("POST")
	r.Use(logger.IDMiddleware, logger.TimeLogMiddleware, logger.AccessLogMiddleware)
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Run()
}
