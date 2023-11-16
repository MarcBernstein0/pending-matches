package challongebracketmatches

import (
	"errors"
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
		baseURL string
		client  *http.Client
		config  struct {
			apiKey string
		}
		contextTimeout time.Duration
	}
)

func New(baseURL, apiKey string, client *http.Client, contextTimeout time.Duration) customClient {
	return customClient{
		baseURL: baseURL,
		client:  client,
		config: struct{ apiKey string }{
			apiKey: apiKey,
		},
		contextTimeout: contextTimeout,
	}
}
