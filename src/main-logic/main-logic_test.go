package mainlogic

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var server *httptest.Server

const (
	MOCK_API_KEY      = "mock_api_key"
	MOCK_API_USERNAME = "mock_api_username"
)

func testApiKeyAuth(apiKey string) bool {
	if len(apiKey) == 0 {
		return false
	} else if apiKey != MOCK_API_KEY {
		return false
	}
	return true
}

func mockFetchDataEndpoint(w http.ResponseWriter, r *http.Request) {
	apiKey := r.URL.Query().Get("api_key")
	if !testApiKeyAuth(apiKey) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sc := http.StatusOK
	w.WriteHeader(sc)

	date := r.URL.Query().Get("created_after")
	if date == "2022-07-16" {
		w.Write([]byte("[]"))
		return
	}

	jsonFile, err := os.Open("./test-data/testTournamentData.json")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened users.json")

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	fmt.Println(string(byteValue))

	w.Write(byteValue)
}

func TestMain(m *testing.M) {
	fmt.Println("Mocking server")
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// mock calls go here
		switch strings.TrimSpace(r.URL.Path) {
		case "/":
			mockFetchDataEndpoint(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	fmt.Println("mocking customClient")

	fmt.Println("run tests")
	m.Run()
}

func TestCustomClient_FetchData(t *testing.T) {
	tt := []struct {
		name      string
		date      string
		fetchData FetchData
		wantData  []Tournament
		wantErr   error
	}{
		{
			name: "response not ok",
			date: time.Now().Local().Format("2006-01-02"),
			fetchData: func(baseURL, username, apiKey string, client *http.Client) *customClient {
				return New(baseURL, username, apiKey, client)
			}(server.URL, "ashdfhsf", "asdfhdsfh", http.DefaultClient),
			wantData: nil,
			wantErr:  fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(http.StatusUnauthorized)),
		},
		{
			name: "data found",
			date: time.Now().Local().Format("2006-01-02"),
			fetchData: func(baseURL, username, apiKey string, client *http.Client) *customClient {
				return New(baseURL, username, apiKey, client)
			}(server.URL, MOCK_API_USERNAME, MOCK_API_KEY, http.DefaultClient),
			wantData: []Tournament{
				{
					ID:       10878303,
					Name:     "BP GGST 3/4 ",
					GameName: "Guilty Gear -Strive-",
				},
			},
			wantErr: nil,
		},
		{
			name: "no data found but response ok",
			date: "2022-07-16",
			fetchData: func(baseURL, username, apiKey string, client *http.Client) *customClient {
				return New(baseURL, username, apiKey, client)
			}(server.URL, MOCK_API_USERNAME, MOCK_API_KEY, http.DefaultClient),
			wantData: nil,
			wantErr:  fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusNotFound)),
		},
	}

	for _, testCase := range tt {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			gotData, gotErr := testCase.fetchData.FetchTournaments(testCase.date)
			assert.Equal(t, testCase.wantData, gotData)
			if testCase.wantErr != nil {
				assert.EqualError(t, gotErr, testCase.wantErr.Error())
			} else {
				assert.NoError(t, gotErr)
			}

		})
	}
}
