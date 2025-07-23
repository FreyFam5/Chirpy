package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.ServeMux{}
	
	apicfg := apiConfig{}

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apicfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apicfg.handlerMetrics)
	mux.HandleFunc("POST /reset", apicfg.handlerReset)

	port := "8080"

	server := &http.Server{
		Handler: &mux,
		Addr: ":" + port,
	}

	log.Printf("starting server for port: %s", port)
	log.Fatal(server.ListenAndServe())
}


// Adds a hit every time the page is visited
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Makes a page that shows the amount fo hits visually
func (cfg * apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}
