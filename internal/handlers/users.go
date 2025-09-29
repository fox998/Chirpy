package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/fox998/Chirpy/internal/auth"
	"github.com/fox998/Chirpy/internal/common"
	"github.com/fox998/Chirpy/internal/config"
	"github.com/fox998/Chirpy/internal/database"
)

func UsersUpdate(c *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bearer, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("Failed to get bearer token %v\n", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		userId, err := auth.ValidateJWT(bearer, c.Env.Secret)
		if err != nil {
			log.Printf("Failed to validate jwt %v\n", err)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		defer r.Body.Close()
		var bodyJson struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err = json.NewDecoder(r.Body).Decode(&bodyJson)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode request body: %s", err.Error()), 400)
			return
		}

		if bodyJson.Email == "" || bodyJson.Password == "" {
			http.Error(w, "Email and password are required", 400)
			return
		}

		password_hash, err := auth.HashPassword(bodyJson.Password)
		if err != nil {
			log.Println("failed to create password hash: ", err.Error())
			http.Error(w, "failed to update user", 500)
			return
		}

		dbUpdatedUser, err := c.Db.UpdateUser(r.Context(), database.UpdateUserParams{
			ID:           userId,
			Email:        bodyJson.Email,
			PasswordHash: password_hash,
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "failed to update user", 500)
			return
		}

		err = json.NewEncoder(w).Encode(common.ResponceUserData{
			ID:          dbUpdatedUser.ID,
			Email:       dbUpdatedUser.Email,
			IsChirpyRed: dbUpdatedUser.IsChirpyRed,
			CreatedAt:   dbUpdatedUser.CretedAt,
			UpdatedAt:   dbUpdatedUser.UpdatedAt,
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "failed to construct responce", 500)
			return
		}
	}
}

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

		writer.WriteHeader(201)
		err = json.NewEncoder(writer).Encode(common.ResponceUserData{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
			CreatedAt:   user.CretedAt,
			UpdatedAt:   user.UpdatedAt,
		})
		if err != nil {
			http.Error(writer, fmt.Sprintf("failed to encode: %s", err.Error()), 500)
			return
		}

	}
}
