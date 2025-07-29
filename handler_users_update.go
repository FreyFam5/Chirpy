package main

import (
	"encoding/json"
	"net/http"

	"github.com/FreyFam5/go/chirpy/internal/auth"
	"github.com/FreyFam5/go/chirpy/internal/database"
)


func (cfg *apiConfig) handlerUsersUpdateLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// - Getting User ID -
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to get token", err)
		return
	} else if token == "" {
		respondWithError(w, http.StatusUnauthorized, "Missing token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to validate token", err)
		return
	}
	// X Getting User ID X

	// - Updating user password & email -
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	dbUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: userID,
		Email: params.Email,
		HashedPassword: hash,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}
	// X Updating user password & email X

	respondWithJSON(w, http.StatusOK, User{
		ID: dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email: dbUser.Email,
	})
}
