package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from database.", err)
		return
	}

	if len(chirps) == 0 {
		respondWithError(w, http.StatusInternalServerError, "No chirps found in database.", err)
		return
	}

	// Sets the new db chirps to the made chirpSlice to control the json key names
	chirpSlice := []chirp{}
	for _, ch := range chirps {
		chirpSlice = append(chirpSlice, chirp{
			ID: ch.ID,
			CreatedAt: ch.CreatedAt,
			UpdatedAt: ch.UpdatedAt,
			Body: ch.Body,
			UserId: ch.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpSlice)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse UUID from path value.", err)
		return 
	}

	foundChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp by id.", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp{
		ID: foundChirp.ID,
		CreatedAt: foundChirp.CreatedAt,
		UpdatedAt: foundChirp.UpdatedAt,
		Body: foundChirp.Body,
		UserId: foundChirp.UserID,
	})
}