package main

import "net/http"

// Resets the hits count to 0 and shows on the page that it was reset
func (cfg * apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Reset Hits To 0"))
}
