package utils

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	Pattern = "g.nsu.ru"
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

func CheckDomain(email string) bool {
	indexRune := strings.IndexRune(email, '@')
	domainStr := email[indexRune+1:]
	if domainStr == Pattern {
		return true
	}
	return false
}
