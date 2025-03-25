package main

import "net/http"

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	content := []byte(http.StatusText(http.StatusOK))
	w.Write(content)
}