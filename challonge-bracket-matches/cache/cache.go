package cache

import (
	"context"
	"fmt"
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
}

func NewCache(cacheTimer, clearCacheTimer time.Duration) *Cache {
	return &Cache{
		data:             map[string]cacheData{},
		updateCacheTimer: cacheTimer,
		clearCacheTimer:  clearCacheTimer,
	}
}

func (c *Cache) UpdateCache(date string, fetchData challongebracketmatches.FetchData) error {
	listTournamentParticipants := []models.TournamentParticipants{}
	fmt.Println("Fetching tournaments")
	tournaments, err := fetchData.FetchTournaments(context.Background(), date)
	if err != nil {
		return err
	}
	fmt.Println("Fetching participants")
	for key, val := range tournaments {
		participants, err := fetchData.FetchParticipants(context.Background(), key, val)
		if err != nil {
			return err
		}
		listTournamentParticipants = append(listTournamentParticipants, participants)
	}
	fmt.Println("Cache is updating")
	c.data[date] = cacheData{
		tournamentsAndParticipants: listTournamentParticipants,
		timeStamp:                  time.Now(),
	}
	fmt.Printf("cache: %+v\n", c)
	return nil
}

func (c *Cache) GetData(date string) []models.TournamentParticipants {
	fmt.Println("Getting data from cache")
	return c.data[date].tournamentsAndParticipants
}

func (c *Cache) ShouldUpdate(date string) bool {
	if data, ok := c.data[date]; ok {
		timeSince := time.Since(data.timeStamp)
		return timeSince >= c.updateCacheTimer
	}
	return true
}

func (c *Cache) IsCacheEmptyDate(date string) bool {
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
