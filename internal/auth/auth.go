package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts an API Key from the header of an HTTP request
// Example: "Authorization: ApiKey {insert API key here}"
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no API key in the Authorization header")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("invalid Authorization header format")
	}
	if vals[0] != "ApiKey" {
		return "", errors.New("invalid first part of the Authorization header format")
	}
	return vals[1], nil
}
