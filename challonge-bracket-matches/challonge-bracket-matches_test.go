package challongebracketmatches

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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
		baseURL: "testEndpoint",
		client:  http.DefaultClient,
		config: struct{ apiKey string }{
			apiKey: "1234567890",
		},
		contextTimeout: 20,
	}
	// When
	res := New("testEndpoint", "1234567890", http.DefaultClient, 20)
	// Then
	require.Equal(t, givenCustomClient, res)
}
