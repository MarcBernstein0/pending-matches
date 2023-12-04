package cache

import (
	"log/slog"
	"sync"
	"time"

	challongebracketmatches "github.com/MarcBernstein0/pending-matches/challonge-bracket-matches"
	"github.com/MarcBernstein0/pending-matches/models"
)

type cacheData struct {
	tournamentsAndParticipants []models.TournamentParticipants
	timeStamp                  time.Time
}

type Cache struct {
	data             map[string]cacheData
	updateCacheTimer time.Duration
	clearCacheTimer  time.Duration
	lastClearCache   time.Time
	logger           *slog.Logger
}

func NewCache(cacheTimer, clearCacheTimer time.Duration, logger *slog.Logger) *Cache {
	return &Cache{
		data:             map[string]cacheData{},
		updateCacheTimer: cacheTimer,
		clearCacheTimer:  clearCacheTimer,
		lastClearCache:   time.Now(),
		logger:           logger,
	}
}

func (c *Cache) UpdateCache(date string, fetchData challongebracketmatches.FetchData) error {
	c.logger.Info("Fetching tournaments") // TODO: Replace print with logging
	tournaments, err := fetchData.FetchTournaments(date)
	if err != nil {
		return err
	}

	if len(tournaments) == 0 {
		return nil
	}

	c.logger.Info("Fetching participants") // TODO: Replace print with logging
	listTournamentParticipants, err := c.getParticipantsConcurrently(tournaments, fetchData)
	if err != nil {
		return err
	}

	c.logger.Info("Cache is updating") // TODO: Replace print with logging
	c.data[date] = cacheData{
		tournamentsAndParticipants: listTournamentParticipants,
		timeStamp:                  time.Now(),
	}
	return nil
}

func (c *Cache) GetData(date string) []models.TournamentParticipants {
	c.logger.Info("Getting data from cache") // TODO: Replace print with logging
	return c.data[date].tournamentsAndParticipants
}

func (c *Cache) ShouldUpdate(date string) bool {
	if data, ok := c.data[date]; ok {
		timeSince := time.Since(data.timeStamp)
		return timeSince >= c.updateCacheTimer
	}
	return true
}

func (c *Cache) IsCacheEmptyAtDate(date string) bool {
	if data, ok := c.data[date]; ok {
		return len(data.tournamentsAndParticipants) == 0
	}
	return true
}

func (c *Cache) ShouldClearCacheData() bool {
	timeSince := time.Since(c.lastClearCache)
	return timeSince >= c.clearCacheTimer
}

func (c *Cache) ClearCache() {
	c.data = map[string]cacheData{}
	c.lastClearCache = time.Now()
}

func (c *Cache) getParticipantsConcurrently(tournaments map[string]string, fetchData challongebracketmatches.FetchData) ([]models.TournamentParticipants, error) {
	var tournamentParticipants []models.TournamentParticipants

	chanResponse := make(chan struct {
		tournamentParticipant *models.TournamentParticipants
		err                   error
	})
	var wg sync.WaitGroup
	for key, val := range tournaments {
		wg.Add(1)
		go func(tournamentId, tournamentGame string, chanResponse chan struct {
			tournamentParticipant *models.TournamentParticipants
			err                   error
		}, wg *sync.WaitGroup) {
			defer wg.Done()
			participants, err := fetchData.FetchParticipants(tournamentId, tournamentGame)
			if err != nil {
				chanResponse <- struct {
					tournamentParticipant *models.TournamentParticipants
					err                   error
				}{
					tournamentParticipant: nil,
					err:                   err,
				}
				return
			}
			chanResponse <- struct {
				tournamentParticipant *models.TournamentParticipants
				err                   error
			}{
				tournamentParticipant: &participants,
				err:                   nil,
			}
		}(key, val, chanResponse, &wg)
	}

	go func() {
		wg.Wait()
		close(chanResponse)
	}()

	for getParticipantResult := range chanResponse {
		if getParticipantResult.err != nil {
			return nil, getParticipantResult.err
		}
		tournamentParticipants = append(tournamentParticipants, *getParticipantResult.tournamentParticipant)
	}

	return tournamentParticipants, nil
}
