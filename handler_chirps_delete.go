package main

import (
	"net/http"

	"github.com/FreyFam5/go/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to get token.", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to validate token.", err)
		return
	}

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirp id from path value.", err)
		return 
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp.", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Not allowed to delete others chirps.", nil)
		return
	}

	if err := cfg.db.DeleteChirpByID(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp.", err)
	}

	w.WriteHeader(http.StatusNoContent)
}