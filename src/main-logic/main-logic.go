package mainlogic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrResponseNotOK error = errors.New("response not ok")
	ErrServerProblem error = errors.New("server error")
)

type (
	Tournaments []struct {
		Tournament Tournament `json:"tournament"`
	}
	Tournament struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		GameName string `json:"game_name"`
	}

	FetchData interface {
		FetchTournaments(date string) ([]Tournament, error)
	}

	customClient struct {
		baseURL string
		client  *http.Client
		config  struct {
			username string
			apiKey   string
		}
	}
)

func New(baseURL, username, apiKey string, client *http.Client) *customClient {
	return &customClient{
		baseURL: baseURL,
		client:  client,
		config: struct {
			username string
			apiKey   string
		}{
			username: username,
			apiKey:   apiKey,
		},
	}
}

func (c *customClient) FetchTournaments(date string) ([]Tournament, error) {
	// ctx, cancel := context.WithTimeout(ctx, c.timeout)
	// defer cancel()

	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	// if err != nil {
	// 	return nil, err
	// }
	req, err := http.NewRequest(http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("api_key", c.config.apiKey)
	q.Add("state", "in_progress")
	fmt.Println(date)
	q.Add("created_after", date)
	req.URL.RawQuery = q.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode))
	}

	var tournaments Tournaments
	err = json.NewDecoder(res.Body).Decode(&tournaments)
	if len(tournaments) == 0 {
		return nil, fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusNotFound))
	}
	if err != nil {
		return nil, fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusInternalServerError))
	}
	fmt.Printf("%+v, %v\n", tournaments, len(tournaments))
	var tournamentList []Tournament
	for _, t := range tournaments {
		tournamentList = append(tournamentList, t.Tournament)
	}
	return tournamentList, err
}
