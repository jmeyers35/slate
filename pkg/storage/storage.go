package storage

import (
	"context"
	"time"
)

type Storage interface {
	UpsertTeam(ctx context.Context, team *Team) error
	GetTeamByESPNID(ctx context.Context, espnID string) (*Team, error)
	GetTeams(ctx context.Context) ([]*Team, error)

	UpsertPlayer(ctx context.Context, player *Player) error

	UpsertGame(ctx context.Context, game *Game) error

	UpsertLine(ctx context.Context, line Line) error

	GetGameIDByTeams(ctx context.Context, season int, homeTeam, awayTeam string, start time.Time) (int, error)
}
