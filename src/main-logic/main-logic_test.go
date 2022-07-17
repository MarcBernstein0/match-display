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

	"github.com/MarcBernstein0/match-display/src/models"
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

func readJsonFile(filename string) ([]byte, error) {
	jsonFile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("Successfully Opened users.json")

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	return byteValue, err

}

func mockFetchTournamentEndpoint(w http.ResponseWriter, r *http.Request) {
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

	byteValue, err := readJsonFile("./test-data/testTournamentData.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// fmt.Println(string(byteValue))

	w.Write(byteValue)
}

func mockFetchParticipantEndpoint(w http.ResponseWriter, r *http.Request) {
	apiKey := r.URL.Query().Get("api_key")
	if !testApiKeyAuth(apiKey) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sc := http.StatusOK
	w.WriteHeader(sc)

	byteValue, err := readJsonFile("./test-data/testParticipantData.json")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// fmt.Println(string(byteValue))

	w.Write(byteValue)
}

func mockFetchParticipantEndpoint2(w http.ResponseWriter, r *http.Request) {
	apiKey := r.URL.Query().Get("api_key")
	if !testApiKeyAuth(apiKey) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sc := http.StatusOK
	w.WriteHeader(sc)

	jsonFile, err := os.Open("./test-data/testParticipantData.json")
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("Successfully Opened users.json")

	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	// fmt.Println(string(byteValue))

	w.Write(byteValue)
}

func TestMain(m *testing.M) {
	fmt.Println("Mocking server")
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// mock calls go here
		switch strings.TrimSpace(r.URL.Path) {
		case "/tournaments.json":
			mockFetchTournamentEndpoint(w, r)
		case "/tournaments/10879090/participants.json":
			mockFetchParticipantEndpoint(w, r)
		case "/tournaments/10879091/participants.json":
			mockFetchParticipantEndpoint2(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	fmt.Println("mocking customClient")

	fmt.Println("run tests")
	m.Run()
}

func TestCustomClient_FetchTournaments(t *testing.T) {
	tt := []struct {
		name      string
		date      string
		fetchData FetchData
		wantData  []models.Tournament
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
			wantData: []models.Tournament{
				{
					ID:       10879090,
					Name:     "test",
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

func TestCustomClient_FetchParticipants(t *testing.T) {
	tt := []struct {
		name      string
		fetchData FetchData
		inputData []models.Tournament
		wantData  []models.GameParticipants
		wantErr   error
	}{
		{
			name: "response not ok",
			fetchData: func(baseURL, username, apiKey string, client *http.Client) *customClient {
				return New(baseURL, username, apiKey, client)
			}(server.URL, "ashdfhsf", "asdfhdsfh", http.DefaultClient),
			inputData: []models.Tournament{
				{
					ID:       10879090,
					Name:     "test",
					GameName: "Guilty Gear -Strive-",
				},
			},
			wantData: nil,
			wantErr:  fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(http.StatusUnauthorized)),
		},
		{
			name: "data found",
			fetchData: func(baseURL, username, apiKey string, client *http.Client) *customClient {
				return New(baseURL, username, apiKey, client)
			}(server.URL, MOCK_API_USERNAME, MOCK_API_KEY, http.DefaultClient),
			inputData: []models.Tournament{
				{
					ID:       10879090,
					Name:     "test",
					GameName: "Guilty Gear -Strive-",
				},
			},
			wantData: []models.GameParticipants{
				{
					GameName:     "Guilty Gear -Strive-",
					TournamentID: 10879090,
					Participant: []models.Participant{
						{
							ID:   166014671,
							Name: "test",
						},
						{
							ID:   166014672,
							Name: "test2",
						},
						{
							ID:   166014673,
							Name: "test3",
						},
						{
							ID:   166014674,
							Name: "test4",
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "multiple tournaments",
			fetchData: func(baseURL, username, apiKey string, client *http.Client) *customClient {
				return New(baseURL, username, apiKey, client)
			}(server.URL, MOCK_API_USERNAME, MOCK_API_KEY, http.DefaultClient),
			inputData: []models.Tournament{
				{
					ID:       10879090,
					Name:     "test",
					GameName: "Guilty Gear -Strive-",
				},
				{
					ID:       10879091,
					Name:     "test2",
					GameName: "DNF Duel",
				},
			},
			wantData: []models.GameParticipants{
				{
					GameName:     "Guilty Gear -Strive-",
					TournamentID: 10879090,
					Participant: []models.Participant{
						{
							ID:   166014671,
							Name: "test",
						},
						{
							ID:   166014672,
							Name: "test2",
						},
						{
							ID:   166014673,
							Name: "test3",
						},
						{
							ID:   166014674,
							Name: "test4",
						},
					},
				},
				{
					GameName:     "DNF Duel",
					TournamentID: 10879091,
					Participant: []models.Participant{
						{
							ID:   166014671,
							Name: "test",
						},
						{
							ID:   166014672,
							Name: "test2",
						},
						{
							ID:   166014673,
							Name: "test3",
						},
						{
							ID:   166014674,
							Name: "test4",
						},
					},
				},
			},
			wantErr: nil,
		},
	}

	for _, testCase := range tt {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			gotData, gotErr := testCase.fetchData.FetchParticipants(testCase.inputData)
			assert.ElementsMatch(t, testCase.wantData, gotData)
			if testCase.wantErr != nil {
				assert.EqualError(t, gotErr, testCase.wantErr.Error())
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
