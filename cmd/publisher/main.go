package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/emmaLP/gs-software-onboarding/internal/config"
	"github.com/emmaLP/gs-software-onboarding/internal/logging"
	"github.com/emmaLP/gs-software-onboarding/internal/publisher"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// handle interrupts and propagate the changes across the publisher pipeline
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cancel()
	}()

	logger, err := logging.New()
	if err != nil {
		log.Fatal("Failed to configure the logger", err)
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Fatal("Failed to perform log sync")
		}
	}(logger)

	configuration, err := config.LoadConfig(".")
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	if err := publisher.ConfigureCron(ctx, logger, configuration); err != nil {
		logger.Fatal("Failed to configure the cron", zap.Error(err))
	}
}
