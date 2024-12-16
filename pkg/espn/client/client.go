package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseSiteAPIURLFormat = "https://site.api.espn.com/apis/site/v2/sports/%s/%s"
	baseCoreAPIURLFormat = "https://sports.core.api.espn.com/v2/sports/%s/leagues/%s"
)

// Client is an interface for an ESPN API client
type Client interface {
	GetRoster(ctx context.Context, teamID TeamID) (RosterResponse, error)
	GetTeams(ctx context.Context) (TeamsResponse, error)
	GetTeam(ctx context.Context, teamID TeamID) (Team, error)
	GetPlayer(ctx context.Context, playerID string) (Athlete, error)
}

func New(config ClientConfiguration) Client {
	return &clientImpl{
		httpClient: &http.Client{},
		siteAPIURL: fmt.Sprintf(baseSiteAPIURLFormat, config.Sport, config.League),
		coreAPIURL: fmt.Sprintf(baseCoreAPIURLFormat, config.Sport, config.League),
	}
}

func NewNFL() Client {
	return New(ClientConfiguration{
		Sport:  SportFootball,
		League: LeagueNFL,
	})
}

type clientImpl struct {
	httpClient *http.Client
	siteAPIURL string
	coreAPIURL string
}

var _ Client = &clientImpl{}

func (c *clientImpl) GetRoster(ctx context.Context, teamID TeamID) (RosterResponse, error) {
	url := fmt.Sprintf("%s/teams/%s/roster", c.siteAPIURL, teamID)
	var roster RosterResponse
	if err := c.doHTTPRequest(ctx, url, http.MethodGet, &roster); err != nil {
		return RosterResponse{}, fmt.Errorf("getting roster: %w", err)
	}
	return roster, nil
}

func (c *clientImpl) GetTeams(ctx context.Context) (TeamsResponse, error) {
	url := fmt.Sprintf("%s/teams", c.siteAPIURL)
	var teams TeamsResponse
	if err := c.doHTTPRequest(ctx, url, http.MethodGet, &teams); err != nil {
		return TeamsResponse{}, fmt.Errorf("getting teams: %w", err)
	}
	return teams, nil
}

func (c *clientImpl) GetTeam(ctx context.Context, teamID TeamID) (Team, error) {
	url := fmt.Sprintf("%s/teams/%s", c.siteAPIURL, teamID)
	var team Team
	if err := c.doHTTPRequest(ctx, url, http.MethodGet, &team); err != nil {
		return Team{}, fmt.Errorf("getting team: %w", err)
	}
	return team, nil
}

func (c *clientImpl) GetPlayer(ctx context.Context, playerID string) (Athlete, error) {
	url := fmt.Sprintf("%s/athletes/%s", c.coreAPIURL, playerID)
	var player Athlete
	if err := c.doHTTPRequest(ctx, url, http.MethodGet, &player); err != nil {
		return Athlete{}, fmt.Errorf("getting player: %w", err)
	}
	return player, nil
}

func (c *clientImpl) doHTTPRequest(ctx context.Context, url string, method string, respDataPointer interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	if err := json.NewDecoder(resp.Body).Decode(respDataPointer); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	return nil
}
