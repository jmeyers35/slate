package storage

import (
	"context"
	"time"
)

// Storage defines the interface for persisting NFL data
type Storage interface {
	// UpsertTeam creates or updates a team record in storage
	UpsertTeam(ctx context.Context, team *Team) error

	// GetTeamByESPNID retrieves a team by its ESPN ID
	GetTeamByESPNID(ctx context.Context, espnID string) (*Team, error)

	// GetTeams retrieves all teams from storage
	GetTeams(ctx context.Context) ([]*Team, error)

	// UpsertPlayer creates or updates a player record in storage
	UpsertPlayer(ctx context.Context, player *Player) error

	// UpsertGame creates or updates a game record in storage
	UpsertGame(ctx context.Context, game *Game) error

	// UpsertLine creates or updates a betting line record in storage
	UpsertLine(ctx context.Context, line Line) error

	// GetGameIDByTeams retrieves a game ID by matching the season, teams involved, and start time
	GetGameIDByTeams(ctx context.Context, season int, homeTeam, awayTeam string, start time.Time) (int, error)
}
