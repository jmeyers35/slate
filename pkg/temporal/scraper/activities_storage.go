package scraper

import (
	"context"

	"github.com/jmeyers35/slate/pkg/storage"
)

type StorageActivities struct {
	Storage storage.Storage
}

type UpsertTeamRequest struct {
	Team storage.Team
}

func (a *StorageActivities) UpsertTeam(ctx context.Context, req UpsertTeamRequest) error {
	return a.Storage.UpsertTeam(ctx, &req.Team)
}

type GetTeamByESPNIDRequest struct {
	ESPNID string
}

func (a *StorageActivities) GetTeamByESPNID(ctx context.Context, req GetTeamByESPNIDRequest) (*storage.Team, error) {
	return a.Storage.GetTeamByESPNID(ctx, req.ESPNID)
}

type UpsertPlayerRequest struct {
	Player storage.Player
}

func (a *StorageActivities) UpsertPlayer(ctx context.Context, req UpsertPlayerRequest) error {
	return a.Storage.UpsertPlayer(ctx, &req.Player)
}
