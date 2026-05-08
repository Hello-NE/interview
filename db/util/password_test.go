package util

import "testing"

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hashedPassword1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}
	hashedPassword2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}
	if hashedPassword1 == hashedPassword2 {
		t.Errorf("HashPassword returned the same hash for the same password, expected different hashes")
	}
}
