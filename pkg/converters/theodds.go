package converters

import (
	"github.com/jmeyers35/slate/pkg/odds/client"
	"github.com/jmeyers35/slate/pkg/storage"
)

type TheOddsConverter struct{}

func (c *TheOddsConverter) ConvertLine(line client.GameLines, gameID int) storage.Line {
	sline := storage.Line{
		GameID:        gameID,
		ProviderID:    line.ProviderID,
		HomeSpread:    line.HomeSpread,
		OverUnder:     line.OverUnder,
		HomeMoneyline: line.HomeMoneyline,
		AwayMoneyline: line.AwayMoneyline,
		LastUpdated:   line.LastUpdated,
	}

	if line.HomeTeamTotal != nil {
		sline.HomeTeamTotal = *line.HomeTeamTotal
	}

	if line.AwayTeamTotal != nil {
		sline.AwayTeamTotal = *line.AwayTeamTotal
	}
	return sline
}
