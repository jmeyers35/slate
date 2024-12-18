package scraper

import (
	"fmt"
	"time"

	oddsactivities "github.com/jmeyers35/slate/pkg/odds/activities"
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
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
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

	var storageActivities *storageactivities.StorageActivities

	// Store the lines in the database
	err = workflow.ExecuteActivity(actx, storageActivities.UpsertLine, storageactivities.UpsertLinesRequest{
		Lines:  linesResp.Lines,
		Week:   request.Week,
		Season: request.Season,
	}).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("executing store game lines activity: %w", err)
	}

	return nil
}
