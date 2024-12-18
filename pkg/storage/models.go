package storage

import "time"

type Team struct {
	ID       string
	TeamCode string
	Name     string
	ESPNID   string
}

type Player struct {
	ID       string
	Name     string
	Position string
	TeamID   string
	ESPNID   string
}

type Game struct {
	ID         string
	Week       int
	Season     int
	HomeTeamID string
	AwayTeamID string
	GameDate   time.Time
	Dome       bool
}

type Line struct {
	LineID        string
	GameID        string
	ProviderID    string
	HomeSpread    float64
	OverUnder     float64
	HomeTeamTotal float64
	AwayTeamTotal float64
	HomeMoneyline int
	AwayMoneyline int
	LastUpdated   time.Time
	CreatedAt     time.Time
}
