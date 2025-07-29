package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	tokenString := strings.Replace(header, "Bearer ", "", 1)
	if tokenString == "" {
		return "", errors.New("couldn't get token string in authorization")
	}

	return tokenString, nil
}