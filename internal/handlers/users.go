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

func Users(c *config.ApiConfig) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var bodyJson struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&bodyJson)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Failed to decode request body: %s", err.Error()), 400)
			return
		}

		passwordHash, err := auth.HashPassword(bodyJson.Password)
		if err != nil {
			log.Println("failed to create password hash: ", err.Error())
			http.Error(writer, "failed to create user", 500)
			return
		}

		user, err := c.Db.CreateUser(r.Context(), database.CreateUserParams{
			Email:        bodyJson.Email,
			PasswordHash: passwordHash,
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(writer, "failed to create userd in db", 500)
			return
		}

		type resBody struct {
			ID        uuid.UUID `json:"id"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
			Email     string    `json:"email"`
		}

		writer.WriteHeader(201)
		err = json.NewEncoder(writer).Encode(resBody{
			ID:        user.ID,
			CreatedAt: user.CretedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		})
		if err != nil {
			http.Error(writer, fmt.Sprintf("failed to encode: %s", err.Error()), 500)
			return
		}

	}
}
