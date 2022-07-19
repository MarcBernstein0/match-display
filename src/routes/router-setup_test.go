package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mainlogic "github.com/MarcBernstein0/match-display/src/main-logic"
	"github.com/stretchr/testify/assert"
)

var (
	server    *httptest.Server
	mockFetch mainlogic.FetchData
)

const (
	MOCK_API_KEY      = "mock_api_key"
	MOCK_API_USERNAME = "mock_api_username"
)

func mockFetchMockDataEndpoint(w http.ResponseWriter, r *http.Request) {
	sc := http.StatusOK
	w.WriteHeader(sc)
	w.Write([]byte("{'test':'test'"))
}

func TestMain(m *testing.M) {
	mockFetch = mainlogic.New(server.URL, MOCK_API_USERNAME, MOCK_API_KEY, http.DefaultClient)

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch strings.TrimSpace(r.URL.Path) {
		case "/":
			mockFetchMockDataEndpoint(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))
}

func TestHealthCheckRoute(t *testing.T) {
	router := RouteSetup(mockFetch)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"status": "UP"}`, w.Body.String())
}

func TestGetMatchesRoute(t *testing.T) {
	// tt := []struct {
	// 	name     string
	// 	date     string
	// 	wantData []models.TournamentMatches
	// 	wantErr  error
	// }{
	// 	{
	// 		name:     "response not ok",
	// 		date:     time.Now().Local().Format("2006-01-02"),
	// 		wantData: nil,
	// 		wantErr:  fmt.Errorf("%w. %s", mainlogic.ErrResponseNotOK, http.StatusText(http.StatusUnauthorized)),
	// 	},
	// }
}
