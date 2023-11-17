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

func TestMain(m *testing.M) {
	fmt.Println("Mock Server")
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		trimPath := strings.TrimSpace(r.URL.Path)
		switch trimPath {
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	fmt.Println("run tests")
	m.Run()
}

func TestCreateCustomClient(t *testing.T) {
	// Given
	givenCustomClient := customClient{
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
		wantData      map[int]string
		wantErr       error
	}{
		{
			testName:      "response ok one tournament no pagination",
			mockDate:      time.Now().Local().Format("2006-01-02"),
			mockFetchData: New(server.URL, "asdfghjkl", http.DefaultClient, 20),
			wantData: map[int]string{
				1: "test",
				2: "test2",
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
