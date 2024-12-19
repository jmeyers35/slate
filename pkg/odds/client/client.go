package client

import (
	"context"
	"fmt"
)

// Client defines the interface for fetching odds data
type Client interface {
	// GetCurrentLines fetches the latest lines for all games in the given week
	GetCurrentLines(ctx context.Context, season int, week int) ([]GameLines, error)

	// GetLineHistory fetches historical line movements for a specific game.
	// Some implementations may not support this.
	GetLineHistory(ctx context.Context, gameID string) ([]GameLines, error)
}

// New creates a new odds client for the specified provider
func New(config Config) (Client, error) {
	switch config.Provider {
	case TheOdds:
		return newTheOddsClient(config)
	// Add other providers as they're implemented
	default:
		return nil, fmt.Errorf("unknown odds provider: %s", config.Provider)
	}
}
