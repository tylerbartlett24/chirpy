package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tylerbartlett24/chirpy/internal/auth"
	"github.com/tylerbartlett24/chirpy/internal/database"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
}


func (cfg *apiConfig) createUsersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, 
			"Error decoding parameters", err)
		return
	}

	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, 
			"Please supply a password.", nil)
		return
	}

	hashPwd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, 
			"Could not create password", err)
		return
	}

	dbParams := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashPwd,
	}
	user, err := cfg.Queries.CreateUser(r.Context(), dbParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}
	
	respBody := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,	
	}
	respondWithJSON(w, http.StatusCreated, respBody)
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, 
			"Error decoding parameters", err)
		return
	}

	fmt.Println(params.Email)
	user, err := cfg.Queries.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Incorrect email.", err)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password.", err)
		return
	}

	respBody := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}
	respondWithJSON(w, http.StatusOK, respBody)
}