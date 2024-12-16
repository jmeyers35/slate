package scraper

import (
	"fmt"
	"time"

	"github.com/jmeyers35/slate/pkg/converters"
	"go.temporal.io/sdk/workflow"
)

type ScrapePlayerRequest struct {
	PlayerESPNID string
	TeamID       string
}

func ScrapePlayer(ctx workflow.Context, request ScrapePlayerRequest) error {
	var espnActivities *ESPNActivities
	logger := workflow.GetLogger(ctx)

	actx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: 30 * time.Second,
		StartToCloseTimeout:    5 * time.Second,
	})

	var getPlayerResponse GetPlayerResponse
	if err := workflow.ExecuteActivity(actx, espnActivities.GetPlayer, GetPlayerRequest{
		PlayerID: request.PlayerESPNID,
	}).Get(ctx, &getPlayerResponse); err != nil {
		return fmt.Errorf("getting player: %w", err)
	}
	logger.Info("Got player", "player", getPlayerResponse.Player)

	converter := converters.ESPNAPIConverter{}
	converted := converter.ConvertAthlete(getPlayerResponse.Player, request.TeamID)
	logger.Info("Converted player", "player", converted)

	var storageActivities *StorageActivities
	if err := workflow.ExecuteActivity(actx, storageActivities.UpsertPlayer, UpsertPlayerRequest{
		Player: converted,
	}).Get(ctx, nil); err != nil {
		return fmt.Errorf("upserting athlete: %w", err)
	}
	return nil
}
