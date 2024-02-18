package challongebracketmatches

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/MarcBernstein0/pending-matches/models"
)

var (
	ErrResponseNotOK error = errors.New("response not ok")
	ErrServerProblem error = errors.New("server error")
	ErrNoData        error = errors.New("no data found")
)

type (
	customClient struct {
		baseURL string
		client  *http.Client
		// apiKeyTravCntlr string
		// apiKeySNS       string
		apiKey string
	}

	FetchData interface {
		// FetchTournaments fetch all tournaments created after a specific date
		// GET https://api.challonge.com/v2.1/tournaments.json?page={}&per_page=25
		FetchTournaments(date string) (map[string]string, error)
		// FetchParticipants fetch all participants for a tournament
		// GET https://api.challonge.com/v2.1/tournaments/{tournaments}/participants.json?page={}&per_page=25
		FetchParticipants(tournamentId, tournamentGame string) (models.TournamentParticipants, error)
		// FetchMatches fetch matches for a tournament
		// GET https://api.challonge.com/v2.1/tournaments/{tournaments}/matches.json?page=1&per_page=25&state=open
		FetchMatches(tournamentParticipants models.TournamentParticipants) (models.TournamentMatches, error)
	}
)

// func New(baseURL, apiKeyTravCntlr, apiKeySNS string, client *http.Client, contextTimeout time.Duration) *customClient {
// 	return &customClient{
// 		baseURL:         baseURL,
// 		client:          client,
// 		apiKeyTravCntlr: apiKeyTravCntlr,
// 		apiKeySNS:       apiKeySNS,
// 	}
// }

func New(baseURL, apiKey string, client *http.Client, contextTimeout time.Duration) *customClient {
	return &customClient{
		baseURL: baseURL,
		client:  client,
		apiKey:  apiKey,
	}
}

// Return map of type int -> string where int is the tournamentId and string is the game name
func (c *customClient) FetchTournaments(date string) (map[string]string, error) {

	resMap := make(map[string]string)

	// dealing with paginated response
	paginationLeft := true
	pageNumber := 1

	for paginationLeft {
		params := map[string]string{
			"state":         "in_progress",
			"created_after": date,
			"page":          strconv.Itoa(pageNumber),
			"per_page":      "25",
		}

		res, err := c.get(http.MethodGet, c.baseURL+"/tournaments.json", nil, params)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode))
		}

		defer res.Body.Close()

		var tournaments models.Tournaments
		err = json.NewDecoder(res.Body).Decode(&tournaments)
		if err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("%w. %s", err, http.StatusText(http.StatusInternalServerError))
		}
		if len(tournaments.Data) == 0 {
			paginationLeft = false
		} else {
			for _, tournament := range tournaments.Data {
				resMap[tournament.Id] = tournament.Attributes.GameName

			}

			pageNumber++
		}
	}

	return resMap, nil
}

// Return a models.TournamentParticipants with a map of participants ids -> participant tags
func (c *customClient) FetchParticipants(tournamentId, tournamentGame string) (models.TournamentParticipants, error) {
	participants := models.TournamentParticipants{
		GameName:     tournamentGame,
		TournamentID: tournamentId,
		Participant:  map[string]string{},
	}

	// dealing with paginated responses
	paginationLeft := true
	pageNumber := 1

	for paginationLeft {
		params := map[string]string{
			"page":     strconv.Itoa(pageNumber),
			"per_page": "25",
		}

		res, err := c.get(http.MethodGet, c.baseURL+"/tournaments/"+tournamentId+"/participants.json", nil, params)
		if err != nil {
			return models.TournamentParticipants{}, err
		}

		if res.StatusCode != http.StatusOK {
			return models.TournamentParticipants{}, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode))
		}

		defer res.Body.Close()

		var participantsChall models.Participants
		err = json.NewDecoder(res.Body).Decode(&participantsChall)
		if err != nil {
			return models.TournamentParticipants{}, fmt.Errorf("%w. %s", err, http.StatusText(http.StatusInternalServerError))
		}

		if len(participantsChall.Data) == 0 {
			paginationLeft = false
		} else {
			for _, participant := range participantsChall.Data {
				participants.Participant[participant.Id] = participant.Attributes.Name
			}
			pageNumber++
		}
	}

	return participants, nil
}

// Return a models.TournamentMatches with a list of Match structs that include player names
func (c *customClient) FetchMatches(tournamentParticipants models.TournamentParticipants) (models.TournamentMatches, error) {
	matchResult := models.TournamentMatches{
		GameName:     tournamentParticipants.GameName,
		TournamentId: tournamentParticipants.TournamentID,
		MatchList:    []models.Match{},
	}

	params := map[string]string{
		"page":     "1",
		"per_page": "50",
		"state":    "open",
	}

	res, err := c.get(http.MethodGet, c.baseURL+"/tournaments/"+matchResult.TournamentId+"/matches.json", nil, params)
	if err != nil {
		return models.TournamentMatches{}, err
	}

	if res.StatusCode != http.StatusOK {
		return models.TournamentMatches{}, fmt.Errorf("%w. %s", ErrResponseNotOK, http.StatusText(res.StatusCode))
	}

	defer res.Body.Close()
	var matches models.Matches
	err = json.NewDecoder(res.Body).Decode(&matches)
	if err != nil {
		return models.TournamentMatches{}, fmt.Errorf("%w. %s", err, http.StatusText(http.StatusInternalServerError))
	}

	stationsMap := getStationsMap(matches)

	if len(matches.Data) == 0 {
		return matchResult, nil
	}

	for _, match := range matches.Data {
		matchData := models.Match{
			Id:                 match.Id,
			Player1Name:        tournamentParticipants.Participant[strconv.Itoa(match.Attributes.PointsByParticipant[0].ParticipantId)],
			Player2Name:        tournamentParticipants.Participant[strconv.Itoa(match.Attributes.PointsByParticipant[1].ParticipantId)],
			Round:              match.Attributes.Round,
			SuggestedPlayOrder: match.Attributes.SuggestedPlayOrder,
			Underway:           !match.Attributes.Timestamps.UnderwayAt.IsZero(),
			Station:            stationsMap[match.Relationship.Station.Data.Id],
		}
		matchResult.MatchList = append(matchResult.MatchList, matchData)
	}

	// sort matches based on SuggestedPlayOrder
	sort.Slice(matchResult.MatchList, func(i, j int) bool {
		return matchResult.MatchList[i].SuggestedPlayOrder <= matchResult.MatchList[j].SuggestedPlayOrder
	})

	return matchResult, nil
}

func (c *customClient) get(method, urlPath string, reqBody io.Reader, params map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, urlPath, reqBody)
	if err != nil {
		// gracefully handle error and pass along
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/vnd.api+json")
	req.Header.Add("Authorization-Type", "v1")
	req.Header.Add("Authorization", c.apiKey)

	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	return c.client.Do(req)
}

// Return a map of station id(string) -> station name(string)
func getStationsMap(matches models.Matches) map[string]string {
	stationMap := make(map[string]string)
	for _, includedInfo := range matches.Included {
		if includedInfo.Type == "station" {
			stationMap[includedInfo.Id] = includedInfo.Attributes.Name
		}
	}

	return stationMap
}
