package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/FreyFam5/go/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db *database.Queries
	fileserverHits atomic.Int32
}

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
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
	}

	mux := http.ServeMux{}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", apicfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)
	mux.HandleFunc("POST /api/users", apicfg.handlerUsers)

	mux.HandleFunc("GET /admin/metrics", apicfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apicfg.handlerReset)

	server := &http.Server{
		Handler: &mux,
		Addr: ":" + port,
	}

	log.Printf("starting server for port: %s", port)
	log.Fatal(server.ListenAndServe())
}
