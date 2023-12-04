package route

import (
	"net/http"

	challongebracketmatches "github.com/MarcBernstein0/pending-matches/challonge-bracket-matches"
	"github.com/MarcBernstein0/pending-matches/challonge-bracket-matches/cache"
	"github.com/go-chi/chi/v5"
)

func RouterSetup(fetchData challongebracketmatches.FetchData, cache *cache.Cache) *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`
			{
				"status": "UP"
			}
			`))
		})
		r.Get("/matches", GetMatches(fetchData, cache))
	})

	return r
}
