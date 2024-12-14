package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	baseSiteAPIURLFormat = "https://site.api.espn.com/apis/site/v2/sports/%s/%s"
	baseCoreAPIURLFormat = "https://sports.core.api.espn.com/v2/sports/%s/leagues/%s"
)

// Client is an interface for an ESPN API client
type Client interface {
	GetRoster(ctx context.Context, teamID string) (RosterResponse, error)
}

type ClientConfiguration struct {
	Sport  Sport
	League League
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

func (c *clientImpl) GetRoster(ctx context.Context, teamID string) (RosterResponse, error) {
	url := fmt.Sprintf("%s/teams/%s/roster", c.siteAPIURL, teamID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return RosterResponse{}, fmt.Errorf("creating request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return RosterResponse{}, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return RosterResponse{}, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

	var roster RosterResponse
	if err := json.NewDecoder(resp.Body).Decode(&roster); err != nil {
		return RosterResponse{}, fmt.Errorf("decoding response: %w", err)
	}

	return roster, nil
}

type RosterResponse struct {
	Timestamp string     `json:"timestamp"`
	Status    string     `json:"status"`
	Season    Season     `json:"season"`
	Athletes  []Athletes `json:"athletes"`
}

type Season struct {
	Year        int    `json:"year"`
	DisplayName string `json:"displayName"`
	// TODO: figure out what this enum means
	Type int    `json:"type"`
	Name string `json:"name"`
}

type Athletes struct {
	Position string    `json:"position"`
	Athletes []Athlete `json:"items"`
}

type Athlete struct {
	ID            string   `json:"id"`
	UID           string   `json:"uid"`
	GUID          string   `json:"guid"`
	FirstName     string   `json:"firstName"`
	LastName      string   `json:"lastName"`
	FullName      string   `json:"fullName"`
	DisplayName   string   `json:"displayName"`
	ShortName     string   `json:"shortName"`
	Weight        float32  `json:"weight"`
	DisplayWeight string   `json:"displayWeight"`
	Height        float32  `json:"height"`
	DisplayHeight string   `json:"displayHeight"`
	Age           int      `json:"age"`
	DateOfBirth   Time     `json:"dateOfBirth"`
	Slug          string   `json:"slug"`
	Jersey        string   `json:"jersey"`
	Position      Position `json:"position"`
	Injuries      []Injury `json:"injuries"`
}

type Position struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	Abbreviation string `json:"abbreviation"`
	// TODO: figure out what this means
	Leaf bool `json:"leaf"`
}

type College struct {
	ID           string `json:"id"`
	Mascot       string `json:"mascot"`
	Name         string `json:"name"`
	ShortName    string `json:"shortName"`
	Abbreviation string `json:"abbreviation"`
	// TODO: add logos when i feel less lazy
}

type Injury struct {
	// Status can be strings like: Out, Injured Reserve, Questionable, etc
	Status string `json:"status"`
	Date   Time   `json:"date"`
}

type Status struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Abbreviation string `json:"abbreviation"`
}

type Time struct {
	time.Time
}

func (ct *Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	// Try different formats
	formats := []string{
		"2006-01-02T07:00Z",
		"2006-01-02T15:04Z",
	}

	var lastErr error
	for _, format := range formats {
		t, err := time.Parse(format, s)
		if err == nil {
			ct.Time = t
			return nil
		}
		lastErr = err
	}

	return fmt.Errorf("could not parse time %q with any known format: %v", s, lastErr)
}
