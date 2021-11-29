package businesslogic

import (
	"fmt"
	"sync"

	"github.com/MarcBernstein0/match-display/ulits/errorhandling"
)

type Tournaments struct {
	TournamentList map[int]tournament `json:"tournament_list"`
}

type tournament struct {
	TournamentID     int            `json:"tournament_id"`
	TournamentGame   string         `json:"tournamnet_game"`
	ParticipantsByID map[int]string `json:"participants_by_id"`
}

type Match struct {
	Player1ID          int    `json:"player1_id"`
	Player1Name        string `json:"player1_name"`
	Player2ID          int    `json:"player2_id"`
	Player2Name        string `json:"player2_name"`
	TournamentID       int    `json:"tournament_id"`
	TournamentGameName string `json:"tournament_game_name"`
}

type Matches struct {
	MatchList []Match
}

/* calls challenonge api to get all running tournaments
   created recently
   args:
   	none

   returns:
	map[int]string	mapping of tournament IDs and name of the game
	error
*/
func getTournaments(date string, client HTTPClient) (*Tournaments, error) {
	// map of TournamentIDs and game names
	tournaments := Tournaments{
		TournamentList: make(map[int]tournament),
	}

	// parameters to pass in
	params := map[string]string{
		"state":         "in_progress",
		"created_after": date,
	}

	// create request to client
	res := challongeApiCall(client, "tournaments", params)
	if ok, err := errorhandling.HandleError("request failed in getTournaments.", res.err); ok {
		return nil, err
	}

	for _, elem := range res.data {
		if TournamentID, ok := elem["tournament"]["id"].(float64); ok {
			if gameName, ok := elem["tournament"]["game_name"].(string); ok {
				tournaments.TournamentList[int(TournamentID)] = tournament{
					TournamentID:     int(TournamentID),
					TournamentGame:   gameName,
					ParticipantsByID: make(map[int]string),
				}

			} else if elem["tournament"]["game_name"] == nil {
				tournaments.TournamentList[int(TournamentID)] = tournament{
					TournamentID:     int(TournamentID),
					TournamentGame:   "",
					ParticipantsByID: make(map[int]string),
				}
			} else {
				return nil, errorhandling.FormatError(fmt.Sprintf("type for game_name did not match what was expected. Expected='string' got=%T", gameName))
			}
		} else {
			return nil, errorhandling.FormatError(fmt.Sprintf("type for tournament ID did not match what was expected. Expected='float64' got=%T", TournamentID))
		}
	}

	return &tournaments, nil
}

func (t *Tournaments) getParticipants(client HTTPClient) error {
	allApiResult := make([]result, 0)

	cResponse := make(chan result)
	var wg sync.WaitGroup
	for k := range t.TournamentList {
		wg.Add(1) // tells the waitgroup that there is no 1 pending operation
		apiPath := fmt.Sprintf("tournaments/%d/participants", k)
		// fmt.Println(v.TournamentGame)
		go challongeApiMultiCall(client, apiPath, nil, cResponse, &wg)
	}

	go func() {
		wg.Wait()
		close(cResponse)
	}()

	for resultsApi := range cResponse {
		if ok, err := errorhandling.HandleError("request failed in getParticipants call.", resultsApi.err); ok {
			return err
		}
		allApiResult = append(allApiResult, resultsApi)
	}

	for _, res := range allApiResult {
		for _, elem := range res.data {
			if TournamentID, ok := elem["participant"]["tournament_id"].(float64); ok {
				if name, ok := elem["participant"]["name"].(string); ok {
					if participantID, ok := elem["participant"]["id"].(float64); ok {
						t.TournamentList[int(TournamentID)].ParticipantsByID[int(participantID)] = name
					} else {
						return errorhandling.FormatError(fmt.Sprintf("type for 'participantID' did not match what was expected. Expected='float64' got=%T", participantID))
					}
				} else {
					return errorhandling.FormatError(fmt.Sprintf("type for 'name' did not match what was expected. Expected='string' got=%T", name))
				}
			} else {
				return errorhandling.FormatError(fmt.Sprintf("type for 'tournament_id' did not match what was expected. Expected='float64' got=%T", TournamentID))
			}

		}
	}
	return nil
}

func (t *Tournaments) getMatches(client HTTPClient) (*Matches, error) {
	// all api results from multiple calls
	allAPIResults := make([]result, 0)

	// slice of matches
	matchList := make([]Match, 0)

	// parameters to pass in
	params := map[string]string{
		"state": "open",
	}
	// https://api.challonge.com/v1/tournaments/{tournament}/matches.{json|xml}
	cResponse := make(chan result)
	var wg sync.WaitGroup
	for k := range t.TournamentList {
		wg.Add(1)
		apiPath := fmt.Sprintf("tournaments/%d/matches", k)
		// fmt.Println(v.TournamentGame)
		go challongeApiMultiCall(client, apiPath, params, cResponse, &wg)
	}

	go func() {
		wg.Wait()
		close(cResponse)
	}()

	for apiResults := range cResponse {
		if ok, err := errorhandling.HandleError("request failed in getMatches", apiResults.err); ok {
			return nil, err
		}
		allAPIResults = append(allAPIResults, apiResults)
	}

	for _, res := range allAPIResults {
		for _, elem := range res.data {
			var match Match
			if TournamentID, ok := elem["match"]["tournament_id"].(float64); ok {
				if player1ID, ok := elem["match"]["player1_id"].(float64); ok {
					match.Player1ID = int(player1ID)
					match.Player1Name = t.TournamentList[int(TournamentID)].ParticipantsByID[int(player1ID)]
				} else {
					return nil, errorhandling.FormatError(fmt.Sprintf("type for 'player1_id' did not match what was expected. Expected='float64' got=%T", player1ID))
				}
				if player2ID, ok := elem["match"]["player2_id"].(float64); ok {
					match.Player2ID = int(player2ID)
					match.Player2Name = t.TournamentList[int(TournamentID)].ParticipantsByID[int(player2ID)]
				} else {
					return nil, errorhandling.FormatError(fmt.Sprintf("type for 'player2_id' did not match what was expected. Expected='float64' got=%T", player2ID))
				}
				match.TournamentGameName = t.TournamentList[int(TournamentID)].TournamentGame
				match.TournamentID = int(TournamentID)
				matchList = append(matchList, match)
			} else {
				return nil, errorhandling.FormatError(fmt.Sprintf("type for 'tournament_id' did not match what was expected. Expected='float64' got=%T", TournamentID))
			}

		}

	}

	return &Matches{
		MatchList: matchList,
	}, nil
}

func GetTournamentData(date string) (*Tournaments, error) {
	fmt.Println("Getting tournament info")
	tournaments, err := getTournaments(date, client)
	if ok, err := errorhandling.HandleError("failed when calling getTournaments", err); ok {
		return nil, err
	}
	// fmt.Println(tournaments)
	err = tournaments.getParticipants(client)
	if ok, err := errorhandling.HandleError("failed when calling getParticipants", err); ok {
		return nil, err
	}
	// fmt.Println(tournaments)

	return tournaments, nil
}

func (t *Tournaments) GetMatches() (*Matches, error) {
	matches, err := t.getMatches(client)
	if ok, err := errorhandling.HandleError("failed when calling getMatches", err); ok {
		return nil, err
	}
	return matches, nil
}
