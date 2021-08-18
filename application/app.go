package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Alexseij/server/handlers/auth"
	"github.com/Alexseij/server/handlers/order"
	"github.com/Alexseij/server/handlers/users"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ReqHandlerFunc func(*mongo.Database, http.ResponseWriter, *http.Request)

type App struct {
	Router *mux.Router
	DB     *mongo.Database
}

func (a *App) Init(dbUser, dbPassword, dbHost, dbName string) {

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

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal("Server dosent connect to database")
	}

	log.Print("Server connectrd to database")

	a.DB = client.Database(dbName)
	log.Print("Current database name : ", a.DB.Name())

	a.Router = mux.NewRouter()
	a.setHandlers()
}

func (a *App) setHandlers() {
	a.Post("/api/orderMake", a.handleReq(order.MakeOrder))
	a.Delete("/api/order/delete/{orderID}", a.handleReq(order.DeleteOrderWithID))
	a.Put("/api/orderUpdate", a.handleReq(order.UpdateOrder))
	a.Get("/api/user/login/{token}", a.handleReq(auth.LoginUser))
	a.Get("/api/order/{orderID}", a.handleReq(order.GetOrder))
	a.Get("/api/user/getCurrentRating/{email}", a.handleReq(users.GetCurrentRating))
	a.Get("/api/GetOrdersForCurrentProvider/{email}", a.handleReq(order.GetOrdersForCurrentProvider))
	a.Get("/api/GetOrdersForCurrentConsumer/{email}", a.handleReq(order.GetOrdersForCurrentConsumer))
	a.Put("/api/user/{email}/rating/{rating}", a.handleReq(users.UpdateRating))
}

func (a *App) Get(path string, handler func(http.ResponseWriter, *http.Request)) {
	a.Router.HandleFunc(path, handler).Methods("GET")
}

func (a *App) Put(path string, handler func(http.ResponseWriter, *http.Request)) {
	a.Router.HandleFunc(path, handler).Methods("PUT")
}

func (a *App) Delete(path string, handler func(http.ResponseWriter, *http.Request)) {
	a.Router.HandleFunc(path, handler).Methods("DELETE")
}

func (a *App) Post(path string, handler func(http.ResponseWriter, *http.Request)) {
	a.Router.HandleFunc(path, handler).Methods("POST")
}

func (a *App) handleReq(handler ReqHandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		handler(a.DB, rw, r)
	}
}
