package main

import (
	"log"
	"github.com/vladryk/server/server"
)

func main() {
	serv := webserver.GetServer("0.0.0.0:8000")
	log.Fatal(serv.ListenAndServe())
}
