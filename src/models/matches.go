package models

type (
	Matches []struct {
		Match Match `json:"match"`
	}
	Match struct {
		ID          int    `json:"id"`
		Player1ID   int    `json:"player1_id"`
		Player1Name string `json:"player1_name"`
		Player2ID   int    `json:"player2_id"`
		Player2Name string `json:"player2_name"`
		Round       int    `json:"round"`
	}
	TournamentMatches struct {
		GameName     string  `json:"game_name"`
		TournamentID int     `json:"tournament_id"`
		MatchList    []Match `json:"match_list"`
	}
)
