package main

import (
	"net/http"
	"time"

	"github.com/FreyFam5/go/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type Token struct {
		Token string `json:"token"`
	}

	// Gets the refresh token's string from the header
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token in header", err)
		return
	}

	// Gets the refresh token from the database
	dbRToken, err := cfg.db.GetRefreshToken(r.Context(), rToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token expired or incorrect", err)
		return
	}

	// Stops the refresh if the users token is revoked
	if dbRToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token revoked", nil)
		return
	} 

	// Stops the refresh if the token is expired
	if time.Now().After(dbRToken.ExpiresAt) {
		respondWithError(w, http.StatusUnauthorized, "Token expired, login to access", nil)
		return
	}

	// Gets the user from the database from the refresh token
	dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), dbRToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user from refresh token", err)
		return
	}

	// Gets and saves the JWT access token for user
	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Duration(3600) * time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Token{Token: token})
}


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token in header", err)
		return
	}

	if err := cfg.db.RevokeToken(r.Context(), rToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
	}

	w.WriteHeader(http.StatusNoContent)
}