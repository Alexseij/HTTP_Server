package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Alexseij/server/models"
	"github.com/Alexseij/server/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	token := &models.Token{}

	err := json.NewDecoder(r.Body).Decode(token)
	if err != nil {
		utils.Respond(rw, utils.Message(false, "Ivalid request"))
		return
	}

	user := &models.User{
		Token: token.Token,
	}

	resp := user.Create(db)
	utils.Respond(rw, resp)
}

func LoginUser(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	token := &models.Token{}

	err := json.NewDecoder(r.Body).Decode(token)
	if err != nil {
		utils.Respond(rw, utils.Message(false, "Invalid request"))
		return
	}

	resp := models.LoginUser(db, token.Token)

	utils.Respond(rw, resp)
}

func GetUser(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	token := &models.Token{}

	err := json.NewDecoder(r.Body).Decode(token)
	if err != nil {
		resp := utils.Message(false, "Ivalid request")
		resp["is_err"] = true
		utils.Respond(rw, resp)
		return
	}

	user, err := models.GetUser(db, token.Token)
	if err != nil {
		resp := utils.Message(false, "Ivalid request to database")
		resp["is_err"] = true
		utils.Respond(rw, resp)
		return
	}

	if user == nil {
		resp := utils.Message(false, "Have no current user with that token")
		resp["is_err"] = false
		utils.Respond(rw, resp)
		return
	}

	resp := utils.Message(true, "User was found")
	resp["email"] = user.Email
	resp["name"] = user.Name
	resp["rating"] = user.Rating
	resp["token"] = user.Token

	utils.Respond(rw, resp)
}
