package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	user := User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode decoder", err)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create user", err)
		return
	}

	// Sets all the db users values to the local user, so we can have control over the json naming
	user.ID = dbUser.ID
	user.CreatedAt = dbUser.CreatedAt
	user.UpdatedAt = dbUser.UpdatedAt
	user.Email = dbUser.Email

	respondWithJSON(w, http.StatusCreated, user)
}