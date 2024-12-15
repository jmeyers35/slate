package storage

import (
	"context"
	"database/sql"
)

type postgresStorage struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) Storage {
	return &postgresStorage{db: db}
}

func (s *postgresStorage) GetTeamByESPNID(ctx context.Context, espnID string) (*Team, error) {
	const query = `
		SELECT team_id, team_code, team_name, espn_id 
		FROM teams
		WHERE espn_id = $1`

	team := &Team{}
	err := s.db.QueryRowContext(ctx, query, espnID).Scan(
		&team.ID,
		&team.TeamCode,
		&team.Name,
		&team.ESPNID,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (s *postgresStorage) UpsertTeam(ctx context.Context, team *Team) error {
	const query = `
        INSERT INTO teams (team_code, team_name, espn_id)
        VALUES ($1, $2, $3)
        ON CONFLICT (espn_id) 
        DO UPDATE SET
            team_code = EXCLUDED.team_code,
            team_name = EXCLUDED.team_name
        RETURNING team_id`

	_, err := s.db.ExecContext(ctx, query,
		team.TeamCode,
		team.Name,
		team.ESPNID,
	)
	return err
}

func (s *postgresStorage) UpsertPlayer(ctx context.Context, player *Player) error {
	const query = `
        INSERT INTO players (name, position, team_id, espn_id)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (espn_id) 
        DO UPDATE SET
            name = EXCLUDED.name,
            position = EXCLUDED.position,
            team_id = EXCLUDED.team_id`
	_, err := s.db.ExecContext(ctx, query,
		player.Name,
		player.Position,
		player.TeamID,
		player.ESPNID)
	return err
}
