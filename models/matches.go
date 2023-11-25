package models

import "time"

type (
	Matches struct {
		Data     []ChallongeMatch `json:"data"`
		Included []Included       `json:"included"`
	}

	ChallongeMatch struct {
		Id           string          `json:"id"`
		Attributes   MatchAttributes `json:"attributes"`
		Relationship Relationships   `json:"relationships"`
	}

	Relationships struct {
		Station Station `json:"station"`
	}

	Station struct {
		Data StationData `json:"data"`
	}

	StationData struct {
		Id string `json:"id"`
	}

	MatchAttributes struct {
		Round               int                   `json:"round"`
		SuggestedPlayOrder  int                   `json:"suggested_play_order"`
		PointsByParticipant []PointsByParticipant `json:"points_by_participant"`
		Timestamps          TimeStamps            `json:"timestamps"`
	}

	TimeStamps struct {
		UnderwayAt time.Time `json:"underway_at"`
	}

	PointsByParticipant struct {
		ParticipantId int `json:"participant_id"`
	}

	Included struct {
		Id         string             `json:"id"`
		Type       string             `json:"type"`
		Attributes IncludedAttributes `json:"included"`
	}

	IncludedAttributes struct {
		Name string `json:"name"`
	}
)

// type AutoGenerated struct {
// 	Data []struct {
// 		ID         string `json:"id"`
// 		Type       string `json:"type"`
// 		Attributes struct {
// 			State               string        `json:"state"`
// 			Round               int           `json:"round"`
// 			Identifier          string        `json:"identifier"`
// 			Scores              string        `json:"scores"`
// 			SuggestedPlayOrder  int           `json:"suggested_play_order"`
// 			ScoreInSets         []interface{} `json:"score_in_sets"`
// 			PointsByParticipant []struct {
// 				ParticipantID int           `json:"participant_id"`
// 				Scores        []interface{} `json:"scores"`
// 			} `json:"points_by_participant"`
// 			Timestamps struct {
// 				StartedAt  time.Time `json:"started_at"`
// 				CreatedAt  time.Time `json:"created_at"`
// 				UpdatedAt  time.Time `json:"updated_at"`
// 				UnderwayAt time.Time `json:"underway_at"`
// 			} `json:"timestamps"`
// 			WinnerID interface{} `json:"winner_id"`
// 			Tie      bool        `json:"tie"`
// 		} `json:"attributes"`
// 		Relationships struct {
// 			Attachments struct {
// 				Data  []interface{} `json:"data"`
// 				Links struct {
// 					Related string `json:"related"`
// 					Meta    struct {
// 						Count int `json:"count"`
// 					} `json:"meta"`
// 				} `json:"links"`
// 			} `json:"attachments"`
// 			Station struct {
// 				Data struct {
// 					ID   string `json:"id"`
// 					Type string `json:"type"`
// 				} `json:"data"`
// 				Links struct {
// 					Related string `json:"related"`
// 				} `json:"links"`
// 			} `json:"station"`
// 		} `json:"relationships"`
// 	} `json:"data"`
// 	Included []struct {
// 		ID         string `json:"id"`
// 		Type       string `json:"type"`
// 		Attributes struct {
// 			ID            int    `json:"id"`
// 			Name          string `json:"name"`
// 			StreamURL     string `json:"stream_url"`
// 			Details       string `json:"details"`
// 			DetailsFormat string `json:"details_format"`
// 		} `json:"attributes"`
// 		Relationships struct {
// 			Match struct {
// 				Data struct {
// 					ID   string `json:"id"`
// 					Type string `json:"type"`
// 				} `json:"data"`
// 				Links struct {
// 					Related string `json:"related"`
// 				} `json:"links"`
// 			} `json:"match"`
// 			StationQueuers struct {
// 				Data  []interface{} `json:"data"`
// 				Links struct {
// 					Related string `json:"related"`
// 					Meta    struct {
// 						Count int `json:"count"`
// 					} `json:"meta"`
// 				} `json:"links"`
// 			} `json:"station_queuers"`
// 		} `json:"relationships"`
// 	} `json:"included"`
// 	Meta struct {
// 		Count int `json:"count"`
// 	} `json:"meta"`
// 	Links struct {
// 		Self string `json:"self"`
// 		Next string `json:"next"`
// 		Prev string `json:"prev"`
// 	} `json:"links"`
// }
