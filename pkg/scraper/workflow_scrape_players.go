package scraper

import (
	"fmt"
	"time"

	espnactivities "github.com/jmeyers35/slate/pkg/espn/activities"
	"github.com/jmeyers35/slate/pkg/espn/client"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/multierr"
)

type ScrapePlayersForTeamRequest struct {
	TeamID   client.TeamID
	DBTeamID string
}

func ScrapePlayersForTeam(ctx workflow.Context, request ScrapePlayersForTeamRequest) error {
	var espnActivities *espnactivities.ESPNActivities

	var playersResponse espnactivities.GetPlayersForTeamResponse
	actx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: 1 * time.Minute,
		StartToCloseTimeout:    10 * time.Second,
	})
	if err := workflow.ExecuteActivity(actx, espnActivities.GetPlayersForTeam, espnactivities.GetPlayersForTeamRequest{
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
			allErrors = multierr.Append(allErrors, err)
		}
	}
	return allErrors
}
