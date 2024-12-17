package storage

import (
	"context"
	"fmt"

	"github.com/jmeyers35/slate/pkg/converters"
	espnclient "github.com/jmeyers35/slate/pkg/espn/client"
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

func (a *StorageActivities) GetTeamsFromStorage(ctx context.Context, req GetTeamsFromStorageRequest) (GetTeamsFromStorageResponse, error) {
	teams, err := a.Storage.GetTeams(ctx)
	if err != nil {
		return GetTeamsFromStorageResponse{}, fmt.Errorf("getting teams: %w", err)
	}
	return GetTeamsFromStorageResponse{
		Teams: teams,
	}, nil
}

// New game storage activity
type StoreGamesRequest struct {
	Schedule espnclient.ScheduleResponse
	Week     int
	Season   int
}

func (a *StorageActivities) StoreGames(ctx context.Context, req StoreGamesRequest) error {
	converter := converters.ESPNAPIConverter{}

	for _, event := range req.Schedule.Events {
		game := converter.ConvertGame(event, req.Week, req.Season)

		// Get team IDs from storage
		homeTeam, err := a.Storage.GetTeamByESPNID(ctx, game.HomeTeamID)
		if err != nil {
			return fmt.Errorf("getting home team: %w", err)
		}
		if homeTeam == nil {
			return fmt.Errorf("home team with ESPN ID %s not found - ensure teams are bootstrapped first", game.HomeTeamID)
		}

		awayTeam, err := a.Storage.GetTeamByESPNID(ctx, game.AwayTeamID)
		if err != nil {
			return fmt.Errorf("getting away team: %w", err)
		}
		if awayTeam == nil {
			return fmt.Errorf("away team with ESPN ID %s not found - ensure teams are bootstrapped first", game.AwayTeamID)
		}

		// Update with internal IDs
		game.HomeTeamID = homeTeam.ID
		game.AwayTeamID = awayTeam.ID

		if err := a.Storage.UpsertGame(ctx, &game); err != nil {
			return fmt.Errorf("upserting game: %w", err)
		}
	}
	return nil
}
