package businesslogic

import (
	"fmt"
	"sync"
)

type tournaments struct {
	tournamentList map[int]tournament
}

type tournament struct {
	tournamentID       int
	tournamentGame     string
	participantsByName map[string]int
	participantsByID   map[int]string
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
	if res.err != nil {
		return nil, fmt.Errorf("request failed in getTournaments.\n%v", res.err)
	}

	for _, elem := range res.data {
		if tournamentID, ok := elem["tournament"]["id"].(float64); ok {

			if gameName, ok := elem["tournament"]["game_name"].(string); ok {
				tournaments.tournamentList[int(tournamentID)] = tournament{
					tournamentID:       int(tournamentID),
					tournamentGame:     gameName,
					participantsByName: make(map[string]int),
					participantsByID:   make(map[int]string),
				}

			} else {
				return nil, fmt.Errorf("type for game_name did not match what was expected. Expected='string' got=%T", gameName)
			}
		} else {
			return nil, fmt.Errorf("type for tournament ID did not match what was expected. Expected='float64' got=%T", tournamentID)
		}
	}

	return &tournaments, nil
}

func (t *tournaments) getParticipants(client HTTPClient) error {
	allApiResult := make([]result, 0)

	cResponse := make(chan result)
	var wg sync.WaitGroup
	for k, v := range t.tournamentList {
		wg.Add(1) // tells the waitgroup that there is no 1 pending operation
		apiPath := fmt.Sprintf("tournaments/%d/participants", k)
		fmt.Println(v.tournamentGame)
		go challongeApiMultiCall(client, apiPath, nil, cResponse, &wg)
	}

	go func() {
		wg.Wait()
		close(cResponse)
	}()

	for resultsApi := range cResponse {
		if resultsApi.err != nil {
			return fmt.Errorf("request failed in getParticipants call.\n%v", resultsApi.err)
		}
		allApiResult = append(allApiResult, resultsApi)
	}

	for _, res := range allApiResult {
		for _, elem := range res.data {
			if tournamentID, ok := elem["participant"]["tournament_id"].(float64); ok {
				if name, ok := elem["participant"]["name"].(string); ok {
					if participantID, ok := elem["participant"]["id"].(float64); ok {
						t.tournamentList[int(tournamentID)].participantsByName[name] = int(participantID)
						t.tournamentList[int(tournamentID)].participantsByID[int(participantID)] = name

					} else {
						return fmt.Errorf("type for 'participantID' did not match what was expected. Expected='float64' got=%T", participantID)
					}
				} else {
					return fmt.Errorf("type for 'name' did not match what was expected. Expected='string' got=%T", name)
				}
			} else {
				return fmt.Errorf("type for 'tournament_id' did not match what was expected. Expected='float64' got=%T", tournamentID)
			}

		}
	}
	return nil
}

func (t *tournaments) getMatches(client HTTPClient) ([]match, error) {

	// slice of matches
	matches := make([]match, 0)

	// parameters to pass in
	params := map[string]string{
		"state": "open",
	}
	// https://api.challonge.com/v1/tournaments/{tournament}/matches.{json|xml}
	apiPath := fmt.Sprintf("tournaments/%d/matches", 10469768)

	res := challongeApiCall(client, apiPath, params)
	if res.err != nil {
		return nil, fmt.Errorf("request failed in getMatches\n%v", res.err)
	}

	for _, elem := range res.data {
		var match match
		if player1ID, ok := elem["match"]["player1_id"].(float64); ok {
			match.player1ID = int(player1ID)
			match.player1Name = t.tournamentList[10469768].participantsByID[int(player1ID)]
		} else {
			return nil, fmt.Errorf("type for 'player1_id' did not match what was expected. Expected='float64' got=%T", player1ID)
		}
		if player2ID, ok := elem["match"]["player2_id"].(float64); ok {
			match.player2ID = int(player2ID)
			match.player2Name = t.tournamentList[10469768].participantsByID[int(player2ID)]
		} else {
			return nil, fmt.Errorf("type for 'player2_id' did not match what was expected. Expected='float64' got=%T", player2ID)
		}
		match.tournamentGamename = t.tournamentList[10469768].tournamentGame
		match.tournamentID = 10469768

		matches = append(matches, match)
	}
	player1 := int(res.data[0]["match"]["player1_id"].(float64))
	fmt.Println(player1)
	fmt.Println(t.tournamentList[10469768].participantsByID[player1])
	return matches, nil
}

func GetTournamentData() {
	tournaments, err := getTournaments(client)
	if err != nil {
		panic(err)
	}
	fmt.Println(tournaments)
	err = tournaments.getParticipants(client)
	if err != nil {
		panic(err)
	}
	fmt.Println(tournaments)
}
