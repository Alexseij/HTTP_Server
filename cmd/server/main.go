package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Alexseij/handlers/auth"
	"github.com/Alexseij/handlers/order"
	"github.com/Alexseij/server/config"
	"github.com/gorilla/mux"
)

var (
	Version = "1.0.0"
)

var flagConfig = flag.String("config", "../config/local.yml", "path to cofig file ")

func main() {

	config, err := config.LoadCfg(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	serverAddr := fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: buildHandler(),
	}
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		os.Exit(-1)
	}
}

func buildHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/new", auth.CreateUser).Methods("POST")
	router.HandleFunc("/api/user/login", auth.LoginUser).Methods("POST")
	router.HandleFunc("/api/order/make", order.MakeOrder).Methods("POST")

	return router
}
