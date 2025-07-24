package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)


func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	reqB := reqBody{}
	// After decode if it returns an error, server responds with a internal error status code and a json error is written
	if err := decoder.Decode(&reqB); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode decoder", err)
		return
	}

	// If the request body is over 140 characters, then sends a bad request status and a json error is written
	if len(reqB.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	cleanedWords := getCleanedBody(reqB.Body, profaneWords)

	type responseBody struct{
			CleanedBody string `json:"cleaned_body"`
		}

	respBody := responseBody{CleanedBody: cleanedWords}
	respondWithJSON(w, http.StatusOK, respBody)
}


func getCleanedBody(body string, profaneWords []string) string {
	words := strings.Split(body, " ")
	for idx, word := range words {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			words[idx] = "****"
		}
	}
	return strings.Join(words, " ")
}

