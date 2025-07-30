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
	secret string
	apiKey string
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	// Dev grab
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	// Secret JWT grab
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET must be set")
	}
	// Polka key grab
	apiKey := os.Getenv("POLKA_KEY")
	if apiKey == "" {
		log.Fatal("POLKA_KEY must be set")
	}
	// Url grab and check
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	// Opens database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("couldn't open %s: %s", dbURL, err)
	}
	dbQueries := database.New(db)

	apicfg := apiConfig{
		db: dbQueries,
		platform: platform,
		secret: secret,
		apiKey: apiKey,
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apicfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/chirps", apicfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apicfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apicfg.handlerGetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apicfg.handlerDeleteChirpByID)

	mux.HandleFunc("POST /api/users", apicfg.handlerUsersCreate)
	mux.HandleFunc("PUT /api/users", apicfg.handlerUsersUpdateLogin)
	mux.HandleFunc("POST /api/login", apicfg.handlerUsersLogin)

	mux.HandleFunc("POST /api/polka/webhooks", apicfg.handlerUpgradeUserToRed)

	mux.HandleFunc("POST /api/refresh", apicfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apicfg.handlerRevoke)

	mux.HandleFunc("GET /admin/metrics", apicfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apicfg.handlerReset)

	server := &http.Server{
		Handler: mux,
		Addr: ":" + port,
	}

	log.Printf("starting server for port: %s", port)
	log.Fatal(server.ListenAndServe())
}
