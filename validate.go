package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	const tooLongMsg = "Chirp is too long."
	const maxChirpLength = 140

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, tooLongMsg)
		return
	}
	params.Body = badWordsFilter(params.Body)
	

	
	respBody := returnVals{
		CleanedBody: params.Body,
	}
	respondWithJSON(w, http.StatusOK, respBody)
}