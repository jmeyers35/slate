package scraper

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type BackfillScheduleRequest struct {
	Season int
}

func BackfillSchedule(ctx workflow.Context, request BackfillScheduleRequest) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("Starting schedule backfill workflow", "season", request.Season)

	// Child workflow options - longer timeout since we're running multiple weeks
	childOpts := workflow.ChildWorkflowOptions{
		WorkflowRunTimeout: time.Hour,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    3,
		},
	}
	ctx = workflow.WithChildOptions(ctx, childOpts)

	// Regular season is weeks 1-18
	const regularSeasonWeeks = 18

	// Process each week serially
	for week := 1; week <= regularSeasonWeeks; week++ {
		logger.Info("Processing week", "week", week)

		err := workflow.ExecuteChildWorkflow(ctx, ScrapeSchedule, ScrapeScheduleRequest{
			Week:   week,
			Season: request.Season,
		}).Get(ctx, nil)

		if err != nil {
			return fmt.Errorf("week %d scrape failed: %w", week, err)
		}
	}

	logger.Info("Schedule backfill complete", "season", request.Season)
	return nil
}
