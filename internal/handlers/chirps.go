package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/fox998/Chirpy/internal/auth"
	"github.com/fox998/Chirpy/internal/config"
	"github.com/fox998/Chirpy/internal/database"
	"github.com/google/uuid"
)

func cleanChirps(body string) string {
	words := strings.Split(body, " ")
	barWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	for i, word := range words {
		if _, found := barWords[strings.ToLower(word)]; found {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

type chips_responceFormat struct {
	Id         string `json:"id"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Body       string `json:"body"`
	User_id    string `json:"user_id"`
}

func chirpsDatabaseToResponce(dbChirp database.Chirp) chips_responceFormat {
	res := chips_responceFormat{}

	res.Id = dbChirp.ID.String()
	res.Body = dbChirp.Body
	res.Created_at = dbChirp.CreatedAt.Format("2021-01-01T00:00:00Z")
	res.Updated_at = dbChirp.UpdatedAt.Format("2021-01-01T00:00:00Z")
	res.User_id = dbChirp.UserID.UUID.String()

	return res
}

func ValidateAuth(c *config.ApiConfig, req *http.Request) (uuid.UUID, error) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		return uuid.UUID{}, err
	}

	jwtUserId, err := auth.ValidateJWT(token, c.Env.Secret)
	if err != nil {
		return uuid.UUID{}, err
	}

	return jwtUserId, nil
}

func Chirps(config *config.ApiConfig) http.HandlerFunc {

	return func(writer http.ResponseWriter, req *http.Request) {

		var reqBody struct {
			Body string `json:"body"`
		}
		defer req.Body.Close()

		err := json.NewDecoder(req.Body).Decode(&reqBody)
		if err != nil {
			log.Printf("Failed to decode req: %v\n", err)
			http.Error(writer, "failed to decode request body", 400)
			return
		}

		if len(reqBody.Body) > 140 {
			http.Error(writer, "body lenth should be less than 140", 400)
			return
		}

		jwtUserId, err := ValidateAuth(config, req)
		if err != nil {
			log.Printf("Failed to get jwt token %v\n", err)
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		dbChirp, err := config.Db.PostChirp(req.Context(), database.PostChirpParams{
			Body: cleanChirps(reqBody.Body),
			UserID: uuid.NullUUID{
				UUID:  jwtUserId,
				Valid: true,
			},
		})

		if err != nil {
			log.Println(err.Error())
			http.Error(writer, "Failed to post chirp", 500)
			return
		}

		writer.WriteHeader(201)
		json.NewEncoder(writer).Encode(chirpsDatabaseToResponce(dbChirp))
	}

}

func AllChirps(config *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirps, err := config.Db.ListChirps(r.Context())
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to obtaine chirps", 500)
			return
		}

		resChirps := make([]chips_responceFormat, len(chirps))
		for i, chirp := range chirps {
			resChirps[i] = chirpsDatabaseToResponce(chirp)
		}

		json.NewEncoder(w).Encode(resChirps)
	}
}

func ChirpsById(config *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "id path param expected: "+err.Error(), 400)
			return
		}

		dbChirp, err := config.Db.GetChirpByID(r.Context(), chirpId)
		if err != nil {
			log.Println(err)
			http.NotFound(w, r)
			return
		}

		json.NewEncoder(w).Encode(chirpsDatabaseToResponce(dbChirp))
	}
}

func DeleteChirps(config *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		chirpId, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "id path param expected: "+err.Error(), 400)
			return
		}

		authUserId, err := ValidateAuth(config, r)
		if err != nil {
			log.Printf("Failed to get jwt token %v\n", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		dbChirp, err := config.Db.GetChirpByID(r.Context(), chirpId)
		if err != nil {
			log.Println(err)
			http.NotFound(w, r)
			return
		}

		if dbChirp.UserID.UUID != authUserId {
			http.Error(w, "You can delete only your own chirps", http.StatusForbidden)
			return
		}

		err = config.Db.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{
			ID:     chirpId,
			UserID: dbChirp.UserID,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, "failed to delete chirp", 500)
			return
		}

		w.WriteHeader(204)
	}
}
