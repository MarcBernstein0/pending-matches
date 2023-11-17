package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	challongebracketmatches "github.com/MarcBernstein0/pending-matches/challonge-bracket-matches"
)

func main() {
	// r := chi.NewRouter()
	// r.Use(middleware.Logger)
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello World!"))
	// })
	// log.Fatal(http.ListenAndServe(":8080", r))

	// port, present := os.LookupEnv("PORT")
	// if !present {
	// 	port = "8080"
	// }

	apiKey, present := os.LookupEnv("API_KEY")
	if !present {
		log.Fatalf("api_key not provided in env")
	}

	customClient := challongebracketmatches.New("https://api.challonge.com/v2.1", apiKey, http.DefaultClient, 20*time.Minute)
	fmt.Println(customClient.FetchTournaments(context.Background(), "2023-07-04"))
}
