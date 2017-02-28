package webserver

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/vladryk/server/server/handlers"
)



func GetServer(addr string) *http.Server {
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.NotFoundHandler)
	r.HandleFunc("/analyze", handlers.AnalyzeHandler)
	srv := &http.Server{
		Handler: r,
		Addr:    addr,
	}
	return srv
}
