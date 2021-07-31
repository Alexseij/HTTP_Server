package auth

import (
	"encoding/json"
	"net/http"

	"github.com/Alexseij/server/models"
	"github.com/Alexseij/server/utils"
)

func CreateUser(rw http.ResponseWriter, r *http.Request) {
	token := &models.Token{}

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
	token := &models.Token{}

	err := json.NewDecoder(r.Body).Decode(token)
	if err != nil {
		utils.Respond(rw, utils.Message(false, "Invalid request"))
		return
	}

	resp := models.LoginUser(token.Token)

	utils.Respond(rw, resp)
}
