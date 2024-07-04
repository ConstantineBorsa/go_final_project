package main

import (
	"encoding/json"
	"net/http"
)

func sendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendErrorResponse(w http.ResponseWriter, status int, response ErrorResponse) {
	sendResponse(w, status, response)
}
