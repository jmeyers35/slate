package scraper

import (
	espnclient "github.com/jmeyers35/slate/pkg/espn/client"
	"github.com/jmeyers35/slate/pkg/storage"
	"go.temporal.io/sdk/worker"
)

func InitWorker(w worker.Worker, storage storage.Storage) {
	w.RegisterWorkflow(ScrapeNFLTeam)
	w.RegisterWorkflow(BackfillTeams)

	nflClient := espnclient.NewNFL()
	espnActivities := &ESPNActivities{
		client: nflClient,
	}
	w.RegisterActivity(espnActivities.GetPlayersForTeam)
	w.RegisterActivity(espnActivities.GetTeam)
	w.RegisterActivity(espnActivities.GetTeams)

	storageActivities := &StorageActivities{
		Storage: storage,
	}
	w.RegisterActivity(storageActivities.UpsertTeam)
	w.RegisterActivity(storageActivities.GetTeamByESPNID)
	w.RegisterActivity(storageActivities.UpsertPlayer)
}
