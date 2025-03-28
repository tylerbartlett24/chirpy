package auth

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/tylerbartlett24/chirpy/internal/database"
)

func TestJWTMatch(t *testing.T) {
	db, err := sql.Open("postgres", 
	"postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable")
	if err != nil {
		t.Errorf("could not open database: %v", err)
	}
	dbQueries := database.New(db)

	userID, err := dbQueries.GetRandomUserId(context.Background())
	if err != nil {
		t.Errorf("could not get user id: %v", err)
	}

	duration, _ := time.ParseDuration("1h")
	token, err := MakeJWT(userID, "hello", duration)
	if err != nil {
		t.Errorf("could not make JWT: %v", err)
	}

	newID, err := ValidateJWT(token, "hello")
	if err != nil {
		t.Errorf("JWT did not match: %v", err)
	}

	if newID != userID {
		t.Errorf("%v(new_id) != %v(user_id)", newID, userID)
	}
}

func TestJWTExpire(t *testing.T) {
	db, err := sql.Open("postgres", 
	"postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable")
	if err != nil {
		t.Errorf("could not open database: %v", err)
	}
	dbQueries := database.New(db)

	userID, err := dbQueries.GetRandomUserId(context.Background())
	if err != nil {
		t.Errorf("could not get user id: %v", err)
	}

	duration, _ := time.ParseDuration("1ms")
	token, err := MakeJWT(userID, "hello", duration)
	if err != nil {
		t.Errorf("could not make JWT: %v", err)
	}

	time.Sleep(2 * time.Millisecond)
	_, err = ValidateJWT(token, "hello")
	if err == nil {
		t.Errorf("Expired token evaluated as not expired.")
	}
}

func TestJWTNoMatch(t *testing.T) {
	db, err := sql.Open("postgres", 
	"postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable")
	if err != nil {
		t.Errorf("could not open database: %v", err)
	}
	dbQueries := database.New(db)

	userID, err := dbQueries.GetRandomUserId(context.Background())
	if err != nil {
		t.Errorf("could not get user id: %v", err)
	}

	duration, _ := time.ParseDuration("1h")
	token, err := MakeJWT(userID, "hello", duration)
	if err != nil {
		t.Errorf("could not make JWT: %v", err)
	}

	_, err = ValidateJWT(token, "goodbye")
	if err == nil {
		t.Error("Error authenticating: incorrect secret key validated")
	}

}