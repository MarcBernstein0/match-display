package models

type (
	Participants []struct {
		Participant Participant `json:"participant"`
	}
	Participant struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	TournamentParticipants struct {
		GameName     string         `json:"game_name"`
		TournamentID int            `json:"tournament_id"`
		Participant  map[int]string `json:"participant"`
	}
)
