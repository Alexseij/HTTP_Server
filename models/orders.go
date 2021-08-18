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

type Order struct {
	OrderID     primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Description string             `json:"description" bson:"description"`
	Name        string             `json:"name" bson:"name"`
	From        string             `json:"from" bson:"from"`
	Provider    string             `json:"provider" bson:"provider"`
	Destination string             `json:"destination" bson:"destination"`
	Status      bool               `json:"status" bson:"status"`
	TimeCreate  primitive.DateTime `bson:"time_create,omitempty"`
	TimeUpdate  primitive.DateTime `bson:"time_update,omitempty"`
}

func (o *Order) MakeOrder(db *mongo.Database) map[string]interface{} {
	ordersCollection := db.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

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

func (o *Order) DeleteOrder(db *mongo.Database) map[string]interface{} {
	ordersCollection := db.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	delResult, err := ordersCollection.DeleteOne(ctx, bson.M{"_id": o.OrderID})
	if err != nil {
		log.Print("file : orders.go , DeleteOne() : ", err)
		return utils.Message(false, "Invalid operation to delete order")
	}
	log.Print("Delete in : ", time.Now(), " count of elements : ", delResult.DeletedCount)

	return utils.Message(true, "Order deleted")
}

func GetOrder(db *mongo.Database, orderID primitive.ObjectID) (o *Order, err error) {
	ordersCollection := db.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

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

func GetOrders(db *mongo.Database) ([]*Order, error) {
	ordersCollection := db.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var orders []*Order

	cursor, err := ordersCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func GetOrdersForCurrentConsumer(db *mongo.Database, email string) ([]*Order, error) {
	ordersCollection := db.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var orders []*Order

	cursor, err := ordersCollection.Find(ctx, bson.M{"status": false, "from": email})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func GetOrdersForCurrentProvider(db *mongo.Database, email string) ([]*Order, error) {
	ordersCollection := db.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var orders []*Order

	cursor, err := ordersCollection.Find(ctx, bson.M{"from": bson.M{"$ne": email}, "status": false})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (o *Order) UpdateOrder(db *mongo.Database, updatedOrder *Order) map[string]interface{} {
	ordersCollection := db.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	o.TimeUpdate = primitive.NewDateTimeFromTime(time.Now())

	_, err := ordersCollection.UpdateOne(ctx, bson.M{"_id": o.OrderID}, bson.M{"$set": updatedOrder})
	if err != nil {
		log.Print("file : orders.go , UpdateOne() : ", err)
		return utils.Message(false, "Ivalid request to database")
	}

	return utils.Message(true, "Item updated !")
}
