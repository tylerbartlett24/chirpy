package auth

import "golang.org/x/crypto/bcrypt"

func CheckPasswordHash(hash, password string) error {
	byteHash := []byte(hash)
	bytePassword := []byte(password)
	return bcrypt.CompareHashAndPassword(byteHash, bytePassword)
}