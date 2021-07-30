package models

import (
	"context"
	"log"
	"time"

	"github.com/Alexseij/server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderID struct {
	OrderID string `json:"order_id"`
}

type Order struct {
	Description string    `json:"description"`
	Name        string    `json:"name"`
	From        string    `json:"from"`
	Destination string    `json:"destination"`
	TimeCreate  time.Time `json:"time_create"`
	TimeUpdate  time.Time `json:"time_update"`

	insertedID string
}

type OrderForUpdate struct {
	Order
	OrderID
}

func (o *Order) MakeOrder() map[string]interface{} {
	ordersCollection := GetDB().Collection("orders")

	o.TimeCreate = time.Now()

	result, err := ordersCollection.InsertOne(context.TODO(), o)
	if err != nil {
		return utils.Message(false, "Invalid request to database")
	}
	o.insertedID = result.InsertedID.(string)

	log.Print("Order inserted in : ", o.TimeCreate, "with id : ", result.InsertedID)

	resp := utils.Message(true, "Order added")
	resp["id"] = result.InsertedID

	return resp
}

func DeleteOrder(order *Order) map[string]interface{} {
	ordersCollection := GetDB().Collection("orders")
	delResult, err := ordersCollection.DeleteOne(context.TODO(), bson.M{"_id": order.insertedID})
	if err != nil {
		log.Print(err)
		return utils.Message(false, "Invalid operation to delete order")
	}
	log.Print("Delete in : ", time.Now(), " count of elements : ", delResult.DeletedCount)

	return utils.Message(true, "Order deleted")
}

func FindOrder(orderID string) (o *Order, err error) {
	ordersCollection := GetDB().Collection("orders")

	var order *Order

	if err := ordersCollection.FindOne(context.TODO(), bson.M{"_id": orderID}).Decode(order); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Print("Order already deleted")
			return nil, nil
		}
		log.Print("Error conntection with database")
		return nil, err
	}
	return order, nil
}

func UpdateOrder(updateOrder *OrderForUpdate) map[string]interface{} {
	ordersCollection := GetDB().Collection("orders")

	_, err := ordersCollection.UpdateOne(context.TODO(), bson.M{"_id": updateOrder.OrderID}, updateOrder.Order)
	if err != nil {
		log.Print(err)
		return utils.Message(false, "Ivalid request to database")
	}

	return utils.Message(true, "Item updated !")
}
