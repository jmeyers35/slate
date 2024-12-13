-- Seed Teams
INSERT INTO teams (team_code, team_name) VALUES
    ('KC', 'Kansas City Chiefs'),
    ('SF', 'San Francisco 49ers'),
    ('BUF', 'Buffalo Bills'),
    ('PHI', 'Philadelphia Eagles'),
    ('DAL', 'Dallas Cowboys'),
    ('BAL', 'Baltimore Ravens');

-- Seed Players
INSERT INTO players (name, position, team_id) VALUES
    ('Patrick Mahomes', 'QB', (SELECT team_id FROM teams WHERE team_code = 'KC')),
    ('Travis Kelce', 'TE', (SELECT team_id FROM teams WHERE team_code = 'KC')),
    ('Christian McCaffrey', 'RB', (SELECT team_id FROM teams WHERE team_code = 'SF')),
    ('Deebo Samuel', 'WR', (SELECT team_id FROM teams WHERE team_code = 'SF')),
    ('Josh Allen', 'QB', (SELECT team_id FROM teams WHERE team_code = 'BUF')),
    ('Khalil Shakir', 'WR', (SELECT team_id FROM teams WHERE team_code = 'BUF')),
    ('Jalen Hurts', 'QB', (SELECT team_id FROM teams WHERE team_code = 'PHI')),
    ('AJ Brown', 'WR', (SELECT team_id FROM teams WHERE team_code = 'PHI')),
    ('Dak Prescott', 'QB', (SELECT team_id FROM teams WHERE team_code = 'DAL')),
    ('CeeDee Lamb', 'WR', (SELECT team_id FROM teams WHERE team_code = 'DAL'));

-- Seed Games (Week 1 matchups)
INSERT INTO games (week, season, home_team_id, away_team_id, game_date, dome, temperature, wind_speed) VALUES
    (1, 2024, 
     (SELECT team_id FROM teams WHERE team_code = 'KC'),
     (SELECT team_id FROM teams WHERE team_code = 'SF'),
     '2024-09-08 13:00:00', false, 75, 8),
    (1, 2024,
     (SELECT team_id FROM teams WHERE team_code = 'BUF'),
     (SELECT team_id FROM teams WHERE team_code = 'DAL'),
     '2024-09-08 16:25:00', false, 72, 12);

-- Seed Vegas Lines
INSERT INTO vegas_lines (game_id, home_team_spread, over_under, home_team_total, away_team_total, last_updated) VALUES
    ((SELECT game_id FROM games WHERE home_team_id = (SELECT team_id FROM teams WHERE team_code = 'KC') AND week = 1 AND season = 2024),
     -3.5, 51.5, 27.5, 24.0, '2024-09-07 12:00:00'),
    ((SELECT game_id FROM games WHERE home_team_id = (SELECT team_id FROM teams WHERE team_code = 'BUF') AND week = 1 AND season = 2024),
     -2.5, 49.5, 26.0, 23.5, '2024-09-07 12:00:00');

-- Seed DraftKings Contests
INSERT INTO draftkings_contests (contest_id, week, season, contest_type, entry_fee, total_entries, max_entries, total_prizes) VALUES
    ('MIL1', 1, 2024, 'GPP', 20.00, 100000, 150, 1000000.00),
    ('DU1', 1, 2024, 'Double Up', 10.00, 1000, 1, 1800.00);

-- Seed DraftKings Players (example salaries and projected ownership)
INSERT INTO draftkings_players (player_id, game_id, week, season, salary, ownership_percentage) VALUES
    ((SELECT player_id FROM players WHERE name = 'Patrick Mahomes'),
     (SELECT game_id FROM games WHERE home_team_id = (SELECT team_id FROM teams WHERE team_code = 'KC') AND week = 1),
     1, 2024, 8200, 15.5),
    ((SELECT player_id FROM players WHERE name = 'Christian McCaffrey'),
     (SELECT game_id FROM games WHERE home_team_id = (SELECT team_id FROM teams WHERE team_code = 'KC') AND week = 1),
     1, 2024, 9000, 22.3),
    ((SELECT player_id FROM players WHERE name = 'Josh Allen'),
     (SELECT game_id FROM games WHERE home_team_id = (SELECT team_id FROM teams WHERE team_code = 'BUF') AND week = 1),
     1, 2024, 7800, 18.2);

