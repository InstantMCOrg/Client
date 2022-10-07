package server

import (
	"encoding/json"
	"net/http"
)

func CreateResponse(w http.ResponseWriter, message string, statusCode int) {
	data, _ := json.Marshal(map[string]string{
		"message": message,
	})
	w.WriteHeader(statusCode)
	w.Write(data)
}
