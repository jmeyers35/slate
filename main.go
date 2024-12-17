package main

import (
	"fmt"
	"time"

	"github.com/jmeyers35/slate/config"
	slatedb "github.com/jmeyers35/slate/db"
	oddsclient "github.com/jmeyers35/slate/pkg/odds/client"
	"github.com/jmeyers35/slate/pkg/scraper"
	"github.com/jmeyers35/slate/pkg/storage"
	_ "github.com/lib/pq"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

func main() {
	appConfig := config.MustLoad()

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Errorf("creating logger: %w", err))
	}

	logger.Debug("starting slate", zap.Any("config", appConfig))

	c, err := client.Dial(client.Options{
		HostPort:  appConfig.TemporalHostPort,
		Namespace: appConfig.TemporalNamespace,
	})
	if err != nil {
		logger.Error("creating temporal client", zap.Error(err))
		return
	}
	defer c.Close()

	scraperWorker := worker.New(c, scraper.DefaultTaskQueueName, worker.Options{})
	db, err := slatedb.InitDB(appConfig)
	if err != nil {
		logger.Error("initializing database", zap.Error(err))
		return
	}
	storage := storage.NewPostgres(db)

	oddsClient, err := oddsclient.New(oddsclient.Config{
		Provider:  oddsclient.TheOdds,
		APIKey:    appConfig.TheOddsAPIKey,
		RateLimit: time.Second,
	})
	if err != nil {
		logger.Error("creating odds client", zap.Error(err))
		return
	}

	scraper.InitWorker(scraperWorker, storage, oddsClient)

	err = scraperWorker.Run(worker.InterruptCh())
	if err != nil {
		logger.Error("running scraper worker", zap.Error(err))
		return
	}
}
