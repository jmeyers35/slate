package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmeyers35/slate/config"
)

func InitDB(config *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", config.PostgresConnString())
	if err != nil {
		return nil, fmt.Errorf("opening postgres connection: %w", err)
	}

	// TODO: tune/make configurable?
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	return db, nil
}
