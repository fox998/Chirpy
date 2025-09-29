package handlers

import (
	"fmt"
	"net/http"

	"github.com/fox998/Chirpy/internal/config"
)

func Metrics(config *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := fmt.Sprintf(
			`
			<html>
			<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
			</body>
			</html>
			`, config.GetHits())

		w.WriteHeader(200)
		w.Write([]byte(body))
	}
}
