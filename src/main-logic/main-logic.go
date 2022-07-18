package mainlogic

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/MarcBernstein0/match-display/src/models"
)

var (
	ErrResponseNotOK error = errors.New("response not ok")
	ErrServerProblem error = errors.New("server error")
)

type (
	participantResult struct {
		tournamentParticipant *models.TournamentParticipants
		error                 error
	}
	FetchData interface {
		// FetchTournaments fetch all tournaments created after a specific date
		// GET https://api.challonge.com/v1/tournaments.{json|xml}
		FetchTournaments(date string) ([]models.Tournament, error)
		// FetchParticipants of a given tournament
		// GET https://api.challonge.com/v1/tournaments/{tournament}/participants.{json|xml}
		FetchParticipants(tournaments []models.Tournament) ([]models.TournamentParticipants, error)
		// FetchMatches of a given tournament
		// GET https://api.challonge.com/v1/tournaments/{tournament}/matches.{json|xml}
		FetchMatches(tournamentParticipants []models.TournamentParticipants) ([]models.TournamentMatches, error)
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
	// fmt.Println(date)
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

func (c *customClient) fetchAllParticipants(tournament models.Tournament, participantResultChan chan<- participantResult, wg *sync.WaitGroup) {
	defer wg.Done()
	tournamentID := tournament.ID
	gameName := tournament.GameName
	url := fmt.Sprintf("%s/tournaments/%v/participants.json", c.baseURL, tournamentID)
	// fmt.Println(url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		// return nil, err
		participantResultChan <- participantResult{
			tournamentParticipant: nil,
			error:                 err,
		}
		return
	}
	q := req.URL.Query()
	q.Add("api_key", c.config.apiKey)
	req.URL.RawQuery = q.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		// return nil, err
		participantResultChan <- participantResult{
			tournamentParticipant: nil,
			error:                 err,
		}
		return
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		// return nil, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode))
		participantResultChan <- participantResult{
			tournamentParticipant: nil,
			error:                 fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode)),
		}
		return
	}

	var participants models.Participants
	err = json.NewDecoder(res.Body).Decode(&participants)
	if len(participants) == 0 {
		// return nil, fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusNotFound))
		participantResultChan <- participantResult{
			tournamentParticipant: nil,
			error:                 fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusNotFound)),
		}
		return
	}
	if err != nil {
		// return nil, fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusInternalServerError))
		participantResultChan <- participantResult{
			tournamentParticipant: nil,
			error:                 fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusInternalServerError)),
		}
		return
	}
	fmt.Printf("%+v, %v\n", participants, len(participants))

	tournamentParticipant := models.TournamentParticipants{
		GameName:     gameName,
		TournamentID: tournamentID,
		Participant:  map[int]string{},
	}
	for _, p := range participants {
		tournamentParticipant.Participant[p.Participant.ID] = p.Participant.Name
	}

	participantResultChan <- participantResult{
		tournamentParticipant: &tournamentParticipant,
		error:                 nil,
	}

}

func (c *customClient) FetchParticipants(tournaments []models.Tournament) ([]models.TournamentParticipants, error) {
	var tournamentParticipants []models.TournamentParticipants

	cResponse := make(chan participantResult)
	var wg sync.WaitGroup
	for _, tournament := range tournaments {
		wg.Add(1) // add one to the waitgroup
		go c.fetchAllParticipants(tournament, cResponse, &wg)
	}

	go func() {
		wg.Wait()
		close(cResponse)
	}()

	for tournamentParticipantResult := range cResponse {
		if tournamentParticipantResult.error != nil {
			return nil, tournamentParticipantResult.error
		}
		tournamentParticipants = append(tournamentParticipants, *tournamentParticipantResult.tournamentParticipant)

	}

	fmt.Printf("Final game participants: %+v", tournamentParticipants)
	return tournamentParticipants, nil
}

func (c *customClient) FetchMatches(tournamentParticipants []models.TournamentParticipants) ([]models.TournamentMatches, error) {

	// TODO: use channels to be able to do this for multiple tournaments all at once
	var allMatches []models.TournamentMatches
	tournamentID := tournamentParticipants[0].TournamentID
	gameName := tournamentParticipants[0].GameName

	participants := tournamentParticipants[0].Participant
	url := fmt.Sprintf("%s/tournaments/%v/matches.json", c.baseURL, tournamentID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("api_key", c.config.apiKey)
	q.Add("state", "open")
	// fmt.Println(date)
	req.URL.RawQuery = q.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode))
	}

	var matches models.Matches
	err = json.NewDecoder(res.Body).Decode(&matches)
	if err != nil {
		return nil, fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusInternalServerError))
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("%w. %s", ErrServerProblem, http.StatusText(http.StatusNotFound))
	}
	fmt.Printf("%+v\n", matches)
	tournamentMatches := models.TournamentMatches{
		GameName:     gameName,
		TournamentID: tournamentID,
		MatchList:    make([]models.Match, 0),
	}
	for _, m := range matches {
		m.Match.Player1Name = participants[m.Match.Player1ID]
		m.Match.Player2Name = participants[m.Match.Player2ID]
		tournamentMatches.MatchList = append(tournamentMatches.MatchList, m.Match)
	}
	fmt.Printf("%+v\n", tournamentMatches)
	allMatches = append(allMatches, tournamentMatches)

	return allMatches, nil
}
