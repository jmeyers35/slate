package converters

import (
	espnclient "github.com/jmeyers35/slate/pkg/espn/client"
	"github.com/jmeyers35/slate/pkg/storage"
)

type ESPNAPIConverter struct{}

func (c ESPNAPIConverter) ConvertTeam(team espnclient.Team) storage.Team {
	return storage.Team{
		ESPNID:   team.Team.ID,
		Name:     team.Team.DisplayName,
		TeamCode: team.Team.Abbreviation,
	}
}

func (c ESPNAPIConverter) ConvertAthlete(athlete espnclient.Athlete, teamID string) storage.Player {
	return storage.Player{
		ESPNID:   athlete.ID,
		Name:     athlete.FullName,
		Position: athlete.Position.Name,
		TeamID:   teamID,
	}
}

func (c ESPNAPIConverter) ConvertGame(event espnclient.Event, week, season int) storage.Game {
	var homeTeamID, awayTeamID string
	var dome bool

	for _, competitor := range event.Competitions[0].Competitors {
		if competitor.HomeAway == "home" {
			homeTeamID = competitor.Team.ID
		} else {
			awayTeamID = competitor.Team.ID
		}
	}

	if event.Venue.Indoor {
		dome = true
	}

	return storage.Game{
		Week:       week,
		Season:     season,
		GameDate:   event.Date.Time,
		Dome:       dome,
		HomeTeamID: homeTeamID,
		AwayTeamID: awayTeamID,
	}
}