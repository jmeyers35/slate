package scraper

import (
	"context"
	"fmt"

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

type GetTeamsFromStorageRequest struct{}

type GetTeamsFromStorageResponse struct {
	Teams []*storage.Team
}

func (a *StorageActivities) GetTeams(ctx context.Context, req GetTeamsRequest) (GetTeamsFromStorageResponse, error) {
	teams, err := a.Storage.GetTeams(ctx)
	if err != nil {
		return GetTeamsFromStorageResponse{}, fmt.Errorf("getting teams: %w", err)
	}
	return GetTeamsFromStorageResponse{
		Teams: teams,
	}, nil
}
