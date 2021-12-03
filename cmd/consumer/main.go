package main

import (
	"context"
	"fmt"
	"log"

	"github.com/emmaLP/gs-software-onboarding/internal/config"
	"github.com/emmaLP/gs-software-onboarding/internal/consumer"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger, err := configureLogger()
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

	if err := consumer.ConfigureCron(ctx, logger, configuration); err != nil {
		logger.Fatal("Failed to configure the cron", zap.Error(err))
	}
}

func configureLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("Unable to create logger: %w", err)
	}

	return logger, nil
}
