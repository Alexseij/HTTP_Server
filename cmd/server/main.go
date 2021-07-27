package main

import (
	"net/http"
	"os"

	"github.com/Alexseij/server/handlers"
	"github.com/gorilla/mux"
)

var (
	Version = "1.0.0"
)

func main() {

	server := &http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: buildHandler(),
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		os.Exit(-1)
	}

}

func buildHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/login", handlers.LoginUser).Methods("POST")
	return router
}
