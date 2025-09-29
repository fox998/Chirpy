package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	headerData := headers.Get("Authorization")
	apiKey, found := strings.CutPrefix(headerData, "ApiKey ")
	if !found || apiKey == "" {
		return "", errors.New("no api key provided")
	}

	return apiKey, nil
}
