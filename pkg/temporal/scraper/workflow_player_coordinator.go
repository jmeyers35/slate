package scraper

import (
	"fmt"
	"time"

	"github.com/jmeyers35/slate/pkg/espn/client"
	storageactivities "github.com/jmeyers35/slate/pkg/storage/activities"
	"go.temporal.io/sdk/workflow"
)

type PlayerCoordinatorRequest struct{}

func PlayerCoordinator(ctx workflow.Context, request PlayerCoordinatorRequest) error {
	var storageActivities *storageactivities.StorageActivities

	var getTeamsResp storageactivities.GetTeamsFromStorageResponse
	actx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: 1 * time.Minute,
		StartToCloseTimeout:    30 * time.Second,
	})

	if err := workflow.ExecuteActivity(actx, storageActivities.GetTeamsFromStorage, nil).Get(ctx, &getTeamsResp); err != nil {
		return fmt.Errorf("getting teams: %w", err)
	}

	for _, team := range getTeamsResp.Teams {
		if err := workflow.ExecuteChildWorkflow(ctx, ScrapePlayersForTeam, ScrapePlayersForTeamRequest{
			TeamID:   client.TeamID(team.ESPNID), // TODO: less jank?
			DBTeamID: team.ID,
		}).Get(ctx, nil); err != nil {
			return fmt.Errorf("scraping players for team: %w", err)
		}
	}

	return nil
}
