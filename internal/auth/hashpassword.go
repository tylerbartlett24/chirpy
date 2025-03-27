package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	byteHash, err := bcrypt.GenerateFromPassword(bytePassword, 0)
	if err != nil {
		return "", err
	}

	stringHash := string(byteHash)
	return stringHash, nil

}