package auth

import (
	"testing"
)


func TestHashPassword(t *testing.T) {
	password := "ThisIsAPassword!"
	_, err := HashPassword(password)
	if err != nil {
		t.Errorf("Couldn't hash password %s: %s", password, err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "ThisIsAPassword"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("Couldn't hash password %s: %s", password, err)
	}

	if err = CheckPasswordHash(password, hash); err != nil {
		t.Errorf("Hash '%s' does not match password '%s': %s", hash, password, err)
	}
}