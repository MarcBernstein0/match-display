package models

type (
	Tournaments []struct {
		Tournament Tournament `json:"tournament"`
	}
	Tournament struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		GameName string `json:"game_name"`
	}
)
