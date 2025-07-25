package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/FreyFam5/go/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db *database.Queries
	platform string
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	// Dev grab
	platform := os.Getenv("PLATFORM")
	// Url grab and check
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("couldn't open %s: %s", dbURL, err)
	}
	dbQueries := database.New(db)

	apicfg := apiConfig{
		db: dbQueries,
		platform: platform,
		fileserverHits: atomic.Int32{},
	}

	mux := http.ServeMux{}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apicfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/chirps", apicfg.handlerCreateChirp)
	mux.HandleFunc("POST /api/users", apicfg.handlerUsersCreate)

	mux.HandleFunc("GET /admin/metrics", apicfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apicfg.handlerReset)

	server := &http.Server{
		Handler: &mux,
		Addr: ":" + port,
	}

	log.Printf("starting server for port: %s", port)
	log.Fatal(server.ListenAndServe())
}
