package models

type (
	Match struct {
		Id                 string `json:"id"`
		Player1Name        string `json:"player1_name"`
		Player2Name        string `json:"player2_name"`
		Round              int    `json:"round"`
		SuggestedPlayOrder int    `json:"suggested_play_order"`
		Underway           bool   `json:"underway"`
		Station            string `json:"station"`
	}

	TournamentMatches struct {
		GameName     string  `json:"game_name"`
		TournamentId string  `json:"tournament_id"`
		MatchList    []Match `json:"match_list"`
	}
)
