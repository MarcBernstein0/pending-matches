package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	challongebracketmatches "github.com/MarcBernstein0/pending-matches/challonge-bracket-matches"
	"github.com/MarcBernstein0/pending-matches/challonge-bracket-matches/cache"
	"github.com/MarcBernstein0/pending-matches/route"
)

func main() {
	slog.Info("Test slog")
	port, present := os.LookupEnv("PORT")
	if !present {
		port = "8080"
	}

	apiKey, present := os.LookupEnv("API_KEY")
	if !present {
		log.Fatalf("api_key not provided in env")
	}
	cacheTimerString, present := os.LookupEnv("CACHE_TIMER")
	if !present {
		cacheTimerString = "3"
	}
	cacheTimer, err := strconv.Atoi(cacheTimerString)
	if err != nil {
		log.Fatalf("cacheTimer could not be read properly\n%s", err)
	}

	cacheLastClearTimerString, present := os.LookupEnv("CACHE_CLEAR_TIMER")
	if !present {
		cacheLastClearTimerString = "5"
	}
	cacheClearTimer, err := strconv.Atoi(cacheLastClearTimerString)
	if err != nil {
		log.Fatalf("cacheTimer could not be read properly\n%s", err)
	}

	customClient := challongebracketmatches.New("https://api.challonge.com/v2.1", apiKey, http.DefaultClient, 20*time.Minute)
	customCache := cache.NewCache(time.Duration(cacheTimer)*time.Minute, time.Duration(cacheClearTimer)*time.Hour)

	r := route.RouterSetup(customClient, customCache)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
