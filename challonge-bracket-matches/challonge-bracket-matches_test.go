package challongebracketmatches

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var server *httptest.Server

const MOCK_API_KEY = "mock api key"

func TestMain(m *testing.M) {
	fmt.Println("Mock Server")
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		trimPath := strings.TrimSpace(r.URL.Path)
		fmt.Println("trim path", trimPath)
		switch trimPath {
		case "/tournaments.json":
			mockFetchTournamentEndpoint(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	fmt.Println("run tests")
	m.Run()
}

func TestCreateCustomClient(t *testing.T) {
	// Given
	givenCustomClient := &customClient{
		baseURL:        "testEndpoint",
		client:         http.DefaultClient,
		apiKey:         "1234567890",
		contextTimeout: 20,
	}
	// When
	res := New("testEndpoint", "1234567890", http.DefaultClient, 20)
	// Then
	require.Equal(t, givenCustomClient, res)
}

func TestFetchTournaments(t *testing.T) {
	// Given
	tt := []struct {
		testName      string
		mockDate      string
		mockFetchData FetchData
		wantData      map[string]string
		wantErr       error
	}{
		{
			testName:      "response ok one tournament no pagination",
			mockDate:      time.Now().Local().Format("2006-01-02"),
			mockFetchData: New(server.URL, "mock api key", http.DefaultClient, 5*time.Second),
			wantData: map[string]string{
				"1": "test",
				"2": "test2",
			},
			wantErr: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// When
			gotData, gotErr := tc.mockFetchData.FetchTournaments(context.Background(), tc.mockDate)

			//Then
			require.Equal(t, tc.wantData, gotData)
			if tc.wantErr != nil {
				require.EqualError(t, gotErr, tc.wantErr.Error())
			} else {
				require.NoError(t, gotErr)
			}
		})
	}
}

// mockFunctions
func mockFetchTournamentEndpoint(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Authorization")
	if !testApiKeyAuth(apiKey) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

	date := r.URL.Query().Get("created_after")
	if date == "2022-07-16" {
		w.Write([]byte("[]"))
		return
	}

	mockReturnValue := `
	{
		"data": [
			{
				"id": "1",
				"type": "type",
				"attributes": {
					"tournament_type": "tournament_type",
					"name": "testName",
					"state": "state",
					"game_name": "test"
				}
			},
			{
				"id": "2",
				"type": "type",
				"attributes": {
					"tournament_type": "tournament_type",
					"name": "testName",
					"state": "state",
					"game_name": "test2"
				}
			}
		],
		"included": [
			{
				"id": "blazblue-central-fiction",
				"type": "game",
				"attributes": {
					"name": "BlazBlue: Central Fiction",
					"aliases": [
						"blazeblue central fiction"
					],
					"verified": true
				}
			}
		]
	}
	`
	byteValue := []byte(mockReturnValue)
	w.Write(byteValue)
}

func testApiKeyAuth(apiKey string) bool {
	return apiKey == MOCK_API_KEY
}
