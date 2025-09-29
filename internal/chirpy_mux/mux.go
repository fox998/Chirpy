package chirpy_mux

import (
	"net/http"

	"github.com/fox998/Chirpy/internal/config"
	"github.com/fox998/Chirpy/internal/handlers"
)

type ChirpyMux struct {
	*http.ServeMux
}

func (mux *ChirpyMux) setupHandlers(config *config.ApiConfig) {
	{
		appHandler := config.MiddlewareMetricsInc(http.FileServer(http.Dir(".")))
		mux.Handle("/app/", http.StripPrefix("/app/", appHandler))
	}

	{
		// mux.HandleFunc("GET /api/healthz", handlers.Healthz)
		mux.HandleFunc("POST /api/users", handlers.Users(config))
		mux.HandleFunc("PUT /api/users", handlers.UsersUpdate(config))

		mux.HandleFunc("POST /api/login", handlers.Login(config))
		mux.HandleFunc("POST /api/refresh", handlers.Refresh(config))
		mux.HandleFunc("POST /api/revoke", handlers.Revoke(config))

		mux.HandleFunc("GET /api/chirps", handlers.AllChirps(config))
		mux.HandleFunc("GET /api/chirps/{id}", handlers.ChirpsById(config))
		mux.HandleFunc("DELETE /api/chirps/{id}", handlers.DeleteChirps(config))
		mux.HandleFunc("POST /api/chirps", handlers.Chirps(config))
	}

	{
		mux.HandleFunc("GET /admin/metrics", handlers.Metrics(config))
		mux.HandleFunc("POST /admin/reset", handlers.Reset(config))
	}
}

func CreateServerMux(config *config.ApiConfig) *ChirpyMux {

	mux := ChirpyMux{
		ServeMux: http.NewServeMux(),
	}

	mux.setupHandlers(config)
	return &mux
}
