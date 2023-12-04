package route

import (
	"log/slog"
	"net/http"

	challongebracketmatches "github.com/MarcBernstein0/pending-matches/challonge-bracket-matches"
	"github.com/MarcBernstein0/pending-matches/challonge-bracket-matches/cache"
	"github.com/go-chi/chi/v5"
)

func RouterSetup(fetchData challongebracketmatches.FetchData, cache *cache.Cache, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		logger.Info("GET /v1/health route setup")
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
			{
				"status": "UP"
			}
			`))
		})
		logger.Info("GET /v1/matches route setup")
		r.Get("/matches", GetMatches(fetchData, cache))
	})

	return r
}
