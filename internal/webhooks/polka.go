package webhooks

import (
	"encoding/json"
	"net/http"

	"github.com/fox998/Chirpy/internal/auth"
	"github.com/fox998/Chirpy/internal/config"
	"github.com/fox998/Chirpy/internal/database"
	"github.com/google/uuid"
)

func Polka(c *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil || apiKey != c.Env.PolkaKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var bodyJson struct {
			Event string `json:"event"`
			Data  struct {
				UserID string `json:"user_id"`
			} `json:"data"`
		}

		defer r.Body.Close()
		err = json.NewDecoder(r.Body).Decode(&bodyJson)
		if err != nil {
			http.Error(w, "Failed to decode request body", 400)
			return
		}

		if bodyJson.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		_, err = c.Db.UpdateUserIsChirpyRed(r.Context(), database.UpdateUserIsChirpyRedParams{
			ID:          uuid.MustParse(bodyJson.Data.UserID),
			IsChirpyRed: true,
		})
		if err != nil {
			http.Error(w, "Failed to update user", 500)
			return
		}

		w.WriteHeader(204)
	}
}
