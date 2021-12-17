package main

import (
	"context"
	"log"

	"github.com/emmaLP/gs-software-onboarding/internal/api"
	"github.com/emmaLP/gs-software-onboarding/internal/config"
	"github.com/emmaLP/gs-software-onboarding/internal/logging"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	server, err := api.NewServer(ctx, logger, configuration)
	if err != nil {
		logger.Fatal("Unable to instantiate api server", zap.Error(err))
	}
	server.StartServer()
}
