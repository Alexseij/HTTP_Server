package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Alexseij/server/models"
	"github.com/Alexseij/server/utils"
)

func CreateUser(rw http.ResponseWriter, r *http.Request) {
	token := &models.TokenType{}

	err := json.NewDecoder(r.Body).Decode(token)
	if err != nil {
		utils.Respond(rw, utils.Message(false, "Ivalid request"))
		return
	}

	user := &models.User{
		Token: token.Token,
	}

	resp := user.Create()
	utils.Respond(rw, resp)
}

func LoginUser(rw http.ResponseWriter, r *http.Request) {
	token := &models.TokenType{}

	err := json.NewDecoder(r.Body).Decode(token)
	if err != nil {
		utils.Respond(rw, utils.Message(false, "Invalid request"))
	}

	resp := models.LoginUser(token.Token)

	utils.Respond(rw, resp)
}

func MakeOrder(rw http.ResponseWriter, r *http.Request) {
	order := &models.Order{}

	if err := json.NewDecoder(r.Body).Decode(order); err != nil {
		utils.Respond(rw, utils.Message(false, "Invalid request"))
	}

	resp := order.MakeOrder()

	utils.Respond(rw, resp)
}
