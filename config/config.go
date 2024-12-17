package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config is the application configuration for slate.
type Config struct {
	PostgresUser             string `envconfig:"POSTGRES_USER" default:"postgres"`
	PostgresPassword         string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	PostgresHost             string `envconfig:"POSTGRES_HOST" default:"localhost"`
	PostgresPort             int    `envconfig:"POSTGRES_PORT" default:"5432"`
	PostgresDB               string `envconfig:"POSTGRES_DB" default:"slate"`
	PostgresSSLMode          string `envconfig:"POSTGRES_SSLMODE" default:"disable"`
	PostgresConnectionString string `envconfig:"POSTGRES_CONN_STRING"`

	TemporalHostPort  string `envconfig:"TEMPORAL_HOSTPORT" default:"localhost:7233"`
	TemporalNamespace string `envconfig:"TEMPORAL_NAMESPACE" default:"default"`

	TheOddsAPIKey string `envconfig:"THEODDS_APIKEY"`
}

func (c Config) PostgresConnString() string {
	if c.PostgresConnectionString != "" {
		return c.PostgresConnectionString
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
		c.PostgresSSLMode,
	)
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
