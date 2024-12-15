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
