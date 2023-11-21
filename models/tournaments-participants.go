package models

type TournamentParticipants struct {
	GameName     string         `json:"game_name"`
	TournamentID int            `json:"tournament_id"`
	Participant  map[int]string `json:"participant"`
}
