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

type GetGameIDRequest struct {
	InternalGameID string
}

type GetGameIDResponse struct {
	ProviderGameID string
}

func (a *OddsActivities) GetGameID(ctx context.Context, in GetGameIDRequest) (GetGameIDResponse, error) {
	providerID, err := a.Client.GetGameID(ctx, in.InternalGameID)
	if err != nil {
		return GetGameIDResponse{}, fmt.Errorf("getting provider game ID: %w", err)
	}

	return GetGameIDResponse{
		ProviderGameID: providerID,
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
