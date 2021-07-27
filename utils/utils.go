package utils

import (
	"encoding/json"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

func Respond(rw http.ResponseWriter, response map[string]interface{}) {
	rw.Header().Add("Content-type", "application/json")
	json.NewEncoder(rw).Encode(response)
}
