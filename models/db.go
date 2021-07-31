package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db *mongo.Database
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbUser := os.Getenv("db_user")
	dbPassword := os.Getenv("db_password")
	dbHost := os.Getenv("db_host")
	dbName := os.Getenv("db_name")

	URI := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority",
		dbUser,
		dbPassword,
		dbHost,
		dbName,
	)

	clientOptions := options.Client().ApplyURI(URI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("file : db.go  , mongo.Connect() : ", err)
	}
	db = client.Database(dbName)
	log.Print("Database created , Name : ", db.Name())
}

func GetDB() *mongo.Database {
	return db
}
