package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alexseij/server/application"
	"github.com/Alexseij/server/config"
	"github.com/joho/godotenv"
)

var (
	Version = "1.0.0"
)

var flagCfg = flag.String("config", "config/local.yml", "Config for starting server")

func main() {

	flag.Parse()

	application := &application.App{}

	cfg, err := config.LoadCfg(*flagCfg)
	if err != nil {
		log.Fatal("error to connect with cfg : ", err)
	}

	if err := godotenv.Load(); err != nil {
		log.Fatal("error with loading .env file : ", err)
	}

	application.Init(cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBName)

	port := os.Getenv("PORT")
	var addr string
	if port == "" {
		addr = ":8080"
	}
	addr = fmt.Sprintf(":%s", port)

	server := &http.Server{
		Addr:         addr,
		Handler:      application.Router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Print(err)
			os.Exit(-1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)

	defer func() {
		if err := application.DB.Client().Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Print("server shutdown")
	os.Exit(0)
}
