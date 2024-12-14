package main

import (
	"context"

	"github.com/jmeyers35/slate/pkg/espn/client"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	nflClient := client.NewNFL()
	ll, _ := zap.NewProduction()

	resp, err := nflClient.GetRoster(ctx, "4")
	if err != nil {
		ll.Error("error getting roster", zap.Error(err))
		return
	}

	allAthletes := []client.Athlete{}
	for _, athleteGroup := range resp.Athletes {
		for _, athlete := range athleteGroup.Athletes {
			allAthletes = append(allAthletes, athlete)
			ll.Info("athlete", zap.String("name", athlete.FullName))
		}
	}

	ll.Info("total athletes", zap.Int("count", len(allAthletes)))
}
