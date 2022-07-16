package models

type (
	Participants []struct {
		Participant Participant `json:"participant"`
	}
	Participant struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	GameParticipants struct {
		GameName     string        `json:"game_name"`
		TournamentID int           `json:"tournament_id"`
		Participant  []Participant `json:"participant"`
	}
)
