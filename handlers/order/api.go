package order

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Alexseij/server/models"
	"github.com/Alexseij/server/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func MakeOrder(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	order := &models.Order{}

	if err := json.NewDecoder(r.Body).Decode(order); err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Invalid request"))
		return
	}

	resp := order.MakeOrder(db)
	utils.Respond(rw, resp)

	if resp["status"].(bool) {
		go deleteOrderAfterDelay(db, time.Minute*2, resp["id"].(primitive.ObjectID), rw)
	}

}

func deleteOrderAfterDelay(db *mongo.Database, delay time.Duration, orderID primitive.ObjectID, rw http.ResponseWriter) {
	time.Sleep(delay)
	deleteOrder(db, rw, orderID)
}

func deleteOrder(db *mongo.Database, rw http.ResponseWriter, orderID primitive.ObjectID) {
	order, err := models.GetOrder(db, orderID)
	if err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Invalid delete from database"))
		return
	}
	if order.Status {
		utils.Respond(rw, utils.Message(false, "Cant delete order"))
		return
	}
	if order == nil {
		utils.Respond(rw, utils.Message(false, "Order alreay deleted"))
		return
	}

	resp := order.DeleteOrder(db)

	utils.Respond(rw, resp)
}

func DeleteOrderWithID(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	primitiveOrderID, err := primitive.ObjectIDFromHex(vars["orderID"])
	if err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Invalid request"))
		return
	}

	deleteOrder(db, rw, primitiveOrderID)
}

func UpdateOrder(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	updateOrder := &models.Order{}

	if err := json.NewDecoder(r.Body).Decode(updateOrder); err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Ivalid request"))
		return
	}

	currentOrder, err := models.GetOrder(db, updateOrder.OrderID)
	if err != nil {
		log.Print("file : order/api.go , UpdateOrder() : ", err)
		utils.Respond(rw, utils.Message(false, "Invalid getting order"))
		return
	}

	resp := currentOrder.UpdateOrder(db, updateOrder)

	utils.Respond(rw, resp)
}

func GetOrders(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	orders, err := models.GetOrders(db)
	if err != nil {
		log.Print("file : orders/api.go , GetOrders() : ", err)
		utils.Respond(rw, utils.Message(false, "Invalid request"))
		return
	}

	resp := utils.Message(true, "Success")
	resp["orders"] = orders
	utils.Respond(rw, resp)
}

func GetOrdersForCurrentConsumer(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	orders, err := models.GetOrdersForCurrentConsumer(db, vars["email"])
	if err != nil {
		log.Print("file : orders/api.go , GetOrdersForCurrenConsumer() : ", err)
		utils.Respond(rw, utils.Message(false, "Ivalid request"))
		return
	}

	resp := utils.Message(true, "Success")
	resp["orders"] = orders
	utils.Respond(rw, resp)
}

func GetOrdersForCurrentProvider(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	orders, err := models.GetOrdersForCurrentProvider(db, vars["email"])
	if err != nil {
		log.Print("file : orders/api.go , GetOrdersForCurrentProvider() : ", err)
		utils.Respond(rw, utils.Message(false, "Invalid request"))
		return
	}

	resp := utils.Message(true, "Success")
	resp["orders"] = orders
	utils.Respond(rw, resp)
}

func GetOrder(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	primitiveOrderID, err := primitive.ObjectIDFromHex(vars["orderID"])
	if err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Fail with creating primitive orderID"))
		return
	}

	order, err := models.GetOrder(db, primitiveOrderID)
	if err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Fail with getting order from db"))
		return
	}

	if order == nil {
		utils.Respond(rw, utils.Message(false, "Order with curren id dosent exits"))
		return
	}

	resp := utils.Message(true, "Order was found")
	resp["order"] = order

	utils.Respond(rw, resp)
}
