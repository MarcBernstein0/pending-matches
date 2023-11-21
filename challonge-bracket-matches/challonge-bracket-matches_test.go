package challongebracketmatches

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var server *httptest.Server

const MOCK_API_KEY = "mock api key"

var mockReturnValue string = `
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
}`

var mockReturnValueMultiTournaments string = `
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
}`

var mockReturnValueMultiTournaments2 string = `
{
	"data": [
		{
			"id": "3",
			"type": "type",
			"attributes": {
				"tournament_type": "tournament_type",
				"name": "testName",
				"state": "state",
				"game_name": "test3"
			}
		},
		{
			"id": "4",
			"type": "type",
			"attributes": {
				"tournament_type": "tournament_type",
				"name": "testName",
				"state": "state",
				"game_name": "test4"
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
}`

var mockReturnValueMultiTournaments3 string = `
{
	"data": [
		{
			"id": "5",
			"type": "type",
			"attributes": {
				"tournament_type": "tournament_type",
				"name": "testName",
				"state": "state",
				"game_name": "test5"
			}
		},
		{
			"id": "6",
			"type": "type",
			"attributes": {
				"tournament_type": "tournament_type",
				"name": "testName",
				"state": "state",
				"game_name": "test6"
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
}`

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
		// {
		// 	testName:      "response not ok, auth error",
		// 	mockDate:      time.Now().Local().Format("2006-01-02"),
		// 	mockFetchData: New(server.URL, "bad api key", http.DefaultClient, 5*time.Second),
		// 	wantData:      nil,
		// 	wantErr:       fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(http.StatusUnauthorized)),
		// },
		// {
		// 	testName:      "response ok one tournament no pagination",
		// 	mockDate:      "2023-07-16",
		// 	mockFetchData: New(server.URL, "mock api key", http.DefaultClient, 5*time.Second),
		// 	wantData: map[string]string{
		// 		"1": "test",
		// 	},
		// 	wantErr: nil,
		// },
		// {
		// 	testName:      "response ok multiple tournament no pagination",
		// 	mockDate:      "2023-07-17",
		// 	mockFetchData: New(server.URL, "mock api key", http.DefaultClient, 5*time.Second),
		// 	wantData: map[string]string{
		// 		"1": "test",
		// 		"2": "test2",
		// 	},
		// 	wantErr: nil,
		// },
		{
			testName:      "response ok multiple tournament and pagination",
			mockDate:      "2023-07-18",
			mockFetchData: New(server.URL, "mock api key", http.DefaultClient, 5*time.Second),
			wantData: map[string]string{
				"1": "test",
				"2": "test2",
				"3": "test3",
				"4": "test4",
				"5": "test5",
				"6": "test6",
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

	// no pagination, one tournaments
	if date == "2023-07-16" {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 1 {
			emptyReturn := `{
				"data": [],
				"included": [],
				"meta": {
					"count": 36
				},
				"links": {
					"self": "https://api.challonge.com/v2.1/tournaments.json?page=2&per_page=25&state=in_progress&created_after=2023-07-22",
					"next": "https://api.challonge.com/v2.1/tournaments.json?page=3&per_page=25",
					"prev": "https://api.challonge.com/v2.1/tournaments.json?page=1&per_page=25"
				}
			}`
			w.Write([]byte(emptyReturn))
			return
		}
		byteValue := []byte(mockReturnValue)
		w.Write(byteValue)
	}

	// no pagination, multi-tournaments
	if date == "2023-07-17" {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 1 {
			emptyReturn := `{
				"data": [],
				"included": [],
				"meta": {
					"count": 36
				},
				"links": {
					"self": "https://api.challonge.com/v2.1/tournaments.json?page=2&per_page=25&state=in_progress&created_after=2023-07-22",
					"next": "https://api.challonge.com/v2.1/tournaments.json?page=3&per_page=25",
					"prev": "https://api.challonge.com/v2.1/tournaments.json?page=1&per_page=25"
				}
			}`
			w.Write([]byte(emptyReturn))
			return
		}
		byteValue := []byte(mockReturnValueMultiTournaments)
		w.Write(byteValue)
	}

	// pagination, multi-tournaments
	if date == "2023-07-18" {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 4 {
			emptyReturn := `{
				"data": [],
				"included": [],
				"meta": {
					"count": 36
				},
				"links": {
					"self": "https://api.challonge.com/v2.1/tournaments.json?page=2&per_page=25&state=in_progress&created_after=2023-07-22",
					"next": "https://api.challonge.com/v2.1/tournaments.json?page=3&per_page=25",
					"prev": "https://api.challonge.com/v2.1/tournaments.json?page=1&per_page=25"
				}
			}`
			w.Write([]byte(emptyReturn))
			return
		}
		if page == 1 {
			byteValue := []byte(mockReturnValueMultiTournaments)
			w.Write(byteValue)
		} else if page == 2 {
			byteValue := []byte(mockReturnValueMultiTournaments2)
			w.Write(byteValue)
		} else if page == 3 {
			byteValue := []byte(mockReturnValueMultiTournaments3)
			w.Write(byteValue)
		}
	}

}

func testApiKeyAuth(apiKey string) bool {
	return apiKey == MOCK_API_KEY
}
