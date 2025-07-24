package main

import (
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

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apicfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apicfg.handlerReset)

	port := "8080"

	server := &http.Server{
		Handler: &mux,
		Addr: ":" + port,
	}

	log.Printf("starting server for port: %s", port)
	log.Fatal(server.ListenAndServe())
}
