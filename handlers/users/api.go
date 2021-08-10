package users

import (
	"net/http"
	"os"
	"strconv"

	"github.com/Alexseij/server/models"
	"github.com/Alexseij/server/utils"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateRating(db *mongo.Database, rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	user, err := models.GetUser(db, vars["email"])
	if err != nil {
		utils.Respond(rw, utils.Message(false, "Ivalid getting user operation"))
		return
	}

	if user == nil {
		utils.Respond(rw, utils.Message(false, "User is empty"))
		return
	}

	rating, err := strconv.Atoi(vars["rating"])
	if err != nil {
		utils.Respond(rw, utils.Message(false, "Invalid rquest"))
		return
	}

	resp := user.UpdateRating(db, rating, os.Getenv("rating_topic_name"), os.Getenv("project_id"))

	utils.Respond(rw, resp)
}
