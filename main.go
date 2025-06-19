package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/Burncoat/Chirpy/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
		fileserverHits atomic.Int32
		db 			   *database.Queries
		platform 	   string
		jwtSecret      string
	}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	dbQueries := database.New(db)

	const filePathRoot = "."
	const port = "8080"
	

	// Pulling from index.html
	mux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db: 			dbQueries,
		platform:	    platform,
		jwtSecret:      jwtSecret,
	}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(handler))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpGetByID)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Serving files form %s on port: %s\n", filePathRoot, port)
	log.Fatal(srv.ListenAndServe())
}


