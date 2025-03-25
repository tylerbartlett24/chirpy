package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const rootFilePath = "."
	const port = ":8080"

	rootHandler := http.StripPrefix("/app", http.FileServer(http.Dir(rootFilePath)))
	serveMux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(rootHandler))
	serveMux.HandleFunc("GET /api/healthz", readyHandler)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	serveMux.HandleFunc("POST /api/validate_chirp", validateHandler)
	
	server := &http.Server{
		Handler: serveMux,
		Addr: port,
	}
	
	fmt.Printf("Serving files from %s on port %s\n", rootFilePath, port)
	log.Fatal(server.ListenAndServe())
}