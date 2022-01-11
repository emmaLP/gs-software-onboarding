package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/emmaLP/gs-software-onboarding/internal/config"
	"github.com/emmaLP/gs-software-onboarding/internal/consumer"
	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
	"github.com/emmaLP/gs-software-onboarding/internal/logging"
	"github.com/emmaLP/gs-software-onboarding/internal/queue"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// handle interrupts and propagate the changes across the consumer pipeline
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
	qClient, err := queue.New(logger, ctx, &configuration.RabbitMq)
	if err != nil {
		logger.Fatal("Failed to instantiate queue client", zap.Error(err))
	}
	defer qClient.CloseConnection()

	grpcClient, err := grpc.NewClient(configuration.GrpcClient.GrpcAddress, logger)
	if err != nil {
		logger.Fatal("Unable to create GRPC client.", zap.Error(err))
	}
	defer grpcClient.Close()
	logger.Info("GRPC client connected to server")
	var wg sync.WaitGroup
	itemChan := make(chan *commonModel.Item)

	consumerClient := consumer.New(logger, grpcClient)
	for i := 0; i < configuration.Consumer.NumberOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			consumerClient.ProcessMessages(ctx, itemChan)
		}()
	}

	err = qClient.ReceiveMessage(itemChan)
	if err != nil {
		logger.Fatal("Unable to consumer messages from rabbitmq.", zap.Error(err))
	}
	close(itemChan)
	wg.Wait()
}
