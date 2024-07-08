package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

func sendResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Error encoding response data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func sendErrorResponse(w http.ResponseWriter, status int, response ErrorResponse) {
	sendResponse(w, status, response)
}
