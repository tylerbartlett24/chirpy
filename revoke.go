package main

import (
	"net/http"

	"github.com/tylerbartlett24/chirpy/internal/auth"
)

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, 
			"Malformed request.", err)
		return
	}

	err = cfg.Queries.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, 
			"Token could not be revoked", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}