package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Writes the given error to the server
func respondWithError(w http.ResponseWriter, statusCode int, msg string, err error) {
	if err != nil {
		log.Printf("error: %s\n", err)
	}

	if statusCode > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorBody struct {
		Error string `json:"error"`
	}
	
	respondWithJSON(w, statusCode, errorBody{Error: msg})
}

// Writes the given parameters json to the server
func respondWithJSON(w http.ResponseWriter, statusCode int, parameter interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(parameter)
	if err != nil {
		log.Printf("error marshalling %v: %s", parameter, err)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(data)
}