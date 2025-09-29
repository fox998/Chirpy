package handlers

import (
	"log"
	"net/http"

	"github.com/fox998/Chirpy/internal/auth"
	"github.com/fox998/Chirpy/internal/config"
	"github.com/google/uuid"
)

func Revoke(c *config.ApiConfig) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		bearer, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("Failed to get bearer refresh token %v\n", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		id, err := uuid.Parse(bearer)
		if err != nil {
			log.Printf("Failed to parse bearer refresh token %v\n", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		err = c.Db.RevokeRefreshToken(r.Context(), id)
		if err != nil {
			log.Printf("Failed to revoke refresh token from db %v\n", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
