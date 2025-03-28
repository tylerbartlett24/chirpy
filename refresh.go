package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/tylerbartlett24/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, 
			"Malformed request.", err)
		return
	}

	refToken, err := cfg.Queries.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, 
			"Refresh token not recognized.", err)
		return
	}

	if time.Now().After(refToken.ExpiresAt) {
		err = errors.New("Refresh token expired.")
		respondWithError(w, http.StatusUnauthorized, 
			err.Error(), err)
		return
	}

	if refToken.RevokedAt.Valid && time.Now().After(refToken.RevokedAt.Time) {
		err = errors.New("Refresh token revoked.")
		respondWithError(w, http.StatusUnauthorized, 
			err.Error(), err)
		return
	}


	newAccToken, err := auth.MakeJWT(refToken.UserID, cfg.Secret, 3600 * time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, 
			"Could not generate access token.", err)
		return
	}

	respBody := response{
		Token: newAccToken,
	}
	respondWithJSON(w, http.StatusOK, respBody)
}