package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alexseij/server/application"
	"github.com/joho/godotenv"
)

var (
	Version = "1.0.0"
)

func main() {

	application := &application.App{}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("file main.go , godotenv.Load() : ", err)
	}

	dbUser := os.Getenv("db_user")
	dbPassword := os.Getenv("db_password")
	dbHost := os.Getenv("db_host")
	dbName := os.Getenv("db_name")

	application.Init(dbUser, dbPassword, dbHost, dbName)

	server := &http.Server{
		Addr:         "localhost:8000",
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

	server.Shutdown(ctx)

	log.Print("server shutdown")
	os.Exit(0)
}
