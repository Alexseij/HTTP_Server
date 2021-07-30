package order

import (
	"encoding/json"
	"net/http"

	"github.com/Alexseij/server/models"
	"github.com/Alexseij/server/utils"
)

func MakeOrder(rw http.ResponseWriter, r *http.Request) {
	order := &models.Order{}

	if err := json.NewDecoder(r.Body).Decode(order); err != nil {
		utils.Respond(rw, utils.Message(false, "Invalid request"))
	}

	resp := order.MakeOrder()

	utils.Respond(rw, resp)
}
