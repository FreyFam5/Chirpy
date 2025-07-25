package main

import (
	"net/http"
)

// Resets the hits count to 0 and shows on the page that it was reset
func (cfg * apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	cfg.fileserverHits.Store(0)

	if err := cfg.db.ResetUsers(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset the database: " + err.Error()))
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset Hits To 0 and reset database to initial state."))
}
