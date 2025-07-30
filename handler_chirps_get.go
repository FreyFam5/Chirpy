package main

import (
	"net/http"
	"sort"

	"github.com/FreyFam5/go/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	var err error
	queryID := r.URL.Query().Get("author_id")
	if queryID == "" { // If query is not input, grabs all chirps
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps from database", err)
			return
		}
	} else { // If query is input, will grab all the chirps with that user id
		userID, err := uuid.Parse(queryID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get uuid from query ID", err)
			return
		}
		chirps, err = cfg.db.GetChirpsByUserID(r.Context(), userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps by id from database", err)
			return
		}
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

	// If the query is desc, sorts the chirps in descending order
	if querySort := r.URL.Query().Get("sort"); querySort == "desc" {
		sort.Slice(chirpSlice, func(i, j int) bool {
			return chirpSlice[i].CreatedAt.After(chirpSlice[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirpSlice)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse UUID from path value", err)
		return 
	}

	foundChirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", err)
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