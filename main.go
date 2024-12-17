package main

import (
	"fmt"

	"github.com/jmeyers35/slate/config"
	slatedb "github.com/jmeyers35/slate/db"
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
	scraper.InitWorker(scraperWorker, storage)

	err = scraperWorker.Run(worker.InterruptCh())
	if err != nil {
		logger.Error("running scraper worker", zap.Error(err))
		return
	}
}
