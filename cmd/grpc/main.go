package main

import (
	"context"
	"log"

	"github.com/emmaLP/gs-software-onboarding/internal/caching"
	"github.com/emmaLP/gs-software-onboarding/internal/config"
	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
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

	databaseClient, err := database.New(ctx, logger, &configuration.Database)
	if err != nil {
		logger.Fatal("Unexpected error when connecting to the database.", zap.Error(err))
	}
	defer databaseClient.CloseConnection(ctx)

	cacheClient, err := caching.New(ctx, configuration.Cache.Address, databaseClient, logger)
	if err != nil {
		logger.Fatal("Unexpected error when connecting to the cache.", zap.Error(err))
	}
	defer cacheClient.Close()
	grpcHandler := grpc.Handler{
		ItemCache: cacheClient,
	}
	server := grpc.NewServer(configuration.Grpc.Port, logger, grpcHandler)
	logger.Debug("Starting grpc server")
	grpcServer, err := server.Start()
	if err != nil {
		logger.Fatal("Failed to start grpc server:", zap.Error(err))
	}
	defer grpcServer.Stop()
}
