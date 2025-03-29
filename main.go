package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/tylerbartlett24/chirpy/internal/database"
)

func main() {
	const rootFilePath = "."
	const port = ":8080"
	
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secret := os.Getenv("SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Could not open database: %v", err)
	}
	dbQueries := database.New(db)

	rootHandler := http.StripPrefix("/app", http.FileServer(http.Dir(rootFilePath)))
	serveMux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		Queries: dbQueries,
		Platform: platform,
		Secret: secret,
	}

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(rootHandler))
	serveMux.HandleFunc("GET /api/healthz", readyHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.createChirpsHandler)
	serveMux.HandleFunc("POST /api/users", apiCfg.createUsersHandler)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.readChirpsHandler)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.readChirpHandler)
	serveMux.HandleFunc("POST /api/login", apiCfg.loginHandler)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.refreshHandler)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.revokeHandler)
	serveMux.HandleFunc("PUT /api/users", apiCfg.updateUsersHandler)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.deleteChirpsHandler)
	
	server := &http.Server{
		Handler: serveMux,
		Addr: port,
	}
	
	fmt.Printf("Serving files from %s on port %s\n", rootFilePath, port)
	log.Fatal(server.ListenAndServe())
}