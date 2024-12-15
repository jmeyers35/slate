package converters

import (
	espnclient "github.com/jmeyers35/slate/pkg/espn/client"
	"github.com/jmeyers35/slate/pkg/storage"
)

type ESPNAPIConverter struct{}

func (c ESPNAPIConverter) ConvertTeam(team espnclient.Team) storage.Team {
	return storage.Team{
		ESPNID:   team.Team.ID,
		Name:     team.Team.Name,
		TeamCode: team.Team.Abbreviation,
	}
}

func (c ESPNAPIConverter) ConvertAthlete(athlete espnclient.Athlete, teamID int) storage.Player {
	return storage.Player{
		ESPNID:   athlete.ID,
		Name:     athlete.FullName,
		Position: athlete.Position.Name,
		TeamID:   teamID,
	}
}
