package main

import (
	"github.com/VTerenya/employees/internal/database"
	"github.com/VTerenya/employees/internal/handler"
	"github.com/VTerenya/employees/internal/repository"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	myData := database.NewDatabase()
	myRepo := repository.NewRepository(myData)
	myH := handler.NewHandler(myRepo)
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
