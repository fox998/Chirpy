package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fox998/Chirpy/internal/auth"
	"github.com/fox998/Chirpy/internal/config"
	"github.com/fox998/Chirpy/internal/database"
	"github.com/google/uuid"
)

func validateExpireIn(inSec uint) time.Duration {
	defaultDuration := time.Hour
	if inSec == 0 || inSec > uint(defaultDuration.Seconds()) {
		return defaultDuration
	}

	return time.Second * time.Duration(inSec)

}

func Login(config *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&reqBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %s", err.Error()), 400)
			return
		}

		dbUser, err := config.Db.GetUserByEmail(r.Context(), reqBody.Email)
		if err != nil {
			log.Printf("failed to get user by emain %s : %v", reqBody.Email, err)
			http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
			return
		}

		ok, err := auth.CheckPasswordHash(reqBody.Password, dbUser.PasswordHash)
		if err != nil {
			log.Println("failed to check password ", err)
			http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
			return
		}

		if !ok {
			http.Error(w, "Incorrect email or password", http.StatusUnauthorized)
			return
		}

		jwt, err := auth.MakeJWT(dbUser.ID, config.Env.Secret, time.Hour)
		if err != nil {
			log.Printf("failed to create jwt for %v: %v", dbUser.ID, err)
			http.Error(w, "Failed to authorathe", 500)
			return
		}

		refreshToken, err := config.Db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			UserID:     uuid.NullUUID{UUID: dbUser.ID, Valid: true},
			ExpairesAt: time.Now().Add(time.Hour * 24 * 60),
		})

		if err != nil {
			log.Printf("Failed to create refresh token %v", err)
			http.Error(w, "Failed to authorathe", 500)
			return
		}

		type resBody struct {
			ID           string `json:"id"`
			CreatedAt    string `json:"created_at"`
			UpdatedAt    string `json:"updated_at"`
			Email        string `json:"email"`
			Token        string `json:"token"`
			Refreshtoken string `json:"refresh_token"`
		}

		err = json.NewEncoder(w).Encode(resBody{
			ID:           dbUser.ID.String(),
			CreatedAt:    dbUser.CretedAt.Format("2021-07-07T00:00:00Z"),
			UpdatedAt:    dbUser.UpdatedAt.Format("2021-07-07T00:00:00Z"),
			Email:        dbUser.Email,
			Token:        jwt,
			Refreshtoken: refreshToken.ID.String(),
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to encode: %s", err.Error()), 500)
			return
		}
	}
}
