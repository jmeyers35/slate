CREATE TABLE teams (
    team_id SERIAL PRIMARY KEY,
    team_code VARCHAR(5) UNIQUE NOT NULL,
    team_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE players (
    player_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    position VARCHAR(5) NOT NULL,
    team_id INTEGER REFERENCES teams(team_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE games (
    game_id SERIAL PRIMARY KEY,
    week INTEGER NOT NULL,
    season INTEGER NOT NULL,
    home_team_id INTEGER REFERENCES teams(team_id),
    away_team_id INTEGER REFERENCES teams(team_id),
    game_date TIMESTAMP NOT NULL,
    dome BOOLEAN,
    temperature INTEGER,
    wind_speed INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE vegas_lines (
    line_id SERIAL PRIMARY KEY,
    game_id INTEGER REFERENCES games(game_id),
    home_team_spread DECIMAL(4,1),
    over_under DECIMAL(4,1),
    home_team_total DECIMAL(4,1),
    away_team_total DECIMAL(4,1),
    last_updated TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE draftkings_contests (
    contest_id VARCHAR(50) PRIMARY KEY,
    week INTEGER NOT NULL,
    season INTEGER NOT NULL,
    contest_type VARCHAR(50) NOT NULL,
    entry_fee DECIMAL(10,2),
    total_entries INTEGER,
    max_entries INTEGER,
    total_prizes DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE draftkings_players (
    dk_player_id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(player_id),
    game_id INTEGER REFERENCES games(game_id),
    week INTEGER NOT NULL,
    season INTEGER NOT NULL,
    salary INTEGER NOT NULL,
    ownership_percentage DECIMAL(5,2),
    actual_points DECIMAL(6,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE player_stats (
    stat_id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(player_id),
    game_id INTEGER REFERENCES games(game_id),
    team_id INTEGER REFERENCES teams(team_id),
    opponent_team_id INTEGER REFERENCES teams(team_id),
    is_home BOOLEAN NOT NULL,
    passing_attempts INTEGER,
    passing_completions INTEGER,
    passing_yards INTEGER,
    passing_touchdowns INTEGER,
    interceptions INTEGER,
    rushing_attempts INTEGER,
    rushing_yards INTEGER,
    rushing_touchdowns INTEGER,
    targets INTEGER,
    receptions INTEGER,
    receiving_yards INTEGER,
    receiving_touchdowns INTEGER,
    fumbles_lost INTEGER,
    two_point_conversions INTEGER,
    snap_count INTEGER,
    snap_percentage DECIMAL(5,2),
    draftkings_points DECIMAL(6,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE injury_reports (
    injury_id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(player_id),
    week INTEGER NOT NULL,
    season INTEGER NOT NULL,
    practice_status VARCHAR(50),
    game_status VARCHAR(50),
    report_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_player_stats_player_game ON player_stats(player_id, game_id);
CREATE INDEX idx_player_stats_game ON player_stats(game_id);
CREATE INDEX idx_draftkings_players_week_season ON draftkings_players(week, season);
CREATE INDEX idx_games_week_season ON games(week, season);
CREATE INDEX idx_players_team ON players(team_id);
CREATE INDEX idx_games_teams ON games(home_team_id, away_team_id);