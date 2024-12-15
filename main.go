package main

import (
	"fmt"

	"github.com/jmeyers35/slate/config"
	slatetemporal "github.com/jmeyers35/slate/pkg/temporal"
	"github.com/jmeyers35/slate/pkg/temporal/scraper"
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

	scraperWorker := worker.New(c, slatetemporal.DefaultTaskQueueName, worker.Options{})
	scraper.InitWorker(scraperWorker)

	err = scraperWorker.Run(worker.InterruptCh())
	if err != nil {
		logger.Error("running scraper worker", zap.Error(err))
		return
	}
}
