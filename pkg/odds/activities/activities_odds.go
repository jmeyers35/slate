package scraper

import (
	"context"
	"fmt"

	"github.com/jmeyers35/slate/pkg/odds/client"
)

type OddsActivities struct {
	Client client.Client
}

type GetCurrentLinesRequest struct {
	Season int
	Week   int
}

type GetCurrentLinesResponse struct {
	Lines []client.GameLines
}

func (a *OddsActivities) GetCurrentLines(ctx context.Context, in GetCurrentLinesRequest) (GetCurrentLinesResponse, error) {
	lines, err := a.Client.GetCurrentLines(ctx, in.Season, in.Week)
	if err != nil {
		return GetCurrentLinesResponse{}, fmt.Errorf("getting current lines: %w", err)
	}

	return GetCurrentLinesResponse{
		Lines: lines,
	}, nil
}

type GetLineHistoryRequest struct {
	GameID string
}

type GetLineHistoryResponse struct {
	Lines []client.GameLines
}

func (a *OddsActivities) GetLineHistory(ctx context.Context, in GetLineHistoryRequest) (GetLineHistoryResponse, error) {
	history, err := a.Client.GetLineHistory(ctx, in.GameID)
	if err != nil {
		return GetLineHistoryResponse{}, fmt.Errorf("getting line history: %w", err)
	}

	return GetLineHistoryResponse{
		Lines: history,
	}, nil
}
