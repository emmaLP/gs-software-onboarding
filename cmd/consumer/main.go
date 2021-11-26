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
	defer func() {
		err = logger.Sync()
	}()
	if err != nil {
		return nil, fmt.Errorf("Failed to perform sync: %w", err)
	}

	return logger, nil
}
