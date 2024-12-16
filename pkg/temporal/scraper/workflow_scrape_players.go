package scraper

import (
	"fmt"

	"github.com/jmeyers35/slate/pkg/espn/client"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/multierr"
)

type ScrapePlayersForTeamRequest struct {
	TeamID   client.TeamID
	DBTeamID string
}

func ScrapePlayersForTeam(ctx workflow.Context, request ScrapePlayersForTeamRequest) error {
	var espnActivities *ESPNActivities

	var playersResponse GetPlayersForTeamResponse
	if err := workflow.ExecuteActivity(ctx, espnActivities.GetPlayersForTeam, GetPlayersForTeamRequest{
		TeamID: request.TeamID,
	}).Get(ctx, &playersResponse); err != nil {
		return fmt.Errorf("getting players for team: %w", err)
	}

	var futures []workflow.Future
	for _, player := range playersResponse.Athletes {
		futures = append(futures, workflow.ExecuteChildWorkflow(ctx, ScrapePlayer, ScrapePlayerRequest{
			PlayerESPNID: player.ID,
			TeamID:       request.DBTeamID,
		}))
	}

	var allErrors error

	for _, future := range futures {
		if err := future.Get(ctx, nil); err != nil {
			multierr.Append(allErrors, err)
		}
	}
	return nil
}
