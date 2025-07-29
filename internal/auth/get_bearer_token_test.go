package auth

import (
	"net/http"
	"testing"
)


func TestGetBearerToken(t *testing.T) {
	tokenString := "ThisIsTheTokenString"

	header := http.Header{}
	header.Set("Authorization", tokenString)

	token, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("Couldn't get bearer token: %s", err)
	}

	if token != tokenString {
		t.Errorf("Token doesn't match token string")
	}
}

func TestGetBearerTokenFail(t *testing.T) {
	tokenString := "ThisIsTheTokenString"

	header := http.Header{}
	header.Set("Authorization", tokenString)

	token, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("Couldn't get bearer token: %s", err)
	}

	if token == "" {
		t.Errorf("Token is not supposed to equal an empty string")
	}
}