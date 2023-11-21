package cache

import (
	"fmt"
	"time"

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

func (c *Cache) UpdateCache(listTournamentParticipants []models.TournamentParticipants, date string) {
	fmt.Println("Cache is updating")
	c.data[date] = cacheData{
		tournamentsAndParticipants: listTournamentParticipants,
		timeStamp:                  time.Now(),
	}
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
