package application

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Alexseij/server/handlers/auth"
	"github.com/Alexseij/server/handlers/order"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("file : db.go  , mongo.Connect() : ", err)
	}

	a.DB = client.Database(dbName)
	log.Print("Database created , Name : ", a.DB.Name())

	a.Router = mux.NewRouter()
	a.setHandlers()
}

func (a *App) setHandlers() {
	a.Post("/api/user/new", a.handleReq(auth.CreateUser))
	a.Post("/api/user/login", a.handleReq(auth.LoginUser))
	a.Post("/api/order/make", a.handleReq(order.MakeOrder))
	a.Delete("/api/order/delete", a.handleReq(order.DeleteOrderWithID))
	a.Put("/api/order/update", a.handleReq(order.UpdateOrder))
	a.Get("/api/order/get", a.handleReq(order.GetOrder))
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
