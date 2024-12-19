package activities

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/multierr"

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

type UpsertGameRequest struct {
	Schedule espnclient.ScheduleResponse
	Week     int
	Season   int
}

func (a *StorageActivities) UpsertGame(ctx context.Context, req UpsertGameRequest) error {
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

// UpsertLinesRequest is the request for storing a vegas line for a game
type UpsertLinesRequest struct {
	Lines []storage.Line
}

// UpsertLine stores game lines in the database
func (a *StorageActivities) UpsertLine(ctx context.Context, req UpsertLinesRequest) error {
	var storeErrs error
	for _, line := range req.Lines {
		if err := a.Storage.UpsertLine(ctx, line); err != nil {
			storeErrs = multierr.Append(storeErrs, fmt.Errorf("upserting game line %v: %w", line, err))
		}
	}
	return storeErrs
}

type GetGameIDByTeamsRequest struct {
	Season   int
	HomeTeam string
	AwayTeam string
	Start    time.Time
}

type GetGameIDByTeamsResponse struct {
	GameID int
}

func (a *StorageActivities) GetGameIDByTeams(ctx context.Context, req GetGameIDByTeamsRequest) (GetGameIDByTeamsResponse, error) {
	gameID, err := a.Storage.GetGameIDByTeams(ctx, req.Season, req.HomeTeam, req.AwayTeam, req.Start)
	if err != nil {
		return GetGameIDByTeamsResponse{}, fmt.Errorf("getting game ID: %w", err)
	}
	return GetGameIDByTeamsResponse{
		GameID: gameID,
	}, nil
}
