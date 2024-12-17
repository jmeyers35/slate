package scraper

import (
	"context"
	"fmt"

	"github.com/jmeyers35/slate/pkg/espn/client"
)

type ESPNActivities struct {
	Client client.Client
}

type GetPlayersForTeamRequest struct {
	TeamID client.TeamID
}

type GetPlayersForTeamResponse struct {
	Athletes []client.Athlete
}

func (a *ESPNActivities) GetPlayersForTeam(ctx context.Context, in GetPlayersForTeamRequest) (GetPlayersForTeamResponse, error) {
	roster, err := a.Client.GetRoster(ctx, in.TeamID)
	if err != nil {
		return GetPlayersForTeamResponse{}, fmt.Errorf("getting roster: %w", err)
	}

	var athletes []client.Athlete

	for _, positionGroup := range roster.Athletes {
		athletes = append(athletes, positionGroup.Athletes...)
	}
	return GetPlayersForTeamResponse{
		Athletes: athletes,
	}, nil
}

type GetTeamRequest struct {
	TeamID client.TeamID
}

type GetTeamResponse struct {
	Team client.Team
}

func (a *ESPNActivities) GetTeam(ctx context.Context, in GetTeamRequest) (GetTeamResponse, error) {
	team, err := a.Client.GetTeam(ctx, in.TeamID)
	if err != nil {
		return GetTeamResponse{}, fmt.Errorf("getting team: %w", err)
	}
	return GetTeamResponse{
		Team: team,
	}, nil
}

type GetTeamsRequest struct{}

type GetTeamsResponse struct {
	Teams []client.Team
}

func (a *ESPNActivities) GetTeamsFromESPN(ctx context.Context, in GetTeamsRequest) (GetTeamsResponse, error) {
	teamsResp, err := a.Client.GetTeams(ctx)
	if err != nil {
		return GetTeamsResponse{}, fmt.Errorf("getting teams: %w", err)
	}

	var teams []client.Team
	for _, sport := range teamsResp.Sports {
		for _, league := range sport.Leagues {
			teams = append(teams, league.Teams...)
		}
	}
	return GetTeamsResponse{
		Teams: teams,
	}, nil
}

type GetPlayerRequest struct {
	PlayerID string
}

type GetPlayerResponse struct {
	Player client.Athlete
}

func (a *ESPNActivities) GetPlayer(ctx context.Context, in GetPlayerRequest) (GetPlayerResponse, error) {
	player, err := a.Client.GetPlayer(ctx, in.PlayerID)
	if err != nil {
		return GetPlayerResponse{}, fmt.Errorf("getting player: %w", err)
	}
	return GetPlayerResponse{
		Player: player,
	}, nil
}

type GetScheduleRequest struct {
	Week   int
	Season int
}

type GetScheduleResponse struct {
	Schedule client.ScheduleResponse
}

func (a *ESPNActivities) GetSchedule(ctx context.Context, in GetScheduleRequest) (GetScheduleResponse, error) {
	schedule, err := a.Client.GetSchedule(ctx, in.Week, in.Season)
	if err != nil {
		return GetScheduleResponse{}, fmt.Errorf("getting schedule: %w", err)
	}
	return GetScheduleResponse{
		Schedule: schedule,
	}, nil
}
