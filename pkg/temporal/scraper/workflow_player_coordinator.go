package scraper

import (
	"fmt"

	"github.com/jmeyers35/slate/pkg/espn/client"
	"go.temporal.io/sdk/workflow"
)

type PlayerCoordinatorRequest struct{}

func PlayerCoordinator(ctx workflow.Context, request PlayerCoordinatorRequest) error {
	var storageActivities *StorageActivities

	var getTeamsResp GetTeamsFromStorageResponse

	if err := workflow.ExecuteActivity(ctx, storageActivities.GetTeams, nil).Get(ctx, &getTeamsResp); err != nil {
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
