package scraper

import (
	"fmt"
	"time"

	espnactivities "github.com/jmeyers35/slate/pkg/espn/activities"
	storageactivities "github.com/jmeyers35/slate/pkg/storage/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type ScrapeScheduleRequest struct {
	Week   int
	Season int
}

func ScrapeSchedule(ctx workflow.Context, request ScrapeScheduleRequest) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting schedule scrape workflow", "week", request.Week, "season", request.Season)

	activityOpts := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	actx := workflow.WithActivityOptions(ctx, activityOpts)

	var espnActivities *espnactivities.ESPNActivities

	var scheduleResp espnactivities.GetScheduleResponse
	err := workflow.ExecuteActivity(actx, espnActivities.GetSchedule, espnactivities.GetScheduleRequest{
		Week:   request.Week,
		Season: request.Season,
	}).Get(ctx, &scheduleResp)
	if err != nil {
		return fmt.Errorf("executing get schedule activity: %w", err)
	}

	var storageActivities *storageactivities.StorageActivities

	err = workflow.ExecuteActivity(actx, storageActivities.StoreGames, storageactivities.StoreGamesRequest{
		Schedule: scheduleResp.Schedule,
		Week:     request.Week,
		Season:   request.Season,
	}).Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("executing store games activity: %w", err)
	}

	return nil
}
