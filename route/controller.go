package route

import (
	"fmt"
	"net/http"

	challongebracketmatches "github.com/MarcBernstein0/pending-matches/challonge-bracket-matches"
	"github.com/MarcBernstein0/pending-matches/challonge-bracket-matches/cache"
)

func GetMatches(fetchData challongebracketmatches.FetchData, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Get Match called")
	}
}
