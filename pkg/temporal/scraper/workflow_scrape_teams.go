package scraper

import (
	"fmt"
	"time"

	espnclient "github.com/jmeyers35/slate/pkg/espn/client"
	"go.temporal.io/sdk/workflow"
)

type ScrapeNFLTeamRequest struct {
	TeamID espnclient.TeamID
}

// ScrapeNFLTeam is a Temporal workflow that scrapes NFL team data for a particular team.
func ScrapeNFLTeam(ctx workflow.Context, request ScrapeNFLTeamRequest) error {
	logger := workflow.GetLogger(ctx)
	actx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: 30 * time.Second,
	})

	var activities *ESPNActivities
	req := GetPlayersForTeamRequest{
		TeamID: request.TeamID,
	}
	var athletesResp GetPlayersForTeamResponse

	if err := workflow.ExecuteActivity(actx, activities.GetPlayersForTeam, req).Get(ctx, &athletesResp); err != nil {
		return fmt.Errorf("getting players for team: %w", err)
	}

	logger.Info("Got athletes", "athletes", athletesResp.Athletes)
	return nil
}
