package order

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Alexseij/server/models"
	"github.com/Alexseij/server/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func MakeOrder(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	order := &models.Order{}

	if err := json.NewDecoder(r.Body).Decode(order); err != nil {
		utils.Respond(rw, utils.Message(false, "Invalid request"))
		return
	}

	resp := order.MakeOrder(db)
	utils.Respond(rw, resp)

	if resp["status"].(bool) {
		go deleteOrderAfterDelay(db, time.Second*30, resp["id"].(primitive.ObjectID), rw)
	}

}

func deleteOrderAfterDelay(db *mongo.Database, delay time.Duration, orderID primitive.ObjectID, rw http.ResponseWriter) {
	done := make(chan struct{})
	go func() {
		time.Sleep(delay)
		done <- struct{}{}
	}()

	select {
	case <-done:
		deleteOrder(db, rw, orderID)
	}
}

func deleteOrder(db *mongo.Database, rw http.ResponseWriter, orderID primitive.ObjectID) {
	order, err := models.FindOrder(db, orderID)
	if err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Invalid delete from database"))
		return
	}

	if order == nil {
		utils.Respond(rw, utils.Message(false, "Order alreay deleted"))
		return
	}

	resp := models.DeleteOrder(db, order)
	utils.Respond(rw, resp)

}

func DeleteOrderWithID(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	orderID := &models.OrderID{}

	if err := json.NewDecoder(r.Body).Decode(orderID); err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Invalid request"))
		return
	}

	primitiveOrderID, err := primitive.ObjectIDFromHex(orderID.ID)
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

	resp := models.UpdateOrder(db, updateOrder)

	utils.Respond(rw, resp)
}

func GetOrder(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {

	orderID := &models.OrderID{}

	if err := json.NewDecoder(r.Body).Decode(orderID); err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Ivalid rquest"))
		return
	}

	primitiveOrderID, err := primitive.ObjectIDFromHex(orderID.ID)
	if err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Fail with creating primitive orderID"))
		return
	}

	order, err := models.FindOrder(db, primitiveOrderID)
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
	resp["id"] = orderID.ID
	resp["description"] = order.Description
	resp["name"] = order.Name
	resp["from"] = order.From
	resp["destination"] = order.Destination
	resp["time_create"] = order.TimeCreate
	resp["time_update"] = order.TimeUpdate

	utils.Respond(rw, resp)
}
