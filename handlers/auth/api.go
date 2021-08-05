package auth

import (
	"log"
	"net/http"

	"github.com/Alexseij/server/models"
	"github.com/Alexseij/server/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user := &models.User{
		Token: vars["token"],
	}

	resp := user.Create(db)
	utils.Respond(rw, resp)
}

func LoginUser(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	resp := models.LoginUser(db, vars["token"])

	utils.Respond(rw, resp)
}

func GetUser(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	log.Print("Icoming request : ", r.Body)

	vars := mux.Vars(r)

	user, err := models.GetUser(db, vars["token"])
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
