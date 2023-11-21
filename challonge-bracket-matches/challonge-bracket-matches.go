package challongebracketmatches

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

	// dealing with paginated response
	paginationLeft := true
	pageNumber := 1

	for paginationLeft {
		params := map[string]string{
			"state":         "in_progress",
			"created_after": date,
			"page":          strconv.Itoa(pageNumber),
			"per_page":      "25",
		}

		res, err := c.get(ctx, http.MethodGet, c.baseURL+"/tournaments.json", nil, params)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode))
		}

		defer res.Body.Close()

		var tournaments models.Tournaments
		err = json.NewDecoder(res.Body).Decode(&tournaments)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("%w. %s", err, http.StatusText(http.StatusInternalServerError))
		}

		if len(tournaments.Data) == 0 {
			paginationLeft = false
		} else {
			for _, tournament := range tournaments.Data {
				resMap[tournament.Id] = tournament.Attributes.GameName
			}

			pageNumber++
		}
	}

	return resMap, nil
}

func (c *customClient) get(ctx context.Context, method, urlPath string, reqBody io.Reader, params map[string]string) (resp *http.Response, err error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, urlPath, reqBody)
	if err != nil {
		// gracefully handle error and pass along
		return nil, err
	}

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
