package businesslogic

import (
	"fmt"
	"sync"

	"github.com/MarcBernstein0/match-display/ulits/errorhandling"
)

type tournaments struct {
	tournamentList map[int]tournament
}

type tournament struct {
	tournamentID     int
	tournamentGame   string
	participantsByID map[int]string
}

type match struct {
	player1ID          int
	player1Name        string
	player2ID          int
	player2Name        string
	tournamentID       int
	tournamentGamename string
}

/* calls challenonge api to get all running tournaments
   created recently
   args:
   	none

   returns:
	map[int]string	mapping of tournament IDs and name of the game
	error
*/
func getTournaments(client HTTPClient) (*tournaments, error) {
	// map of tournamentIDs and game names
	tournaments := tournaments{
		tournamentList: make(map[int]tournament),
	}

	// parameters to pass in
	params := map[string]string{
		"state": "in_progress",
	}

	// create request to client
	res := challongeApiCall(client, "tournaments", params)
	if ok, err := errorhandling.HandleError("request failed in getTournaments.", res.err); ok {
		return nil, err
	}

	for _, elem := range res.data {
		if tournamentID, ok := elem["tournament"]["id"].(float64); ok {
			if gameName, ok := elem["tournament"]["game_name"].(string); ok {
				tournaments.tournamentList[int(tournamentID)] = tournament{
					tournamentID:     int(tournamentID),
					tournamentGame:   gameName,
					participantsByID: make(map[int]string),
				}

			} else if elem["tournament"]["game_name"] == nil {
				tournaments.tournamentList[int(tournamentID)] = tournament{
					tournamentID:     int(tournamentID),
					tournamentGame:   "",
					participantsByID: make(map[int]string),
				}
			} else {
				return nil, errorhandling.FormatError(fmt.Sprintf("type for game_name did not match what was expected. Expected='string' got=%T", gameName))
			}
		} else {
			return nil, errorhandling.FormatError(fmt.Sprintf("type for tournament ID did not match what was expected. Expected='float64' got=%T", tournamentID))
		}
	}

	return &tournaments, nil
}

func (t *tournaments) getParticipants(client HTTPClient) error {
	allApiResult := make([]result, 0)

	cResponse := make(chan result)
	var wg sync.WaitGroup
	for k := range t.tournamentList {
		wg.Add(1) // tells the waitgroup that there is no 1 pending operation
		apiPath := fmt.Sprintf("tournaments/%d/participants", k)
		// fmt.Println(v.tournamentGame)
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
			if tournamentID, ok := elem["participant"]["tournament_id"].(float64); ok {
				if name, ok := elem["participant"]["name"].(string); ok {
					if participantID, ok := elem["participant"]["id"].(float64); ok {
						t.tournamentList[int(tournamentID)].participantsByID[int(participantID)] = name
					} else {
						return errorhandling.FormatError(fmt.Sprintf("type for 'participantID' did not match what was expected. Expected='float64' got=%T", participantID))
					}
				} else {
					return errorhandling.FormatError(fmt.Sprintf("type for 'name' did not match what was expected. Expected='string' got=%T", name))
				}
			} else {
				return errorhandling.FormatError(fmt.Sprintf("type for 'tournament_id' did not match what was expected. Expected='float64' got=%T", tournamentID))
			}

		}
	}
	return nil
}

func (t *tournaments) getMatches(client HTTPClient) ([]match, error) {
	// all api results from multiple calls
	allAPIResults := make([]result, 0)

	// slice of matches
	matches := make([]match, 0)

	// parameters to pass in
	params := map[string]string{
		"state": "open",
	}
	// https://api.challonge.com/v1/tournaments/{tournament}/matches.{json|xml}
	cResponse := make(chan result)
	var wg sync.WaitGroup
	for k := range t.tournamentList {
		wg.Add(1)
		apiPath := fmt.Sprintf("tournaments/%d/matches", k)
		// fmt.Println(v.tournamentGame)
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
			var match match
			if tournamentID, ok := elem["match"]["tournament_id"].(float64); ok {
				if player1ID, ok := elem["match"]["player1_id"].(float64); ok {
					match.player1ID = int(player1ID)
					match.player1Name = t.tournamentList[int(tournamentID)].participantsByID[int(player1ID)]
				} else {
					return nil, errorhandling.FormatError(fmt.Sprintf("type for 'player1_id' did not match what was expected. Expected='float64' got=%T", player1ID))
				}
				if player2ID, ok := elem["match"]["player2_id"].(float64); ok {
					match.player2ID = int(player2ID)
					match.player2Name = t.tournamentList[int(tournamentID)].participantsByID[int(player2ID)]
				} else {
					return nil, errorhandling.FormatError(fmt.Sprintf("type for 'player2_id' did not match what was expected. Expected='float64' got=%T", player2ID))
				}
				match.tournamentGamename = t.tournamentList[int(tournamentID)].tournamentGame
				match.tournamentID = int(tournamentID)
				matches = append(matches, match)
			} else {
				return nil, errorhandling.FormatError(fmt.Sprintf("type for 'tournament_id' did not match what was expected. Expected='float64' got=%T", tournamentID))
			}

		}

	}

	return matches, nil
}

func GetTournamentData() (*tournaments, error) {
	fmt.Println("Getting tournament info")
	tournaments, err := getTournaments(client)
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

func (t *tournaments) GetMatches() ([]match, error) {
	matches, err := t.getMatches(client)
	if ok, err := errorhandling.HandleError("failed when calling getMatches", err); ok {
		return nil, err
	}
	return matches, nil
}
