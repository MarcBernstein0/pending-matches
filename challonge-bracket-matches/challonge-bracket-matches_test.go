package challongebracketmatches

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/MarcBernstein0/pending-matches/models"
	"github.com/stretchr/testify/assert"
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
		// mock endpoint for get tournaments
		case "/tournaments.json":
			mockFetchTournamentEndpoint(w, r)
		// mock endpoint for get participants
		case "/tournaments/1234/participants.json":
			mockFetchParticipantEndpoint(w, r)
		case "/tournaments/112358/participants.json":
			mockFetchParticipantEndpoint(w, r)
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
			testName:      "response not ok, auth error",
			mockDate:      time.Now().Local().Format("2006-01-02"),
			mockFetchData: New(server.URL, "bad api key", http.DefaultClient, 5*time.Second),
			wantData:      nil,
			wantErr:       fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(http.StatusUnauthorized)),
		},
		{
			testName:      "response ok one tournament no pagination",
			mockDate:      "2023-07-16",
			mockFetchData: New(server.URL, "mock api key", http.DefaultClient, 5*time.Second),
			wantData: map[string]string{
				"1": "test",
			},
			wantErr: nil,
		},
		{
			testName:      "response ok multiple tournament no pagination",
			mockDate:      "2023-07-17",
			mockFetchData: New(server.URL, "mock api key", http.DefaultClient, 5*time.Second),
			wantData: map[string]string{
				"1": "test",
				"2": "test2",
			},
			wantErr: nil,
		},
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

func TestFetchParticipants(t *testing.T) {
	tt := []struct {
		testName      string
		mockFetchData FetchData
		inputData     struct {
			tournamentId   string
			tournamentGame string
		}
		wantData models.TournamentParticipants
		wantErr  error
	}{
		{
			testName:      "response not ok",
			mockFetchData: New(server.URL, "bad api key", http.DefaultClient, 5*time.Second),
			inputData: struct {
				tournamentId   string
				tournamentGame string
			}{
				tournamentId:   "1234",
				tournamentGame: "testGameName",
			},
			wantData: models.TournamentParticipants{},
			wantErr:  fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(http.StatusUnauthorized)),
		},
		{
			testName:      "data found no pagination",
			mockFetchData: New(server.URL, "mock api key", http.DefaultClient, 5*time.Second),
			inputData: struct {
				tournamentId   string
				tournamentGame string
			}{
				tournamentId:   "1234",
				tournamentGame: "test",
			},
			wantData: models.TournamentParticipants{
				GameName:     "test",
				TournamentID: "1234",
				Participant: map[string]string{
					"1": "testName1",
					"2": "testName2",
					"3": "testName3",
					"4": "testName4",
				},
			},
			wantErr: nil,
		},
		{
			testName:      "data found pagination",
			mockFetchData: New(server.URL, "mock api key", http.DefaultClient, 5*time.Second),
			inputData: struct {
				tournamentId   string
				tournamentGame string
			}{
				tournamentId:   "112358",
				tournamentGame: "test",
			},
			wantData: models.TournamentParticipants{
				GameName:     "test",
				TournamentID: "112358",
				Participant: map[string]string{
					"1": "testName1",
					"2": "testName2",
					"3": "testName3",
					"4": "testName4",
					"5": "testName5",
					"6": "testName6",
					"7": "testName7",
					"8": "testName8",
				},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// t.Parallel()

			gotData, gotErr := tc.mockFetchData.FetchParticipants(context.Background(), tc.inputData.tournamentId, tc.inputData.tournamentGame)
			assert.Equal(t, tc.wantData, gotData)
			if tc.wantErr != nil {
				assert.EqualError(t, gotErr, tc.wantErr.Error())
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}

// helper functions
func testApiKeyAuth(apiKey string) bool {
	return apiKey == MOCK_API_KEY
}

func readJsonFile(filename string) ([]byte, error) {
	jsonFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	return byteValue, err

}

// mock endpoints
func mockFetchTournamentEndpoint(w http.ResponseWriter, r *http.Request) {
	emptyReturn, _ := readJsonFile("./mock-api-responses/mock-tournament-response-empty.json")

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
		if page > 1 {
			w.Write(emptyReturn)
			return
		}
		byteValue, _ := readJsonFile("./mock-api-responses/mock-tournament-response.json")
		w.Write(byteValue)
	}

	// no pagination, multi-tournaments
	if date == "2023-07-17" {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page > 1 {
			w.Write(emptyReturn)
			return
		}
		byteValue, _ := readJsonFile("./mock-api-responses/mock-tournament-multi-response.json")
		w.Write(byteValue)
	}

	// pagination, multi-tournaments
	if date == "2023-07-18" {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 4 {
			w.Write(emptyReturn)
		}
		if page == 1 {
			byteValue, _ := readJsonFile("./mock-api-responses/mock-tournament-multi-response.json")
			w.Write(byteValue)
		} else if page == 2 {
			byteValue, _ := readJsonFile("./mock-api-responses/mock-tournament-multi-response-page2.json")
			w.Write(byteValue)
		} else if page == 3 {
			byteValue, _ := readJsonFile("./mock-api-responses/mock-tournament-multi-response-page3.json")
			w.Write(byteValue)
		}
	}

}

func mockFetchParticipantEndpoint(w http.ResponseWriter, r *http.Request) {
	emptyReturn, _ := readJsonFile("./mock-api-responses/mock-tournament-response-empty.json")

	apiKey := r.Header.Get("Authorization")
	if !testApiKeyAuth(apiKey) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

	if strings.Contains(r.URL.Path, "1234") {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page > 1 {
			w.Write(emptyReturn)
			return
		}
		byteValue, _ := readJsonFile("./mock-api-responses/mock-participant-response.json")
		w.Write(byteValue)
	}
	if strings.Contains(r.URL.Path, "112358") {
		fmt.Println("multi-page-print")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 3 {
			w.Write(emptyReturn)
		}
		if page == 1 {
			byteValue, _ := readJsonFile("./mock-api-responses/mock-participant-response.json")
			w.Write(byteValue)
		}
		if page == 2 {
			byteValue, _ := readJsonFile("./mock-api-responses/mock-participant-response-page2.json")
			w.Write(byteValue)
		}
	}

}
