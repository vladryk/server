package main

import (
	"log"
	"github.com/vladryk/webserver/webserver"
)

func main() {
	serv := webserver.GetServer("127.0.0.1:8000")
	log.Fatal(serv.ListenAndServe())
}
