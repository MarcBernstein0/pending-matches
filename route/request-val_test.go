package route

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateRequestVal(t *testing.T) {
	// Given
	tt := []struct {
		testName string
		mockData url.Values
		wantData RequestValues
		wantErr  error
	}{
		{
			testName: "no errors",
			mockData: url.Values{
				"date":          []string{"2006-01-02"},
				"tournamentOrg": []string{"traveling_controller"},
				"games":         []string{"game1,game2,game3"},
			},
			wantData: RequestValues{
				Date:          "2006-01-02",
				TournamentOrg: TRAVELING_CONTROLLER,
				GameList:      []string{"game1", "game2", "game3"},
			},
			wantErr: nil,
		},
		{
			testName: "no date provided",
			mockData: url.Values{
				"tournamentOrg": []string{"traveling_controller"},
				"games":         []string{"game1,game2,game3"},
			},
			wantData: RequestValues{},
			wantErr:  ErrorDateNotProvided,
		},
		{
			testName: "date not formatted correctly",
			mockData: url.Values{
				"date":          []string{"06-01-02"},
				"tournamentOrg": []string{"traveling_controller"},
				"games":         []string{"game1,game2,game3"},
			},
			wantData: RequestValues{},
			wantErr:  ErrorDateIncorrectFormat,
		},
		{
			testName: "tournament org not provided",
			mockData: url.Values{
				"date":  []string{"2006-01-02"},
				"games": []string{"game1,game2,game3"},
			},
			wantData: RequestValues{},
			wantErr:  ErrorTournamentOrgNotProvided,
		},
	}

	for _, tc := range tt {
		t.Run(tc.testName, func(t *testing.T) {
			// When
			gotData, gotErr := CreateRequestValues(tc.mockData)

			// Then
			assert.Equal(t, tc.wantData, gotData)
			if tc.wantErr != nil {
				assert.EqualError(t, gotErr, tc.wantErr.Error())
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
