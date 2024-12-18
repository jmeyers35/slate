package client

import "time"

// GameLines represents odds data for a single game
type GameLines struct {
	// GameID is our internal game ID that maps to the games table
	GameID string
	// ProviderID is the game ID from the odds provider
	ProviderID string
	// HomeTeamName is the name of the home team
	HomeTeamName string
	// AwayTeamName is the name of the away team
	AwayTeamName string
	// Timestamp indicates when these odds were captured
	Timestamp time.Time

	// HomeSpread represents the spread for the home team. Negative means home team is favored
	HomeSpread float64
	// HomeMoneyline is the moneyline for the home team in American odds format (+150, -110, etc)
	HomeMoneyline int
	// AwayMoneyline is the moneyline for the away team in American odds format
	AwayMoneyline int
	// OverUnder represents the total points line
	OverUnder float64

	// HomeTeamTotal is the derived team total based on spread and over/under
	HomeTeamTotal *float64
	// AwayTeamTotal is the derived team total based on spread and over/under
	AwayTeamTotal *float64

	// Source indicates which book these lines are from
	Source string
	// LastUpdated indicates when the lines were last updated at the source
	LastUpdated time.Time
}
