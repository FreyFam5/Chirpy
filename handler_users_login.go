package main

import (
	"encoding/json"
	"net/http"

	"github.com/FreyFam5/go/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	user := User{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode decoder", err)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), user.Email) 
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	if err = auth.CheckPasswordHash(user.HashedPassword, dbUser.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	})
}