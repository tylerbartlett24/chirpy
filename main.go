package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const rootFilePath = "."
	const port = ":8080"

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(rootFilePath))))
	serveMux.HandleFunc("/healthz", readyHandler)
	
	server := &http.Server{
		Handler: serveMux,
		Addr: port,
	}
	
	fmt.Printf("Serving files from %s on port %s\n", rootFilePath, port)
	log.Fatal(server.ListenAndServe())
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	content := []byte(http.StatusText(http.StatusOK))
	w.Write(content)
}