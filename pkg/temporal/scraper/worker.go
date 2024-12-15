package scraper

import (
	espnclient "github.com/jmeyers35/slate/pkg/espn/client"
	"go.temporal.io/sdk/worker"
)

func InitWorker(w worker.Worker) {
	w.RegisterWorkflow(ScrapeNFLTeam)

	nflClient := espnclient.NewNFL()
	activities := &ESPNActivities{
		client: nflClient,
	}
	w.RegisterActivity(activities.GetPlayersForTeam)
}
