package mainlogic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MarcBernstein0/match-display/src/models"
)

var (
	ErrResponseNotOK error = errors.New("response not ok")
	ErrServerProblem error = errors.New("server error")
)

type (
	FetchData interface {
		// FetchTournaments fetch all tournaments created after a specific date
		// GET https://api.challonge.com/v1/tournaments.{json|xml}
		FetchTournaments(date string) ([]models.Tournament, error)
		// FetchParicipants of a given tournament
		// GET https://api.challonge.com/v1/tournaments/{tournament}/participants.{json|xml}
		FetchParicipants(tournaments []models.Tournament) ([]models.GameParticipants, error)
		// FetchMatches of a given tournament
		// GET https://api.challonge.com/v1/tournaments/{tournament}/matches.{json|xml}

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

func (c *customClient) FetchTournaments(date string) ([]models.Tournament, error) {
	// ctx, cancel := context.WithTimeout(ctx, c.timeout)
	// defer cancel()

	// req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	// if err != nil {
	// 	return nil, err
	// }
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/tournaments.json", nil)
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

	var tournaments models.Tournaments
	err = json.NewDecoder(res.Body).Decode(&tournaments)
	if len(tournaments) == 0 {
		return nil, fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusNotFound))
	}
	if err != nil {
		return nil, fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusInternalServerError))
	}
	fmt.Printf("%+v, %v\n", tournaments, len(tournaments))
	var tournamentList []models.Tournament
	for _, t := range tournaments {
		tournamentList = append(tournamentList, t.Tournament)
	}
	return tournamentList, err
}

func (c *customClient) FetchParicipants(tournaments []models.Tournament) ([]models.GameParticipants, error) {
	// final list of all participants for all tournaments
	var gameParticipants []models.GameParticipants

	tournametID := tournaments[0].ID
	gameName := tournaments[0].GameName
	url := fmt.Sprintf("%s/tournaments/%v/participants.json", c.baseURL, tournametID)
	// fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("api_key", c.config.apiKey)
	req.URL.RawQuery = q.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode))
	}

	var participants models.Participants
	err = json.NewDecoder(res.Body).Decode(&participants)
	if len(participants) == 0 {
		return nil, fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusNotFound))
	}
	if err != nil {
		return nil, fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusInternalServerError))
	}
	fmt.Printf("%+v, %v\n", participants, len(participants))

	gameParticipant := models.GameParticipants{
		GameName:     gameName,
		TournamentID: tournametID,
		Participant:  make([]models.Participant, 0),
	}
	for _, p := range participants {
		gameParticipant.Participant = append(gameParticipant.Participant, p.Participant)
	}

	gameParticipants = append(gameParticipants, gameParticipant)
	fmt.Println(gameParticipants)
	return gameParticipants, nil
}
