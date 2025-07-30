package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/FreyFam5/go/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUserToRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	// Checks api key so no one unauthorized can upgrade a user to chirp red
	if apiKey, err := auth.GetAPIKey(r.Header); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get api key", err)
		return
	} else if apiKey != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "Incorrect api key", nil)
		return
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if _, err := cfg.db.UpgradeUser(r.Context(), params.Data.UserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Couldn't find given user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}