package businesslogic

import (
	"fmt"
	"sync"
)

type tournament struct {
	tournamentID   int
	tournamentGame string
	participants   map[string]int
}

// type matche struct {
// 	player1ID   int
// 	player1Name string
// 	player2ID   int
// 	player2Name string

// }

/* calls challenonge api to get all running tournaments
   created recently
   args:
   	none

   returns:
	map[int]string	mapping of tournament IDs and name of the game
	error
*/
func getTournaments(client HTTPClient) (map[int]tournament, error) {
	// map of tournamentIDs and game names
	tournaments := make(map[int]tournament, 0)

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
				tournaments[int(tournamentID)] = tournament{
					tournamentID:   int(tournamentID),
					tournamentGame: gameName,
					participants:   make(map[string]int),
				}

			} else {
				return nil, fmt.Errorf("type for game_name did not match what was expected. Expected='string' got=%T", gameName)
			}
		} else {
			return nil, fmt.Errorf("type for tournament ID did not match what was expected. Expected='float64' got=%T", tournamentID)
		}
	}

	return tournaments, nil
}

func getParticipants(tournaments map[int]tournament, client HTTPClient) (map[int]tournament, error) {
	allApiResult := make([]result, 0)

	cResponse := make(chan result)
	var wg sync.WaitGroup
	for k, v := range tournaments {
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
			return nil, fmt.Errorf("request failed in getParticipants call.\n%v", resultsApi.err)
		}
		allApiResult = append(allApiResult, resultsApi)
	}

	for _, res := range allApiResult {
		for _, elem := range res.data {
			if tournamentID, ok := elem["participant"]["tournament_id"].(float64); ok {
				if name, ok := elem["participant"]["name"].(string); ok {
					if participantID, ok := elem["participant"]["id"].(float64); ok {

						tournaments[int(tournamentID)].participants[name] = int(participantID)
					} else {
						return nil, fmt.Errorf("type for 'participantID' did not match what was expected. Expected='float64' got=%T", participantID)
					}
				} else {
					return nil, fmt.Errorf("type for 'name' did not match what was expected. Expected='string' got=%T", name)
				}
			} else {
				return nil, fmt.Errorf("type for 'tournament_id' did not match what was expected. Expected='float64' got=%T", tournamentID)
			}

		}
	}
	return tournaments, nil
}

func GetTournamentData() {
	fmt.Println(getTournaments(client))
	// fmt.Println(getParticipants(
	// 	map[int]string{
	// 		3953832:  "Test",
	// 		10469768: "Test2",
	// 	},
	// 	client,
	// ))

}
