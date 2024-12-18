package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type postgresStorage struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) Storage {
	return &postgresStorage{db: db}
}

func (s *postgresStorage) GetTeams(ctx context.Context) ([]*Team, error) {
	const query = `
		SELECT team_id, team_code, team_name, espn_id
		FROM teams`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teams []*Team
	for rows.Next() {
		team := &Team{}
		err := rows.Scan(
			&team.ID,
			&team.TeamCode,
			&team.Name,
			&team.ESPNID,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
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

func (s *postgresStorage) UpsertGame(ctx context.Context, game *Game) error {
	const query = `
		INSERT INTO games (
			week, season, home_team_id, away_team_id, game_date, dome
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (week, season, home_team_id, away_team_id)
		DO UPDATE SET
			game_date = EXCLUDED.game_date,
			dome = EXCLUDED.dome
		RETURNING game_id`

	_, err := s.db.ExecContext(ctx, query,
		game.Week,
		game.Season,
		game.HomeTeamID,
		game.AwayTeamID,
		game.GameDate,
		game.Dome,
	)
	return err
}

func (s *postgresStorage) UpsertLine(ctx context.Context, line Line) error {
	const query = `
		INSERT INTO vegas_lines (
			line_id, game_id, home_team_spread, over_under, home_team_total, away_team_total, last_updated
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (game_id)
		DO UPDATE SET
			home_team_spread = EXCLUDED.home_team_spread,
			over_under = EXCLUDED.over_under,
			home_team_total = EXCLUDED.home_team_total,
			away_team_total = EXCLUDED.away_team_total,
			last_updated = EXCLUDED.last_updated`

	if _, err := s.db.ExecContext(ctx, query,
		line.LineID,
		line.GameID,
		line.HomeSpread,
		line.OverUnder,
		line.HomeTeamTotal,
		line.AwayTeamTotal,
		line.LastUpdated,
	); err != nil {
		return fmt.Errorf("upserting game lines: %w", err)
	}

	return nil
}

func (s *postgresStorage) GetGameIDByTeams(ctx context.Context, season, week int, homeTeam, awayTeam string) (string, error) {
	query := `
		SELECT g.game_id
		FROM games g
		JOIN teams ht ON g.home_team_id = ht.team_id
		JOIN teams at ON g.away_team_id = at.team_id
		WHERE g.season = $1
		AND g.week = $2
		AND (
			(ht.team_name = $3 AND at.team_name = $4)
			OR 
			(ht.team_code = $3 AND at.team_code = $4)
		)`

	var gameID string
	err := s.db.QueryRowContext(ctx, query, season, week, homeTeam, awayTeam).Scan(&gameID)
	if err != nil {
		return "", fmt.Errorf("getting game ID: %w", err)
	}

	return gameID, nil
}
