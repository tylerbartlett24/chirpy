package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tylerbartlett24/chirpy/internal/auth"
	"github.com/tylerbartlett24/chirpy/internal/database"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpsHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Body string `json:"body"`
	}
	
	const maxChirpLength = 140

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, 
			"No token provided.", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, 
			"Could not validate token.", err)
		return
	}

	if len(params.Body) > maxChirpLength {
		err = errors.New("chirp too long")
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	params.Body = badWordsFilter(params.Body)
	
	chirpParams := database.CreateChirpParams{
		Body: params.Body,
		UserID: userID,
	}
	newChirp, err := cfg.Queries.CreateChirp(r.Context(), chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respBody := Chirp{
		ID: newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body: params.Body,
		UserID: userID,
	}
	respondWithJSON(w, http.StatusCreated, respBody)
}

func (cfg *apiConfig) readChirpsHandler(w http.ResponseWriter, r *http.Request) {
	chirps , err := cfg.Queries.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
		return
	}

	length := len(chirps)
	respBody := make([]Chirp, length)
	for i, chirp := range chirps {
		respBody[i] = Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		}
	}
	respondWithJSON(w, http.StatusOK, respBody)
}

func (cfg *apiConfig) readChirpHandler(w http.ResponseWriter, r *http.Request) {
	stringUUID := r.PathValue("chirpID")
	id, err := uuid.Parse(stringUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Invalid Id", err)
	   return
	}
	
	dbChirp , err := cfg.Queries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound,
			"Could not find chirp", err)
		return
	}

	respBody := Chirp{
		ID: dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body: dbChirp.Body,
		UserID: dbChirp.UserID,
	}
	respondWithJSON(w, http.StatusOK, respBody)
}