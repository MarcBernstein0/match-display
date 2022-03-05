package businesslogic

import (
	"fmt"
	"sort"
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
	Player1Name        string `json:"player1_name"`
	Player2Name        string `json:"player2_name"`
	Round              int    `json:"round"`
	TournamentGameName string `json:"tournament_game_name"`
}

type Matches struct {
	MatchList []Match `json:"match_list"`
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

	params := map[string]string{
		"state": "in_progress",
	}

	if len(date) != 0 {
		params["created_after"] = date
	}

	// parameters to pass in

	// create request to client
	res := challongeApiCall(client, "tournaments", params)
	if err := errorhandling.HandleError("request failed in getTournaments.", res.err); err != nil {
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
		if err := errorhandling.HandleError("request failed in getParticipants call.", resultsApi.err); err != nil {
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
		if err := errorhandling.HandleError("request failed in getMatches", apiResults.err); err != nil {
			return nil, err
		}
		allAPIResults = append(allAPIResults, apiResults)
	}

	for _, res := range allAPIResults {
		// Test Comment
		for _, elem := range res.data {
			var match Match
			if TournamentID, ok := elem["match"]["tournament_id"].(float64); ok {
				if player1ID, ok := elem["match"]["player1_id"].(float64); ok {
					match.Player1Name = t.TournamentList[int(TournamentID)].ParticipantsByID[int(player1ID)]
				} else {
					return nil, errorhandling.FormatError(fmt.Sprintf("type for 'player1_id' did not match what was expected. Expected='float64' got=%T", player1ID))
				}
				if player2ID, ok := elem["match"]["player2_id"].(float64); ok {
					match.Player2Name = t.TournamentList[int(TournamentID)].ParticipantsByID[int(player2ID)]
				} else {
					return nil, errorhandling.FormatError(fmt.Sprintf("type for 'player2_id' did not match what was expected. Expected='float64' got=%T", player2ID))
				}
				if round, ok := elem["match"]["round"].(float64); ok {
					match.Round = int(round)
				} else {
					return nil, errorhandling.FormatError(fmt.Sprintf("type for 'player2_id' did not match what was expected. Expected='float64' got=%T", round))
				}
				// fmt.Printf("Round=%v type=%T\n", elem["match"]["round"], elem["match"]["round"])

				match.TournamentGameName = t.TournamentList[int(TournamentID)].TournamentGame
				matchList = append(matchList, match)
			} else {
				return nil, errorhandling.FormatError(fmt.Sprintf("type for 'tournament_id' did not match what was expected. Expected='float64' got=%T", TournamentID))
			}

		}

	}

	// order slices
	sort.Slice(matchList, func(i, j int) bool {
		// if the player names are the same, sort by round
		if matchList[i].Player1Name == matchList[j].Player1Name {
			// if the rounds are the same, sory by game
			if matchList[i].Round >= matchList[j].Round {
				return matchList[i].TournamentGameName <= matchList[j].TournamentGameName
			}
			return matchList[i].Round >= matchList[j].Round
		}
		return matchList[i].Player1Name <= matchList[j].Player1Name
		// return false
	})

	return &Matches{
		MatchList: matchList,
	}, nil
}

// func GetTournamentData(date string) (*Tournaments, error) {
func GetTournamentData(date string) (*Tournaments, error) {

	fmt.Println("Getting tournament info")
	tournaments, err := getTournaments(date, client)
	if err = errorhandling.HandleError("failed when calling getTournaments", err); err != nil {
		return nil, err
	}
	// fmt.Println(tournaments)
	err = tournaments.getParticipants(client)
	if err = errorhandling.HandleError("failed when calling getParticipants", err); err != nil {
		return nil, err
	}
	// fmt.Println(tournaments)

	return tournaments, nil
}

func (t *Tournaments) GetMatches() (*Matches, error) {
	matches, err := t.getMatches(client)
	if err = errorhandling.HandleError("failed when calling getMatches", err); err != nil {
		return nil, err
	}
	return matches, nil
}
