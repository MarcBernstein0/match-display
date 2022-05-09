package mainlogic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const API_URL = "https://api.challonge.com/v1"

var configuration config

type (
	Tournaments struct {
		Tournament struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"tournament"`
	}
	Participants struct {
		Participant struct {
			Id           int    `json:"id"`
			TournamentId int    `json:"tournament_id"`
			Name         string `json:"name"`
		} `json:"participant"`
	}
	config struct {
		username string
		apiKey   string
	}
	client struct {
		baseURL string
		client  *http.Client
		timeout time.Duration
	}

	Fetch interface {
		FetchTournamentData(ctx context.Context, date string) ([]Tournaments, error)
		FetchParticipantData(ctx context.Context) ([]Participants, error)
	}
)

func init() {
	configuration.username = os.Getenv("USER_NAME")
	configuration.apiKey = os.Getenv("API_KEY")
}

func New(baseURL string, httpClient *http.Client, timeout time.Duration) *client {
	return &client{
		baseURL: baseURL,
		client:  httpClient,
		timeout: timeout,
	}
}

func (c *client) FetchTournamentData(ctx context.Context, date string) ([]Tournaments, error) {
	// get tournament
	var tournaments []Tournaments

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("api_key", configuration.apiKey)
	q.Add("state", "in_progress")
	fmt.Println(date)
	q.Add("created_after", date)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&tournaments)

	return tournaments, nil
}

func (c *client) FetchParticipantData(ctx context.Context) ([]Participants, error) {
	return nil, nil
}
