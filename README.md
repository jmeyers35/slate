# slate

A DraftKings fantasy football data pipeline and modeling system.

## Deployment

Deployed bits are hosted on Fly.io, including the DB.

## Data Dependencies

The system has several data dependencies that must be satisfied in a specific order for workflows to function correctly.

### Bootstrap Requirements

1. **Teams** - Required before running any game-related workflows
   - Teams must be bootstrapped from ESPN before processing games
   - Teams map ESPN IDs to internal IDs used throughout the system
   - Run the backfill_teams workflow to populate teams

2. **Games** - Depends on Teams
   - Requires team data to be present
   - Uses ESPN IDs to map to internal team IDs
   - Stores unique games by (week, season, home_team_id, away_team_id)

3. **Players** - Depends on Teams
   - Players are associated with Teams via team_id
   - Player data includes position and basic info from ESPN

4. **Player Stats** - Depends on Players and Games
   - Links players and games through their respective IDs
   - Tracks performance metrics and DraftKings points

5. **Vegas Lines** - Depends on Games
   - References games via game_id
   - Tracks spreads and over/unders

6. **Injury Reports** - Depends on Players
   - Links to players via player_id
   - Includes practice and game status

### Recommended Bootstrap Order
1. Run backfill_teams workflow to populate teams
2. Run scrape_schedule workflow to fetch and store games
3. Run player_coordinator workflow to populate player data
4. Other workflows (stats, lines, injuries) can run once their dependencies are met