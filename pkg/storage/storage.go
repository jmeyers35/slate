package storage

import "context"

type Storage interface {
	UpsertTeam(ctx context.Context, team *Team) error
	GetTeamByESPNID(ctx context.Context, espnID string) (*Team, error)

	UpsertPlayer(ctx context.Context, player *Player) error
}
