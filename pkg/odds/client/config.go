package client

import "time"

// Provider identifies which odds provider is being used
type Provider string

const (
	// TheOdds represents the-odds-api.com provider
	TheOdds Provider = "theodds"
)

// Config contains the configuration for an odds data provider
type Config struct {
	// Provider identifies which odds provider to use
	Provider Provider
	// APIKey is the authentication key for the provider
	APIKey string
	// RateLimit is the minimum time between API calls
	RateLimit time.Duration
}