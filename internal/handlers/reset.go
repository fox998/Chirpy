package handlers

import (
	"fmt"
	"net/http"

	"github.com/fox998/Chirpy/internal/config"
	"github.com/fox998/Chirpy/internal/config/env/platform"
)

func Reset(config *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if config.Env.Platform != platform.Dev {
			http.Error(w, "Forbiden", http.StatusForbidden)
			return
		}

		err := config.Db.ResetUsers(r.Context())
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to reset users: %s", err.Error()), 500)
			return
		}

		w.WriteHeader(200)
		config.ResetHits()
	}
}
