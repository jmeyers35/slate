package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	theOddsBaseURL = "https://api.the-odds-api.com/"

	sportNFL = "americanfootball_nfl"
)

type theOddsClient struct {
	client    *http.Client
	baseURL   string
	apiKey    string
	rateLimit time.Duration
	lastCall  time.Time
}

// TheOdds API response types
type oddsResponse struct {
	Success bool       `json:"success"`
	Data    []gameOdds `json:"data"`
}

type gameOdds struct {
	ID           string      `json:"id"`
	SportKey     string      `json:"sport_key"`
	CommenceTime time.Time   `json:"commence_time"`
	HomeTeam     string      `json:"home_team"`
	AwayTeam     string      `json:"away_team"`
	Bookmakers   []bookmaker `json:"bookmakers"`
}

type bookmaker struct {
	Key        string    `json:"key"`
	Title      string    `json:"title"`
	LastUpdate time.Time `json:"last_update"`
	Markets    []market  `json:"markets"`
}

const (
	marketKeySpreads string = "spreads"
	marketKeyH2H     string = "h2h"
	marketKeyTotals  string = "totals"
)

type market struct {
	Key      string    `json:"key"`
	Outcomes []outcome `json:"outcomes"`
}

type outcome struct {
	Name  string  `json:"name"`
	Price int     `json:"price"`
	Point float64 `json:"point,omitempty"`
}

func newTheOddsClient(config Config) (*theOddsClient, error) {
	return &theOddsClient{
		client:    &http.Client{Timeout: 10 * time.Second},
		baseURL:   theOddsBaseURL,
		apiKey:    config.APIKey,
		rateLimit: config.RateLimit,
	}, nil
}

func (c *theOddsClient) GetCurrentLines(ctx context.Context, season int, week int) ([]GameLines, error) {
	// Respect rate limiting
	if time.Since(c.lastCall) < c.rateLimit {
		time.Sleep(c.rateLimit - time.Since(c.lastCall))
	}

	// Make request to The Odds API
	url := fmt.Sprintf("%s/v4/sports/americanfootball_nfl/odds/?apikey=%s&regions=us&markets=spreads,totals,h2h",
		c.baseURL, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	var oddsResp oddsResponse
	if err := json.NewDecoder(resp.Body).Decode(&oddsResp); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	// Convert response to our GameLines type
	var lines []GameLines
	for _, game := range oddsResp.Data {
		// For now, we'll just use DraftKings odds if available
		// TODO: aggregate multiple bookmakers
		dkOdds := findBookmaker(game.Bookmakers, bookmakerDraftKings)
		if dkOdds == nil {
			continue
		}

		gl := GameLines{
			ProviderID:  game.ID,
			Timestamp:   time.Now(),
			Source:      bookmakerDraftKings,
			LastUpdated: dkOdds.LastUpdate,
		}

		// Parse markets for spread, moneyline, and totals
		for _, market := range dkOdds.Markets {
			switch market.Key {
			case marketKeySpreads:
				gl.HomeSpread = parseSpread(market.Outcomes, game.HomeTeam)
			case marketKeyH2H:
				parseMoneyline(&gl, market.Outcomes, game.HomeTeam)
			case marketKeyTotals:
				gl.OverUnder = parseTotal(market.Outcomes)
			}
		}

		// Calculate implied totals if we have both spread and over/under
		if gl.HomeSpread != 0 && gl.OverUnder != 0 {
			home := (gl.OverUnder + gl.HomeSpread) / 2
			away := (gl.OverUnder - gl.HomeSpread) / 2
			gl.HomeTeamTotal = &home
			gl.AwayTeamTotal = &away
		}

		lines = append(lines, gl)
	}

	c.lastCall = time.Now()
	return lines, nil
}

// GetLineHistory is not supported by the free tier of The Odds API
func (c *theOddsClient) GetLineHistory(ctx context.Context, gameID string) ([]GameLines, error) {
	return nil, fmt.Errorf("line history not supported by The Odds API free tier")
}

// GetGameID attempts to map our internal game ID to The Odds API game ID
func (c *theOddsClient) GetGameID(ctx context.Context, internalGameID string) (string, error) {
	// This would need to be implemented based on how we want to handle ID mapping
	// Could store mappings in database or try to match based on team names/time
	return "", fmt.Errorf("game ID mapping not implemented")
}

const (
	bookmakerDraftKings string = "draftkings"
)

func findBookmaker(bookmakers []bookmaker, key string) *bookmaker {
	for _, b := range bookmakers {
		if b.Key == key {
			return &b
		}
	}
	return nil
}

func parseSpread(outcomes []outcome, homeTeam string) float64 {
	for _, o := range outcomes {
		if o.Name == homeTeam {
			return o.Point
		}
	}
	return 0
}

func parseMoneyline(gl *GameLines, outcomes []outcome, homeTeam string) {
	for _, o := range outcomes {
		if o.Name == homeTeam {
			gl.HomeMoneyline = o.Price
		} else {
			gl.AwayMoneyline = o.Price
		}
	}
}

func parseTotal(outcomes []outcome) float64 {
	for _, o := range outcomes {
		if o.Name == "Over" {
			return o.Point
		}
	}
	return 0
}
