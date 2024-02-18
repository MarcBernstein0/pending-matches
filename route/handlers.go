package route

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"

	challongebracketmatches "github.com/MarcBernstein0/pending-matches/challonge-bracket-matches"
	"github.com/MarcBernstein0/pending-matches/challonge-bracket-matches/cache"
	"github.com/MarcBernstein0/pending-matches/models"
	"github.com/go-chi/httplog/v2"
)

func GetMatches(fetchData challongebracketmatches.FetchData, cache *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get logger
		logger := httplog.LogEntry(r.Context())

		// set json response header
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// check if cache should be cleared
		if cache.ShouldClearCacheData() {
			cache.ClearCache()
		}

		requestValues, err := models.CreateRequestValues(r.URL.Query())
		if err != nil {
			requestQueryParamErr := ErrorBadRequest(err.Error(), err)
			requestQueryParamErr.LogError(logger)
			requestQueryParamErr.JSONError(w)
			return
		}

		// Get tournaments and participants
		var tournamentsAndParticipants []models.TournamentParticipants
		// check if cache is empty or time limit has been exceeded
		if cache.IsCacheEmptyAtDate(requestValues.Date) || cache.ShouldUpdate(requestValues.Date) {
			// update cache
			err := cache.UpdateCache(requestValues.Date, fetchData)
			if err != nil {
				cacheUpdateError := ErrorInternal("Error in getting tournament data", err)
				cacheUpdateError.LogError(logger)
				cacheUpdateError.JSONError(w)
				return
			}
		}

		tournamentsAndParticipants = cache.GetData(requestValues.Date, requestValues.GameList)

		matches, err := getMatchesConcurrently(tournamentsAndParticipants, fetchData)
		if err != nil {
			getMatchesErr := ErrorInternal("Error in getting match data", err)
			getMatchesErr.LogError(logger)
			getMatchesErr.JSONError(w)
			return
		}

		json.NewEncoder(w).Encode(matches)
	}
}

func getMatchesConcurrently(tournamentsAndParticipants []models.TournamentParticipants, fetchData challongebracketmatches.FetchData) ([]models.TournamentMatches, error) {
	matches := []models.TournamentMatches{}

	chanResponse := make(chan struct {
		tournamentMatches *models.TournamentMatches
		err               error
	})
	var wg sync.WaitGroup
	for _, elem := range tournamentsAndParticipants {
		wg.Add(1)
		go func(tournament models.TournamentParticipants, chanResponse chan struct {
			tournamentMatches *models.TournamentMatches
			err               error
		}) {
			defer wg.Done()
			match, err := fetchData.FetchMatches(tournament)
			if err != nil {
				chanResponse <- struct {
					tournamentMatches *models.TournamentMatches
					err               error
				}{
					tournamentMatches: nil,
					err:               err,
				}
				return
			}
			chanResponse <- struct {
				tournamentMatches *models.TournamentMatches
				err               error
			}{
				tournamentMatches: &match,
				err:               nil,
			}
		}(elem, chanResponse)
	}

	go func() {
		wg.Wait()
		close(chanResponse)
	}()

	for getMatchesResult := range chanResponse {
		if getMatchesResult.err != nil {
			return nil, getMatchesResult.err
		}
		matches = append(matches, *getMatchesResult.tournamentMatches)
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].GameName < matches[j].GameName
	})

	return matches, nil
}
