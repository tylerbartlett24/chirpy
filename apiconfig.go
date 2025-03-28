package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/tylerbartlett24/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	Queries *database.Queries
	Platform string
	Secret string
}

const metricsHTML = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	text := fmt.Sprintf(metricsHTML, cfg.fileserverHits.Load())
	content := []byte(text)
	w.Write(content)
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.Queries.Reset(r.Context())
	if err != nil {
		log.Printf("Could not reset database: %v", err)
		w.WriteHeader(400)
		return
	}
	w.WriteHeader(http.StatusOK)
}

