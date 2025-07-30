package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	apiKeyString := strings.Replace(header, "ApiKey ", "", 1)
	if apiKeyString == "" {
		return "", errors.New("couldn't get apiKey string in authorization")
	}

	return apiKeyString, nil
}