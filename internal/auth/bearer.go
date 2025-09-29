package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return "", errors.New("authorization header is empty")
	}

	token, found := strings.CutPrefix(bearer, "Bearer ")
	if !found {
		return "", fmt.Errorf("no bearer prefix in authorization header")
	}

	return token, nil
}
