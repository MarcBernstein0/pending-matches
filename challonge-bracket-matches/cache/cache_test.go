package cache

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	challongebracketmatches "github.com/MarcBernstein0/pending-matches/challonge-bracket-matches"
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
		switch trimPath {
		case "/tournaments.json":
			mockFetchTournamentEndpoint(w, r)
		case "/tournaments/1/participants.json":
			mockFetchParticipantsEndpoint(w, r)
		case "/tournaments/2/participants.json":
			mockFetchParticipantsEndpoint(w, r)
		case "/tournaments/3/participants.json":
			mockFetchParticipantsEndpoint(w, r)
		case "/tournaments/4/participants.json":
			mockFetchParticipantsEndpoint(w, r)
		case "/tournaments/5/participants.json":
			mockFetchParticipantsEndpoint(w, r)
		case "/tournaments/6/participants.json":
			mockFetchParticipantsEndpoint(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	fmt.Println("run tests")
	m.Run()

}

func TestCreateCache(t *testing.T) {
	// Given
	givenCache := &Cache{
		data:             map[string]cacheData{},
		updateCacheTimer: 5 * time.Minute,
		clearCacheTimer:  5 * time.Minute,
	}
	// When
	mockCache := NewCache(5*time.Minute, 5*time.Minute)
	// Then
	require.Equal(t, givenCache, mockCache)
}

func TestUpdateCache(t *testing.T) {
	mockCache := NewCache(5*time.Minute, 5*time.Hour)
	// Given
	tt := []struct {
		name          string
		date          string
		mockFetchData challongebracketmatches.FetchData
		wantData      []models.TournamentParticipants
		wantErr       error
	}{
		{
			name:          "response not ok",
			date:          "2022-07-16",
			mockFetchData: challongebracketmatches.New(server.URL, "bad api key", http.DefaultClient, 5*time.Second),
			wantData:      nil,
			wantErr:       fmt.Errorf("%w. %s", challongebracketmatches.ErrResponseNotOK, http.StatusText(http.StatusUnauthorized)),
		},
		{
			name:          "response not ok but no tournaments",
			date:          "2022-07-16",
			mockFetchData: challongebracketmatches.New(server.URL, MOCK_API_KEY, http.DefaultClient, 5*time.Second),
			wantData:      []models.TournamentParticipants{},
			wantErr:       nil,
		},
		{
			name:          "response ok",
			date:          "2023-11-25",
			mockFetchData: challongebracketmatches.New(server.URL, MOCK_API_KEY, http.DefaultClient, 5*time.Second),
			wantData: []models.TournamentParticipants{
				{
					GameName:     "test",
					TournamentID: "1234",
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
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// When
			gotErr := mockCache.UpdateCache(tc.date, tc.mockFetchData)
			fmt.Println(gotErr)

			// Then
			if tc.wantErr != nil {
				assert.EqualError(t, gotErr, tc.wantErr.Error())
			} else {
				assert.ElementsMatch(t, tc.wantData, mockCache.GetData(tc.date))
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

func mockFetchTournamentEndpoint(w http.ResponseWriter, r *http.Request) {
	emptyReturn, _ := readJsonFile("../mock-api-responses/mock-tournament-response-empty.json")

	apiKey := r.Header.Get("Authorization")
	if !testApiKeyAuth(apiKey) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)

	date := r.URL.Query().Get("created_after")
	if date == "2022-07-16" {
		w.Write(emptyReturn)
		return
	}

	if date == "2023-11-25" {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 4 {
			w.Write(emptyReturn)
		}
		if page == 1 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-tournament-multi-response.json")
			w.Write(byteValue)
		} else if page == 2 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-tournament-multi-response-page2.json")
			w.Write(byteValue)
		} else if page == 3 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-tournament-multi-response-page3.json")
			w.Write(byteValue)
		}
	}
}

func mockFetchParticipantsEndpoint(w http.ResponseWriter, r *http.Request) {
	emptyReturn, _ := readJsonFile("../mock-api-responses/mock-tournament-response-empty.json")

	apiKey := r.Header.Get("Authorization")
	if !testApiKeyAuth(apiKey) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	if strings.Contains(r.URL.Path, "1") {
		// fmt.Println("multi-page-print")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 3 {
			w.Write(emptyReturn)
		}
		if page == 1 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-participant-response.json")
			w.Write(byteValue)
		}
		if page == 2 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-participant-response-page2.json")
			w.Write(byteValue)
		}
	}
	if strings.Contains(r.URL.Path, "2") {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page > 1 {
			w.Write(emptyReturn)
			return
		}
		byteValue, _ := readJsonFile("../mock-api-responses/mock-participant-response.json")
		w.Write(byteValue)
	}
	if strings.Contains(r.URL.Path, "3") {
		// fmt.Println("multi-page-print")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 3 {
			w.Write(emptyReturn)
		}
		if page == 1 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-participant-response.json")
			w.Write(byteValue)
		}
		if page == 2 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-participant-response-page2.json")
			w.Write(byteValue)
		}
	}
	if strings.Contains(r.URL.Path, "4") {
		// fmt.Println("multi-page-print")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 2 {
			w.Write(emptyReturn)
		}
		if page == 1 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-participant-response.json")
			w.Write(byteValue)
		}
	}
	if strings.Contains(r.URL.Path, "5") {
		// fmt.Println("multi-page-print")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 3 {
			w.Write(emptyReturn)
		}
		if page == 1 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-participant-response.json")
			w.Write(byteValue)
		}
		if page == 2 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-participant-response-page2.json")
			w.Write(byteValue)
		}
	}
	if strings.Contains(r.URL.Path, "6") {
		// fmt.Println("multi-page-print")
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page >= 3 {
			w.Write(emptyReturn)
		}
		if page == 1 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-participant-response.json")
			w.Write(byteValue)
		}
		if page == 2 {
			byteValue, _ := readJsonFile("../mock-api-responses/mock-participant-response-page2.json")
			w.Write(byteValue)
		}
	}
}

// func TestShouldUpdate(t *testing.T) {
// 	t.Run("Test that method returns true when timer has exceeded limit", func(t *testing.T) {
// 		// Given
// 		mockCache := NewCache(2*time.Microsecond, 2*time.Microsecond)
// 		mockCache.UpdateCache([]models.TournamentParticipants{
// 			{
// 				GameName:     "Guilty Gear -Strive-",
// 				TournamentID: 10879090,
// 				Participant: map[int]string{
// 					166014671: "test",
// 					166014672: "test2",
// 					166014673: "test3",
// 					166014674: "test4",
// 				},
// 			},
// 		}, "2006-01-02")
// 		// When
// 		time.Sleep(5 * time.Millisecond)
// 		// Then
// 		assert.Equal(t, true, mockCache.ShouldUpdate("2006-01-02"))
// 	})

// 	t.Run("Test that method returns false when timer has note exceeded limit", func(t *testing.T) {
// 		// Given
// 		mockCache := NewCache(5*time.Millisecond, 5*time.Millisecond)
// 		mockCache.UpdateCache([]models.TournamentParticipants{
// 			{
// 				GameName:     "Guilty Gear -Strive-",
// 				TournamentID: 10879090,
// 				Participant: map[int]string{
// 					166014671: "test",
// 					166014672: "test2",
// 					166014673: "test3",
// 					166014674: "test4",
// 				},
// 			},
// 		}, "2006-01-02")
// 		// When
// 		time.Sleep(2 * time.Microsecond)
// 		// Then
// 		assert.Equal(t, false, mockCache.ShouldUpdate("2006-01-02"))
// 	})
// }
