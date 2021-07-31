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
