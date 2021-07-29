package models

import (
	"context"
	"log"
	"time"

	"github.com/Alexseij/server/utils"
)

type Order struct {
	Description string    `json:"description"`
	Name        string    `json:"name"`
	From        string    `json:"from"`
	Destination string    `json:"destination"`
	TimeCreate  time.Time `json:"time_create"`
	TimeUpdate  time.Time `json:"time_update"`

	stop chan struct{}
}

func stop(stop chan struct{}) {
	stop <- struct{}{}
}

func (o *Order) MakeOrder() map[string]interface{} {
	ordersCollection := GetDB().Collection("orders")

	o.TimeCreate = time.Now()

	result, err := ordersCollection.InsertOne(context.TODO(), o)
	if err != nil {
		return utils.Message(false, "Invalid request to database")
	}

	log.Print("Order inserted in : ", o.TimeCreate, "with id : ", result.InsertedID)
	return utils.Message(true, "Order added")
}

func DeleteOrder(o) {

}

func (o *Order) deleteOrderWithDelay(delay time.Duration) map[string]interface{} {

	done := make(chan struct{})

	go func() {
		time.Sleep(delay)
		done <- struct{}{}
	}()

	select {
	case <-done:
		ordersCollection := GetDB().Collection("orders")
		_, err := ordersCollection.DeleteOne(context.TODO(), o)
		if err != nil {
			log.Print(err)
			return utils.Message(false, "Invalid delete document")
		}
		return utils.Message(true, "Document deleted")
	case <-o.stop:
		return utils.Message(true, "Document deleted")
	}
}
