package scraper

import (
	"fmt"
	"time"

	"github.com/jmeyers35/slate/pkg/converters"
	espnclient "github.com/jmeyers35/slate/pkg/espn/client"
	"go.temporal.io/sdk/workflow"
)

type ScrapeNFLTeamRequest struct {
	TeamID espnclient.TeamID
}

// ScrapeNFLTeam is a Temporal workflow that scrapes NFL team data for a particular team
// and writes it to storage.
func ScrapeNFLTeam(ctx workflow.Context, request ScrapeNFLTeamRequest) error {
	logger := workflow.GetLogger(ctx)
	actx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: 30 * time.Second,
		StartToCloseTimeout:    5 * time.Second,
	})

	var espnActivities *ESPNActivities
	req := GetTeamRequest{
		TeamID: request.TeamID,
	}

	var getTeamResponse GetTeamResponse
	if err := workflow.ExecuteActivity(actx, espnActivities.GetTeam, req).Get(ctx, &getTeamResponse); err != nil {
		return fmt.Errorf("getting team: %w", err)
	}

	logger.Info("Got team", "team", getTeamResponse.Team)

	var storageActivities *StorageActivities
	converter := converters.ESPNAPIConverter{}
	teamReq := UpsertTeamRequest{
		Team: converter.ConvertTeam(getTeamResponse.Team),
	}
	if err := workflow.ExecuteActivity(actx, storageActivities.UpsertTeam, teamReq).Get(ctx, nil); err != nil {
		return fmt.Errorf("upserting team: %w", err)
	}
	return nil
}
