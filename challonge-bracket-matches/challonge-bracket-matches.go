package challongebracketmatches

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
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
		FetchTournaments(ctx context.Context, date string) (map[int]string, error)
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
func (c *customClient) FetchTournaments(ctx context.Context, date string) (map[int]string, error) {
	params := map[string][]string{
		"state":         {"in_progress"},
		"created_after": {date},
	}

	ctx, cancelCtx := context.WithTimeout(ctx, c.contextTimeout)
	defer cancelCtx()

	requestURL := c.baseURL + "/tournaments.json"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.get(req, params)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		fmt.Println(err)
		return nil, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode))
	}

	defer res.Body.Close()
	fmt.Println(res.Body)

	return nil, nil
}

func (c *customClient) get(req *http.Request, params map[string]string) (resp *http.Response, err error) {
	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization-Type", "v1")
	req.Header.Add("Authorization", c.apiKey)

	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL)
	fmt.Println(req.Header)
	fmt.Println(req)

	return c.client.Do(req)
}
