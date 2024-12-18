package activities

import (
	"context"
	"fmt"
	"time"

	"github.com/jmeyers35/slate/pkg/storage"
	"go.uber.org/multierr"

	"github.com/jmeyers35/slate/pkg/odds/client"
)

// StoreGameLinesRequest is the request for storing game lines
type StoreGameLinesRequest struct {
	Lines  []client.GameLines
	Week   int
	Season int
}

// StoreGameLines stores game lines in the database
func (a *StorageActivities) StoreGameLines(ctx context.Context, req StoreGameLinesRequest) error {
	var err error
	for _, line := range req.Lines {
		// Get internal game ID based on season, week, and team IDs
		gameID, err := a.Storage.GetGameIDByTeams(ctx, req.Season, req.Week, line.HomeTeamName, line.AwayTeamName)
		if err != nil {
			return fmt.Errorf("getting game ID: %w", err)
		}

		storageLine := storage.Line{
			GameID:        gameID,
			ProviderID:    line.ProviderID,
			HomeSpread:    line.HomeSpread,
			OverUnder:     line.OverUnder,
			HomeMoneyline: line.HomeMoneyline,
			AwayMoneyline: line.AwayMoneyline,
			LastUpdated:   line.LastUpdated,
			CreatedAt:     time.Now(),
		}

		// Set team totals if available
		if line.HomeTeamTotal != nil {
			storageLine.HomeTeamTotal = *line.HomeTeamTotal
		}
		if line.AwayTeamTotal != nil {
			storageLine.AwayTeamTotal = *line.AwayTeamTotal
		}

		if err := a.Storage.UpsertLine(ctx, storageLine); err != nil {
			err = multierr.Append(err, fmt.Errorf("upserting game line: %w", err))
		}
	}
	return err
}
