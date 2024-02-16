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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httplog/v2"
)

func main() {
	logger := httplog.NewLogger("match-display", httplog.Options{
		// JSON:             true,
		LogLevel: slog.LevelInfo,
		Concise:  true,
		// RequestHeaders:   true,
		MessageFieldName: "message",
		// TimeFieldFormat: time.RFC850,
		Tags: map[string]string{
			"version": "v1.0-81aa4244d9fc8076a",
			"env":     "dev",
		},
		QuietDownRoutes: []string{
			"/",
			"/ping",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})
	slog.SetDefault(logger.Logger)

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
	customCache := cache.NewCache(time.Duration(cacheTimer)*time.Minute, time.Duration(cacheClearTimer)*time.Hour, logger.Logger)

	// chi service
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			httplog.LogEntrySetField(ctx, "user", slog.StringValue("user1"))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	api := route.RouterSetup(customClient, customCache)

	r.Mount("/", api)
	logger.Info("pending match server started")
	log.Fatal(http.ListenAndServe(":"+port, r))
}
