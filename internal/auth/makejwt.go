package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	currentDate := jwt.NewNumericDate(time.Now().UTC())
	claims := jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: currentDate,
		ExpiresAt: jwt.NewNumericDate(currentDate.Add(expiresIn)),
		Subject: userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	byteKey := []byte(tokenSecret)
	signedToken, err := token.SignedString(byteKey)
	if err != nil {
		return "", err
	}
	return signedToken, err
}