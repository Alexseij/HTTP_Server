package order

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Alexseij/server/models"
	"github.com/Alexseij/server/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MakeOrder(rw http.ResponseWriter, r *http.Request) {
	order := &models.Order{}

	if err := json.NewDecoder(r.Body).Decode(order); err != nil {
		utils.Respond(rw, utils.Message(false, "Invalid request"))
		return
	}

	resp := order.MakeOrder()
	utils.Respond(rw, resp)

	if resp["status"].(bool) {
		go deleteOrderAfterDelay(time.Second*30, resp["id"].(primitive.ObjectID), rw)
	}

}

func deleteOrderAfterDelay(delay time.Duration, orderID primitive.ObjectID, rw http.ResponseWriter) {
	done := make(chan struct{})
	go func() {
		time.Sleep(delay)
		done <- struct{}{}
	}()

	select {
	case <-done:
		deleteOrder(rw, orderID)
	}
}

func deleteOrder(rw http.ResponseWriter, orderID primitive.ObjectID) {
	order, err := models.FindOrder(orderID)
	if err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Invalid delete from database"))
		return
	}

	if order == nil {
		utils.Respond(rw, utils.Message(false, "Order alreay deleted"))
		return
	}

	resp := models.DeleteOrder(order)
	utils.Respond(rw, resp)

}

func DeleteOrderWithID(rw http.ResponseWriter, r *http.Request) {
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

	deleteOrder(rw, primitiveOrderID)
}

func UpdateOrder(rw http.ResponseWriter, r *http.Request) {
	updateOrder := &models.Order{}

	if err := json.NewDecoder(r.Body).Decode(updateOrder); err != nil {
		log.Print(err)
		utils.Respond(rw, utils.Message(false, "Ivalid request"))
		return
	}

	resp := models.UpdateOrder(updateOrder)

	utils.Respond(rw, resp)
}
