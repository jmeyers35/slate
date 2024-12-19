package scraper

import (
	"fmt"
	"time"

	"github.com/jmeyers35/slate/pkg/converters"
	oddsactivities "github.com/jmeyers35/slate/pkg/odds/activities"
	"github.com/jmeyers35/slate/pkg/storage"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"

	storageactivities "github.com/jmeyers35/slate/pkg/storage/activities"
)

type ScrapeOddsRequest struct {
	Week   int
	Season int
}

func ScrapeOdds(ctx workflow.Context, request ScrapeOddsRequest) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting odds scrape workflow", "week", request.Week, "season", request.Season)

	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: 5 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 1,
		},
	}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	var oddsActivities *oddsactivities.OddsActivities

	// Get current lines from the odds provider
	var linesResp oddsactivities.GetCurrentLinesResponse
	err := workflow.ExecuteActivity(actx, oddsActivities.GetCurrentLines, oddsactivities.GetCurrentLinesRequest{
		Week:   request.Week,
		Season: request.Season,
	}).Get(ctx, &linesResp)
	if err != nil {
		return fmt.Errorf("executing get current lines activity: %w", err)
	}

	logger.Info("Got current lines", "lines", linesResp.Lines)

	var storageActivities *storageactivities.StorageActivities

	var convertedLines []storage.Line
	var converter converters.TheOddsConverter

	for _, line := range linesResp.Lines {
		var getGameIDResponse storageactivities.GetGameIDByTeamsResponse
		if err := workflow.ExecuteActivity(actx, storageActivities.GetGameIDByTeams, storageactivities.GetGameIDByTeamsRequest{
			Season:   request.Season,
			HomeTeam: line.HomeTeamName,
			AwayTeam: line.AwayTeamName,
			Start:    line.GameTime,
		}).Get(ctx, &getGameIDResponse); err != nil {
			return fmt.Errorf("getting game ID by teams for line %v: %w", line, err)
		}
		convertedLines = append(convertedLines, converter.ConvertLine(line, getGameIDResponse.GameID))
	}

	// Store the lines in the database
	err = workflow.ExecuteActivity(actx, storageActivities.UpsertLine, storageactivities.UpsertLinesRequest{
		Lines: convertedLines,
	}).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("executing store game lines activity: %w", err)
	}

	return nil
}
