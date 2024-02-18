package models

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	TRAVELING_CONTROLLER = iota
	SNS
)

var (
	ErrorDateNotProvided     = errors.New("date query parameter not provided")
	ErrorDateIncorrectFormat = errors.New("incorrect date format")
)

type (
	TournamentOrganizer int

	RequestValues struct {
		Date     string
		GameList []string
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

	gamesListStr := urlValues.Get("games")
	fmt.Println("gamesList", gamesListStr)
	var gamesList []string
	if gamesListStr != "" {
		gamesList = strings.Split(gamesListStr, ",")
	}

	return RequestValues{
		Date:     dateStr,
		GameList: gamesList,
	}, nil
}
