package models

import (
	"context"
	"log"
	"os"

	"github.com/Alexseij/server/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/idtoken"
)

type Token struct {
	Token string `json:"token"`
}

type User struct {
	Token  string `json:"token"`
	Rating int    `json:"rating"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

const (
	DefaultRating = 5
)

var (
	clientID = os.Getenv("client_id")
)

func validateToken(token string) (*idtoken.Payload, error) {

	tokenValidator, err := idtoken.NewValidator(context.Background())
	if err != nil {
		log.Print(err)
		return nil, err
	}
	payload, err := tokenValidator.Validate(context.Background(), token, clientID)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return payload, nil
}

func (u *User) Validate() (*idtoken.Payload, map[string]interface{}, bool) {

	payload, err := validateToken(u.Token)
	if err != nil {
		log.Print(err)
		return nil, utils.Message(false, "Incorrect token"), false
	}

	return payload, utils.Message(true, "Accepted"), true

}

func (u *User) Create() map[string]interface{} {
	payload, resp, ok := u.Validate()
	if !ok {
		return resp
	}

	users := GetDB().Collection("users")
	ctx := context.TODO()

	user, err := GetUser(u.Token)
	if err != nil {
		log.Print(err)
		return utils.Message(false, "Incorrect database query")
	}

	if user != nil {
		return utils.Message(false, "User already exist")
	}

	u.Email = payload.Claims["email"].(string)
	u.Name = payload.Claims["name"].(string)
	u.Rating = DefaultRating

	result, err := users.InsertOne(ctx, u)
	if err != nil {
		log.Print(err)
		return utils.Message(false, "Ivalid request to database")
	}

	log.Print("User created with id :", result.InsertedID)

	return utils.Message(true, "User created !")
}

func GetUser(token string) (*User, error) {
	usersCollection := GetDB().Collection("users")
	ctx := context.TODO()

	user := &User{}

	if err := usersCollection.FindOne(ctx, bson.M{"token": token}).Decode(user); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		log.Print(err)
		return nil, err
	}

	return user, nil
}

func LoginUser(token string) map[string]interface{} {
	user, err := GetUser(token)
	if err != nil {
		log.Print(err)
		return utils.Message(false, "Incorrect database query")
	}

	if user == nil {
		return utils.Message(false, "User dosent exist")
	}

	log.Print("User : ", user.Token, "Logined.")

	return utils.Message(true, "User login into account")

}
