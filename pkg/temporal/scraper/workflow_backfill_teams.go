package scraper

import (
	"fmt"
	"time"

	espnclient "github.com/jmeyers35/slate/pkg/espn/client"
	"go.temporal.io/sdk/workflow"
)

type BackfillTeamsRequest struct{}

func BackfillTeams(ctx workflow.Context, request BackfillTeamsRequest) error {
	var espnActivities *ESPNActivities

	var resp GetTeamsResponse
	actx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: 1 * time.Minute,
		StartToCloseTimeout:    10 * time.Second,
	})
	if err := workflow.ExecuteActivity(actx, espnActivities.GetTeamsFromESPN, GetTeamsRequest{}).Get(ctx, &resp); err != nil {
		return fmt.Errorf("getting teams: %w", err)
	}

	// Start a child workflow for each team to backfill into the DB
	for _, team := range resp.Teams {
		req := ScrapeNFLTeamRequest{
			TeamID: espnclient.TeamID(team.Team.ID),
		}
		if err := workflow.ExecuteChildWorkflow(ctx, ScrapeNFLTeam, req).Get(ctx, nil); err != nil {
			return fmt.Errorf("scraping team %s: %w", team.Team.Name, err)
		}
	}
	return nil
}
