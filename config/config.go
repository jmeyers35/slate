package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

func init() {

}

// Config is the application configuration for slate.
type Config struct {
	TemporalHostPort  string `envconfig:"TEMPORAL_HOSTPORT" default:"localhost:7233"`
	TemporalNamespace string `envconfig:"TEMPORAL_NAMESPACE" default:"default"`
}

func MustLoad() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		panic(fmt.Errorf("loading config: %w", err))
	}
	return &cfg
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	return &cfg, nil
}
