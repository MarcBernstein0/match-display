package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	mainlogic "github.com/MarcBernstein0/match-display/src/main-logic"
	"github.com/MarcBernstein0/match-display/src/models"
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

func mockFetchMockDataEndpoint(w http.ResponseWriter, r *http.Request) {
	sc := http.StatusOK
	w.WriteHeader(sc)
	w.Write([]byte("{'test':'test'"))
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
	if date == "2022-07-20" {
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

// func mockFetchMatchDataEndpoint(w http.ResponseWriter, r *http.Request) {
// 	apiKey := r.URL.Query().Get("api_key")
// 	if !testApiKeyAuth(apiKey) {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	sc := http.StatusOK
// 	w.WriteHeader(sc)
// 	date := r.URL.Query().Get("created_after")
// 	if date == "2022-07-20" {
// 		w.Write([]byte("[]"))
// 		return
// 	}
// 	byteValue, err := readJsonFile("./test-data/testTournamentData.json")
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		w.Write([]byte(err.Error()))
// 		return
// 	}

// 	w.Write(byteValue)
// }

func TestMain(m *testing.M) {

	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch strings.TrimSpace(r.URL.Path) {
		case "/":
			mockFetchMockDataEndpoint(w, r)
		case "/tournaments.json":
			mockFetchTournamentEndpoint(w, r)
		// case "/":
		// 	mockFetchMatchDataEndpoint(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	mockFetch = mainlogic.New(server.URL, MOCK_API_USERNAME, MOCK_API_KEY, http.DefaultClient)
	m.Run()
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
	tt := []struct {
		name       string
		date       string
		statusCode int
		wantData   []models.TournamentMatches
		expectErr  bool
		wantErr    models.ErrorResponse
	}{
		{
			name:       "response not ok",
			date:       "",
			statusCode: http.StatusBadRequest,
			wantData:   nil,
			expectErr:  true,
			wantErr: models.ErrorResponse{
				Message:      "did not fill out required 'date' query field",
				ErrorMessage: "Key: 'Date.Date' Error:Field validation for 'Date' failed on the 'required' tag",
			},
		},
		{
			name:       "response not ok",
			date:       "2022-07-20",
			statusCode: http.StatusInternalServerError,
			wantData:   nil,
			expectErr:  true,
			wantErr: models.ErrorResponse{
				Message:      "failed to get tournament data",
				ErrorMessage: fmt.Errorf("%w. %s", mainlogic.ErrNoData, http.StatusText(http.StatusNotFound)).Error(),
			},
		},
		// {
		// 	name:       "single tournament",
		// 	date:       "2022-07-19",
		// 	statusCode: http.StatusOK,
		// 	wantData: []models.TournamentMatches{
		// 		{
		// 			GameName:     "Guilty Gear -Strive-",
		// 			TournamentID: 10879090,
		// 			MatchList: []models.Match{
		// 				{
		// 					ID:          267800918,
		// 					Player1ID:   166014671,
		// 					Player1Name: "test",
		// 					Player2ID:   166014674,
		// 					Player2Name: "test4",
		// 					Round:       1,
		// 				},
		// 				{
		// 					ID:          267800919,
		// 					Player1ID:   166014672,
		// 					Player1Name: "test2",
		// 					Player2ID:   166014673,
		// 					Player2Name: "test3",
		// 					Round:       1,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantErr: nil,
		// },
		// {
		// 	name:       "multiple tournaments",
		// 	date:       time.Now().Local().Format("2006-01-02"),
		// 	statusCode: http.StatusOK,
		// 	wantData: []models.TournamentMatches{
		// 		{
		// 			GameName:     "Guilty Gear -Strive-",
		// 			TournamentID: 10879090,
		// 			MatchList: []models.Match{
		// 				{
		// 					ID:          267800918,
		// 					Player1ID:   166014671,
		// 					Player1Name: "test",
		// 					Player2ID:   166014674,
		// 					Player2Name: "test4",
		// 					Round:       1,
		// 				},
		// 				{
		// 					ID:          267800919,
		// 					Player1ID:   166014672,
		// 					Player1Name: "test2",
		// 					Player2ID:   166014673,
		// 					Player2Name: "test3",
		// 					Round:       1,
		// 				},
		// 			},
		// 		},
		// 		{
		// 			GameName:     "DNF Duel",
		// 			TournamentID: 10879091,
		// 			MatchList: []models.Match{
		// 				{
		// 					ID:          267800918,
		// 					Player1ID:   166014671,
		// 					Player1Name: "test",
		// 					Player2ID:   166014674,
		// 					Player2Name: "test4",
		// 					Round:       1,
		// 				},
		// 				{
		// 					ID:          267800919,
		// 					Player1ID:   166014672,
		// 					Player1Name: "test2",
		// 					Player2ID:   166014673,
		// 					Player2Name: "test3",
		// 					Round:       1,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantErr: nil,
		// },
	}

	router := RouteSetup(mockFetch)

	for _, testCase := range tt {
		t.Run(testCase.name, func(t *testing.T) {
			// t.Parallel()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/matches", nil)
			q := req.URL.Query()
			q.Add("date", testCase.date)
			req.URL.RawQuery = q.Encode()

			router.ServeHTTP(w, req)
			assert.Equal(t, testCase.statusCode, w.Code)
			if testCase.expectErr {
				var gotErr models.ErrorResponse
				err := json.NewDecoder(w.Body).Decode(&gotErr)
				if err != nil {
					t.Fatalf("failed to decode error response %v", err)
				}
				assert.Equal(t, testCase.wantErr, gotErr)
			} else {
				var gotData []models.TournamentMatches
				err := json.NewDecoder(w.Body).Decode(&gotData)
				if err != nil {
					t.Fatalf("failed to decode error response %v", err)
				}
				assert.ElementsMatch(t, testCase.wantData, gotData)
			}
		})
	}

}
