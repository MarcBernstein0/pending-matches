package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	challongebracketmatches "github.com/MarcBernstein0/pending-matches/challonge-bracket-matches"
	"github.com/MarcBernstein0/pending-matches/challonge-bracket-matches/cache"
	"github.com/MarcBernstein0/pending-matches/models"
)

func GetMatches(fetchData challongebracketmatches.FetchData, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// set json response header
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// check if cache should be cleared
		if cache.ShouldClearCacheData() {
			cache.ClearCache()
		}

		matches := make([]models.TournamentMatches, 0)

		dateStr := r.URL.Query().Get("date")
		if dateStr == "" {
			dateErr := errors.New("date query parameter not provided")
			noDateProvidedErr := ErrorBadRequest(dateErr.Error(), dateErr)
			noDateProvidedErr.LogError()
			noDateProvidedErr.JSONError(w)
			return
		}
		if _, err := time.Parse("2006-01-02", dateStr); err != nil {
			dateStrNotFormattedProperly := ErrorBadRequest("Date query parameter not formatted properly. Expect formatting YYYY-MM-DD", err)
			dateStrNotFormattedProperly.LogError()
			dateStrNotFormattedProperly.JSONError(w)
			return
		}
		fmt.Println("Get Match called", dateStr)

		// Get tournaments and participants
		var tournamentsAndParticipants []models.TournamentParticipants
		// check if cache is empty or time limit has been exceeded
		if cache.IsCacheEmptyAtDate(dateStr) || cache.ShouldUpdate(dateStr) {
			// update cache
			err := cache.UpdateCache(dateStr, fetchData)
			if err != nil {
				cacheUpdateError := ErrorInternal("Error in getting tournament data", err)
				cacheUpdateError.LogError()
				cacheUpdateError.JSONError(w)
				return
			}
		}

		tournamentsAndParticipants = cache.GetData(dateStr)

		for _, elem := range tournamentsAndParticipants {
			match, err := fetchData.FetchMatches(elem)
			if err != nil {
				getMatchesErr := ErrorInternal("Error in getting match data", err)
				getMatchesErr.LogError()
				getMatchesErr.JSONError(w)
				return
			}
			matches = append(matches, match)
		}
		// fmt.Println(tournamentsAndParticipants)
		json.NewEncoder(w).Encode(matches)
	}
}
