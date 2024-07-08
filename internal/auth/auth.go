package auth

import (
	"errors"
	"net/http"
	"strings"
)

var errorAuth = errors.New("Authorization Header Not Correct")

func GetAPIKEY(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errorAuth
	}

	splitAuth := strings.Split(authHeader, " ")

	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errorAuth
	}

	return splitAuth[1], nil

}
