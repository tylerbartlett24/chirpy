package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHead := headers.Get("Authorization")
	if authHead == "" {
		return "", errors.New("NO authorization header in request.")
	}

	splitAuth := strings.Split(authHead, " ")
	if len(splitAuth) != 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}