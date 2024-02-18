package models

import (
	"errors"
	"net/url"
	"strings"
	"time"
)

const (
	TRAVELING_CONTROLLER = iota
	SNS
)

var (
	ErrorDateNotProvided          = errors.New("date query parameter not provided")
	ErrorDateIncorrectFormat      = errors.New("incorrect date format")
	ErrorTournamentOrgNotProvided = errors.New("no tournament organizer provided")
)

type (
	TournamentOrganizer int

	RequestValues struct {
		Date          string
		GameList      []string
		TournamentOrg TournamentOrganizer
	}
)

func CreateRequestValues(urlValues url.Values) (RequestValues, error) {

	dateStr := urlValues.Get("date")
	if dateStr == "" {
		return RequestValues{}, ErrorDateNotProvided
	}
	if _, err := time.Parse("2006-01-02", dateStr); err != nil {
		return RequestValues{}, ErrorDateIncorrectFormat
	}

	tournamentOrg := urlValues.Get("tournamentOrg")
	if tournamentOrg == "" {
		return RequestValues{}, ErrorTournamentOrgNotProvided
	}

	gamesListStr := urlValues.Get("games")
	var gamesList []string
	if gamesListStr != "" {
		gamesList = strings.Split(gamesListStr, ",")
	}

	return RequestValues{
		Date:     dateStr,
		GameList: gamesList,
		TournamentOrg: func(tournamentOrg string) TournamentOrganizer {
			if tournamentOrg == "sns" {
				return SNS
			}
			return TRAVELING_CONTROLLER
		}(tournamentOrg),
	}, nil
}
