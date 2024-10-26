package nag

import (
	"encoding/json"
	"fmt"
	"nba-sql/nag/params"
	"net/http"
)

// LeagueStandingsV3 wraps request to and response from leaguestandingsv3 endpoint.
type LeagueStandingsV3 struct {
	*Client
	LeagueID   string
	Season     string
	SeasonType params.SeasonType

	Response *Response
}

// NewLeagueStandingsV3 creates a default LeagueStandingsV3 instance.
func NewLeagueStandingsV3() *LeagueStandingsV3 {
	return &LeagueStandingsV3{
		Client: NewDefaultClient(),

		LeagueID: params.LeagueID.Default(),

		Season:     params.CurrentSeason,
		SeasonType: params.DefaultSeasonType,
	}
}

// Get sends a GET request to leaguestandingsv3 endpoint.
func (c *LeagueStandingsV3) Get() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/leaguestandingsv3", c.BaseURL.String()), nil)
	if err != nil {
		return err
	}

	req.Header = DefaultStatsHeader

	q := req.URL.Query()
	q.Add("LeagueID", c.LeagueID)
	q.Add("Season", c.Season)
	q.Add("SeasonType", string(c.SeasonType))

	req.URL.RawQuery = q.Encode()

	b, err := c.Do(req)
	if err != nil {
		return err
	}

	var res Response
	if err := json.Unmarshal(b, &res); err != nil {
		return err
	}

	c.Response = &res
	return nil
}
