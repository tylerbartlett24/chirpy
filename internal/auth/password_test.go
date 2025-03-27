package auth

import "testing"

func TestPasswordMatch(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("hash didn't work: %v", err)
		return
	}

	err = CheckPasswordHash(hash, "password")
	if err != nil {
		t.Errorf("hash did not match: %v", err)
	}
}

func TestPasswordNoMatch(t *testing.T) {
	password := "password"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("hash didn't work: %v", err)
		return
	}

	err = CheckPasswordHash(hash, "fuck")
	if err == nil {
		t.Error("hash matched for different strings")
	}
}