package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FreyFam5/go/chirpy/internal/auth"
	"github.com/FreyFam5/go/chirpy/internal/database"
)


func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	param := parameter{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameter", err)
		return
	}
	// Gets user by email
	dbUser, err := cfg.db.GetUserByEmail(r.Context(), param.Email) 
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// Checks to see if password is correct
	if err = auth.CheckPasswordHash(param.Password, dbUser.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	// Gets the JWT access token for user
	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Duration(3600) * time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make JWT", err)
		return
	}
	// Gets the refresh tokens string
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't make refresh token string", err)
		return
	}
	// Gets and saves refresh token for 60 days
	dbRT, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: dbUser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token to database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID: dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email: dbUser.Email,
			IsChirpyRed: dbUser.IsChirpyRed,
		},
		Token: token,
		RefreshToken: dbRT.Token,
	})
}