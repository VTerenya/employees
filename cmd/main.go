package main

import (
	"github.com/VTerenya/employees/internal/handler"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	myH := handler.NewHandler()
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
