package models

import (
	"context"
	"log"
	"time"

	"github.com/Alexseij/server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderID struct {
	ID string `json:"id" bson:"_id"`
}

type Order struct {
	OrderID     primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Description string             `json:"description" bson:"description"`
	Name        string             `json:"name" bson:"name"`
	From        string             `json:"from" bson:"from"`
	Destination string             `json:"destination" bson:"destination"`
	TimeCreate  primitive.DateTime `bson:"time_create,omitempty"`
	TimeUpdate  primitive.DateTime `bson:"time_update,omitempty"`
}

func (o *Order) MakeOrder() map[string]interface{} {
	ordersCollection := GetDB().Collection("orders")
	ctx := context.TODO()

	o.TimeCreate = primitive.NewDateTimeFromTime(time.Now())
	o.TimeUpdate = primitive.NewDateTimeFromTime(time.Now())
	o.OrderID = primitive.NewObjectID()

	result, err := ordersCollection.InsertOne(ctx, o)
	if err != nil {
		log.Print("file : orders.go , InserOne() : ", err)
		return utils.Message(false, "Invalid request to database")
	}

	log.Print("Order inserted in : ", o.TimeCreate, " with id : ", result.InsertedID)

	resp := utils.Message(true, "Order added")
	resp["id"] = result.InsertedID

	return resp
}

func DeleteOrder(order *Order) map[string]interface{} {
	ordersCollection := GetDB().Collection("orders")
	ctx := context.TODO()

	delResult, err := ordersCollection.DeleteOne(ctx, bson.M{"_id": order.OrderID})
	if err != nil {
		log.Print("file : orders.go , DeleteOne() : ", err)
		return utils.Message(false, "Invalid operation to delete order")
	}
	log.Print("Delete in : ", time.Now(), " count of elements : ", delResult.DeletedCount)

	return utils.Message(true, "Order deleted")
}

func FindOrder(orderID primitive.ObjectID) (o *Order, err error) {
	ordersCollection := GetDB().Collection("orders")
	ctx := context.TODO()
	order := &Order{}

	if err := ordersCollection.FindOne(ctx, bson.M{"_id": orderID}).Decode(order); err != nil {
		if err == mongo.ErrNoDocuments {
			log.Print("Order already deleted")
			return nil, nil
		}
		log.Print("file : orders.go , FindOne() : ", err)
		return nil, err
	}
	return order, nil
}

func UpdateOrder(updateOrder *Order) map[string]interface{} {
	ordersCollection := GetDB().Collection("orders")
	ctx := context.TODO()

	updateOrder.TimeUpdate = primitive.NewDateTimeFromTime(time.Now())

	_, err := ordersCollection.UpdateOne(ctx, bson.M{"_id": updateOrder.OrderID}, bson.M{"$set": updateOrder})
	if err != nil {
		log.Print("file : orders.go , UpdateOne() : ", err)
		return utils.Message(false, "Ivalid request to database")
	}

	return utils.Message(true, "Item updated !")
}
