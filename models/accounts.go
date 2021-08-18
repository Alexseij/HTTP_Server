package models

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Alexseij/server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/idtoken"
)

type User struct {
	Token          string `json:"token" bson:"token"`
	ConsumerRating int    `json:"consumer_rating" bson:"consumer_rating"`
	ProviderRating int    `json:"provider_rating" bson:"provider_rating"`
	Name           string `json:"name" bson:"name"`
	Email          string `json:"email" bson:"email"`
	NowProvider    bool   `json:"now_provider" bson:"now_provider"`
}

const (
	DefaultRating = 5
)

func validateToken(token string) (*idtoken.Payload, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	payload, err := idtoken.Validate(ctx, token, os.Getenv("client_id"))
	if err != nil {
		return nil, err
	}
	return payload, nil
}

func GetUser(db *mongo.Database, email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	usersCollection := db.Collection("users")

	user := &User{}

	if err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *User) createUser(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	usersCollection := db.Collection("users")

	result, err := usersCollection.InsertOne(ctx, u)
	if err != nil {
		return err
	}
	log.Print("User created with id : ", result.InsertedID)

	return nil
}

func LoginUser(db *mongo.Database, token string) map[string]interface{} {
	payload, err := validateToken(token)
	if err != nil {
		log.Print(err)
		resp := utils.Message(false, "Error with validate token")
		resp["user"] = nil
		return resp
	}

	if ok := utils.CheckDomain(payload.Claims["email"].(string)); !ok {
		resp := utils.Message(false, "Current email domain , can't use in this system")
		resp["user"] = nil
		return resp
	}

	user, err := GetUser(db, payload.Claims["email"].(string))
	if err != nil {
		log.Print(err)
		resp := utils.Message(false, "Error with getting user from collection")
		resp["user"] = nil
		return resp
	}
	if user == nil {
		log.Print("in condition user == nil , user : ", user)
		user = &User{}
		user.Email = payload.Claims["email"].(string)
		user.Name = payload.Claims["name"].(string)
		user.ConsumerRating = DefaultRating
		user.ProviderRating = DefaultRating
		user.NowProvider = false
		user.Token = token

		if err := user.createUser(db); err != nil {
			log.Print(err)
			resp := utils.Message(false, "Error with create user")
			resp["user"] = nil
			return resp
		}
	}
	resp := utils.Message(true, "User sing in")
	resp["user"] = user
	return resp
}

func (u *User) UpdateRating(db *mongo.Database, currentRating int, topicName, projectID string) map[string]interface{} {
	usersCollection := db.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := usersCollection.UpdateOne(ctx, bson.M{"token": u.Token}, bson.M{"$set": bson.M{"rating": currentRating}})
	if err != nil {
		log.Print("file accounts.go , UpdateRating() : ", err)
		return utils.Message(false, "Incorrect update query to database")
	}

	return utils.Message(true, "Successful update")

}
