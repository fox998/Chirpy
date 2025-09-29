package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/fox998/Chirpy/internal/auth"
	"github.com/fox998/Chirpy/internal/config"
	"github.com/google/uuid"
)

func Refresh(c *config.ApiConfig) http.HandlerFunc {

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

		dbRefreshToken, err := c.Db.GetRefreshTokenById(r.Context(), id)
		if err != nil {
			log.Printf("Failed to get refresh token from db %v\n", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		if dbRefreshToken.RevokedAt.Valid {
			log.Printf("Refresh token is revoked %v\n", dbRefreshToken.RevokedAt.Time)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		if dbRefreshToken.ExpairesAt.Before(time.Now()) {
			log.Printf("Refresh token is expired %v\n", dbRefreshToken.ExpairesAt)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		jwt, err := auth.MakeJWT(dbRefreshToken.UserID.UUID, c.Env.Secret, time.Hour)
		if err != nil {
			log.Printf("Failed to create jwt for %v: %v", dbRefreshToken.UserID.UUID, err)
			http.Error(w, "Failed to authorathe", 500)
			return
		}

		respBody := struct {
			JWT string `json:"token"`
		}{
			JWT: jwt,
		}

		err = json.NewEncoder(w).Encode(respBody)
		if err != nil {
			log.Printf("Failed to encode response body: %v", err)
			http.Error(w, "Failed to encode response body", 500)
			return
		}
	}
}
