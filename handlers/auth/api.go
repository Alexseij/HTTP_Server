package auth

import (
	"log"
	"net/http"

	"github.com/Alexseij/server/models"
	"github.com/Alexseij/server/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoginUser(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Print("Request with token : ", vars["token"])
	resp := models.LoginUser(db, vars["token"])

	utils.Respond(rw, resp)
}
