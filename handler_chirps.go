package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/FreyFam5/go/chirpy/internal/auth"
	"github.com/FreyFam5/go/chirpy/internal/database"
	"github.com/google/uuid"
)


type chirp struct {
		ID 				uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body 			string 		`json:"body"`
		UserId 		uuid.UUID `json:"user_id"`
	}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	c, ok := validateChirp(w, r, cfg)
	if !ok {
		return
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: c.Body,
		UserID: c.UserId,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	// Sets the new db chirp to the made chirp to control the json key names
	c.ID = dbChirp.ID
	c.CreatedAt = dbChirp.CreatedAt
	c.UpdatedAt = dbChirp.UpdatedAt
	c.Body = dbChirp.Body
	c.UserId = dbChirp.UserID

	respondWithJSON(w, http.StatusCreated, c)
}

// Cleans the chirp given and returns it if it was valid
func validateChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) (chirp, bool){
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	c := chirp{}
	// After decode if it returns an error, server responds with a internal error status code and a json error is written
	if err := decoder.Decode(&c); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode decoder", err)
		return chirp{}, false
	}

	// If the request body is over 140 characters, then sends a bad request status and a json error is written
	if len(c.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return chirp{}, false
	}

	// Gets the the token from the requests header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get token", err)
		return chirp{}, false
	}

	// Gets the user id of the associated token
	c.UserId, err = auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not authorized to view this chirp", err)
		return chirp{}, false
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	c.Body = getCleanedBody(c.Body, profaneWords)
	
	return c, true
}

// Returns a cleaned version of the given string
func getCleanedBody(body string, profaneWords []string) string {
	words := strings.Split(body, " ")
	for idx, word := range words {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			words[idx] = "****"
		}
	}
	return strings.Join(words, " ")
}
