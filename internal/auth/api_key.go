package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts the API key from the Authorization header
// Expected format: "ApiKey <api_key>"
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no authorization header included")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("malformed authorization header")
	}

	if parts[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return parts[1], nil
}
