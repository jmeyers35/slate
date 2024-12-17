package scraper

import (
	espnactivities "github.com/jmeyers35/slate/pkg/espn/activities"
	espnclient "github.com/jmeyers35/slate/pkg/espn/client"
	"github.com/jmeyers35/slate/pkg/storage"
	storageactivities "github.com/jmeyers35/slate/pkg/storage/activities"
	"go.temporal.io/sdk/worker"
)

func InitWorker(w worker.Worker, storage storage.Storage) {
	w.RegisterWorkflow(ScrapeNFLTeam)
	w.RegisterWorkflow(BackfillTeams)
	w.RegisterWorkflow(PlayerCoordinator)
	w.RegisterWorkflow(ScrapePlayersForTeam)
	w.RegisterWorkflow(ScrapePlayer)
	w.RegisterWorkflow(ScrapeSchedule)
	w.RegisterWorkflow(BackfillSchedule)

	nflClient := espnclient.NewNFL()
	espnActivities := &espnactivities.ESPNActivities{
		Client: nflClient,
	}
	w.RegisterActivity(espnActivities.GetPlayersForTeam)
	w.RegisterActivity(espnActivities.GetTeam)
	w.RegisterActivity(espnActivities.GetTeamsFromESPN)
	w.RegisterActivity(espnActivities.GetPlayer)
	w.RegisterActivity(espnActivities.GetSchedule)

	storageActivities := &storageactivities.StorageActivities{
		Storage: storage,
	}
	w.RegisterActivity(storageActivities.UpsertTeam)
	w.RegisterActivity(storageActivities.GetTeamByESPNID)
	w.RegisterActivity(storageActivities.UpsertPlayer)
	w.RegisterActivity(storageActivities.GetTeamsFromStorage)
	w.RegisterActivity(storageActivities.StoreGames)
}
