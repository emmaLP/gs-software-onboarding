package main

import (
	"log"

	"github.com/emmaLP/gs-software-onboarding/internal/api"
	"github.com/emmaLP/gs-software-onboarding/internal/config"
	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
	"github.com/emmaLP/gs-software-onboarding/internal/logging"
	"go.uber.org/zap"
)

func main() {
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

	grpcClient, err := grpc.NewClient(configuration.Api.GrpcAddress, logger)
	if err != nil {
		logger.Fatal("Unable to create GRPC client.", zap.Error(err))
	}
	defer grpcClient.Close()

	server, err := api.NewServer(logger, grpcClient)
	if err != nil {
		logger.Fatal("Unable to instantiate api server", zap.Error(err))
	}
	server.StartServer(configuration.Api.Address)
}
