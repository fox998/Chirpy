package config

import (
	"database/sql"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/fox998/Chirpy/internal/config/env"
	"github.com/fox998/Chirpy/internal/database"
)

type ApiConfig struct {
	fileserverHits atomic.Int32
	Db             *database.Queries
	Env            env.EnvParams
}

func NewApiConfig() *ApiConfig {

	envParams := env.GetEnvParams()
	db, err := sql.Open("postgres", envParams.DbUrl)
	if err != nil {
		log.Fatalf("failed open DB with %s. %s", envParams.DbUrl, err.Error())
	}

	return &ApiConfig{
		Db:  database.New(db),
		Env: envParams,
	}
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cnf *ApiConfig) GetHits() int32 {
	return cnf.fileserverHits.Load()
}

func (cnf *ApiConfig) ResetHits() {
	cnf.fileserverHits.Store(0)
}
