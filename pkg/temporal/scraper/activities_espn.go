package scraper

import (
	"context"
	"fmt"

	espnclient "github.com/jmeyers35/slate/pkg/espn/client"
)

type ESPNActivities struct {
	client espnclient.Client
}

type GetPlayersForTeamRequest struct {
	TeamID espnclient.TeamID
}

type GetPlayersForTeamResponse struct {
	Athletes []espnclient.Athlete
}

func (a *ESPNActivities) GetPlayersForTeam(ctx context.Context, in GetPlayersForTeamRequest) (GetPlayersForTeamResponse, error) {
	roster, err := a.client.GetRoster(ctx, in.TeamID)
	if err != nil {
		return GetPlayersForTeamResponse{}, fmt.Errorf("getting roster: %w", err)
	}

	var athletes []espnclient.Athlete

	for _, positionGroup := range roster.Athletes {
		athletes = append(athletes, positionGroup.Athletes...)
	}
	return GetPlayersForTeamResponse{
		Athletes: athletes,
	}, nil
}

type GetTeamRequest struct {
	TeamID espnclient.TeamID
}

type GetTeamResponse struct {
	Team espnclient.Team
}

func (a *ESPNActivities) GetTeam(ctx context.Context, in GetTeamRequest) (GetTeamResponse, error) {
	team, err := a.client.GetTeam(ctx, in.TeamID)
	if err != nil {
		return GetTeamResponse{}, fmt.Errorf("getting team: %w", err)
	}
	return GetTeamResponse{
		Team: team,
	}, nil
}

type GetTeamsRequest struct{}

type GetTeamsResponse struct {
	Teams []espnclient.Team
}

func (a *ESPNActivities) GetTeamsFromESPN(ctx context.Context, in GetTeamsRequest) (GetTeamsResponse, error) {
	teamsResp, err := a.client.GetTeams(ctx)
	if err != nil {
		return GetTeamsResponse{}, fmt.Errorf("getting teams: %w", err)
	}

	var teams []espnclient.Team
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
	Player espnclient.Athlete
}

func (a *ESPNActivities) GetPlayer(ctx context.Context, in GetPlayerRequest) (GetPlayerResponse, error) {
	player, err := a.client.GetPlayer(ctx, in.PlayerID)
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
	Schedule espnclient.ScheduleResponse
}

func (a *ESPNActivities) GetSchedule(ctx context.Context, in GetScheduleRequest) (GetScheduleResponse, error) {
	schedule, err := a.client.GetSchedule(ctx, in.Week, in.Season)
	if err != nil {
		return GetScheduleResponse{}, fmt.Errorf("getting schedule: %w", err)
	}
	return GetScheduleResponse{
		Schedule: schedule,
	}, nil
}
