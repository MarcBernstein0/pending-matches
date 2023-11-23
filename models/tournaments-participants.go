package models

type TournamentParticipants struct {
	GameName     string            `json:"game_name"`
	TournamentID string            `json:"tournament_id"`
	Participant  map[string]string `json:"participant"`
}
