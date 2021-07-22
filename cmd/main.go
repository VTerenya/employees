package main

import (
	"log"
	"net/http"

	"github.com/VTerenya/employees/internal/handler"
	"github.com/VTerenya/employees/internal/repository"
	"github.com/VTerenya/employees/internal/service"
)

func main() {
	mux := http.NewServeMux()
	myData := repository.NewDataBase()
	myRepo := repository.NewRepository(myData)
	myServ := service.NewService(myRepo)
	myH := handler.NewHandler(myServ)
	mux.Handle("/positions", myH)
	mux.Handle("/position/", myH)
	mux.Handle("/position", myH)
	mux.Handle("/employees", myH)
	mux.Handle("/employee/", myH)
	mux.Handle("/employee", myH)
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
