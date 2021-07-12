package main

import (
	"github.com/VTerenya/employees/iternal/app/apiserver"
	"log"
)

func main() {
	server := apiserver.New()
	if err := server.Start(); err !=nil{
		log.Fatal(err)
	}

}
