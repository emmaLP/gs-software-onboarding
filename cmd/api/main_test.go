//go:build integration

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/emmaLP/gs-software-onboarding/internal/caching"
	"github.com/emmaLP/gs-software-onboarding/internal/config"
	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
	"github.com/emmaLP/gs-software-onboarding/internal/logging"
	"github.com/emmaLP/gs-software-onboarding/pkg/test"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	log.Println("Starting API integration tests")
	grpcPort := strconv.Itoa(grpcRandomPort())
	os.Setenv("GRPC_PORT", grpcPort)

	ctx := context.TODO()
	_, dbConfig, err := test.SetupMongoContainer(ctx)
	if err != nil {
		log.Println("FAIL - unable to setup mongo")
		os.Exit(1)
	}

	os.Setenv("DATABASE_NAME", "int-test-api")
	os.Setenv("DATABASE_USERNAME", dbConfig.User)
	os.Setenv("DATABASE_PASSWORD", dbConfig.Password)
	os.Setenv("DATABASE_HOST", dbConfig.Host)
	os.Setenv("DATABASE_PORT", fmt.Sprint(dbConfig.Port))

	redis, err := test.SetupRedis(ctx)
	if err != nil {
		log.Println("FAIL - unable to setup redis")
		os.Exit(1)
	}
	os.Setenv("CACHE_ADDRESS", redis.URI)
	os.Setenv("GRPC_ADDRESS", fmt.Sprintf("localhost:%s", grpcPort))

	os.Setenv("API_ADDRESS", fmt.Sprintf(":%d", apiRandomPort()))

	go func() {
		log.Println("Starting a grpc server in go routine")
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

		server := grpc.NewServer(configuration.Grpc.Port, logger, grpc.NewHandler(cacheClient, databaseClient, logger))
		grpcServer, err := server.Start()
		if err != nil {
			logger.Fatal("Failed to start grpc server:", zap.Error(err))
		}
		defer grpcServer.Stop()
		log.Println("GRPC Server exiting")
	}()
	go func() {
		log.Println("Starting api server")
		main()
		log.Println("API Server exiting")
	}()
	// Adds buffer to allow api and grpc servers to get started
	time.Sleep(10 * time.Second)
	os.Exit(m.Run())
}

func grpcRandomPort() int {
	return randomPort(51000, 55000)
}

func apiRandomPort() int {
	return randomPort(8081, 9000)
}

func randomPort(min, max int) int {
	return rand.Intn(max-min) + min
}
