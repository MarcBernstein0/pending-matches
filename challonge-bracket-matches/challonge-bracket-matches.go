package challongebracketmatches

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/MarcBernstein0/pending-matches/models"
)

var (
	ErrResponseNotOK error = errors.New("response not ok")
	ErrServerProblem error = errors.New("server error")
	ErrNoData        error = errors.New("no data found")
)

type (
	customClient struct {
		baseURL        string
		client         *http.Client
		apiKey         string
		contextTimeout time.Duration
	}

	FetchData interface {
		// FetchTournaments fetch all tournaments created after a specific date
		// GET https://api.challonge.com/v2.1/tournaments.json?page=1&per_page=25
		FetchTournaments(ctx context.Context, date string) (map[string]string, error)
	}
)

func New(baseURL, apiKey string, client *http.Client, contextTimeout time.Duration) *customClient {
	return &customClient{
		baseURL:        baseURL,
		client:         client,
		apiKey:         apiKey,
		contextTimeout: contextTimeout,
	}
}

// Return map of type int -> string where int is the tournamentId and string is the game name
func (c *customClient) FetchTournaments(ctx context.Context, date string) (map[string]string, error) {

	resMap := make(map[string]string)

	params := map[string]string{
		"state":         "in_progress",
		"created_after": date,
	}

	ctx, cancelCtx := context.WithTimeout(ctx, c.contextTimeout)
	defer cancelCtx()

	requestURL := c.baseURL + "/tournaments.json"
	fmt.Println(requestURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.get(req, params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode))
	}
	fmt.Println(res)

	defer res.Body.Close()

	var tournaments models.Tournaments
	err = json.NewDecoder(res.Body).Decode(&tournaments)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("%w. %s", err, http.StatusText(http.StatusInternalServerError))
	}
	fmt.Printf("tournaments value %+v\n", tournaments)
	for _, tournament := range tournaments.Data {
		resMap[tournament.Id] = tournament.Attributes.GameName
	}

	return resMap, nil
}

func (c *customClient) get(req *http.Request, params map[string]string) (resp *http.Response, err error) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization-Type", "v1")
	req.Header.Add("Authorization", c.apiKey)

	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	return c.client.Do(req)
}
