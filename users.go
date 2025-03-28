package main

import (
	"encoding/json"
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
	Token string `json:"token"`
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
		respondWithError(w, http.StatusInternalServerError, 
			"Couldn't create user", err)
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
		Expires int `json:"expires_in_seconds"`
	}
	
	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, 
			"Error decoding parameters", err)
		return
	}

	// 1 hour is default and maximum expiration time
	if params.Expires == 0 || params.Expires > 3600 {
		params.Expires = 3600
	}
	user, err := cfg.Queries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email.", err)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect password.",
		 err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.Secret, 
		(time.Duration(params.Expires) * time.Second))
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, 
			"Failed to generate token.", err)
		return
	}

	respBody := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	}
	respondWithJSON(w, http.StatusOK, respBody)
}