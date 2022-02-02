package publisher

import (
	"context"
	"fmt"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/emmaLP/gs-software-onboarding/internal/queue"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func ConfigureCron(ctx context.Context, logger *zap.Logger, config *model.Configuration) error {
	c := cron.New()

	var err error
	queueClient, err := queue.New(logger, ctx, &config.RabbitMq)
	if err != nil {
		return fmt.Errorf("Unexpected error when connecting to the database. %w", err)
	}
	defer queueClient.CloseConnection()
	service, err := NewService(logger, config, queueClient)
	if err != nil {
		return fmt.Errorf("An error occurred when trying to instantiate the publisher service: %w", err)
	}
	storyProcessing := func() {
		err = service.processStories()
	}
	storyProcessing()
	if err != nil {
		return fmt.Errorf("Error occurred processing stories: %w", err)
	}

	_, err = c.AddFunc(config.Publisher.CronSchedule, storyProcessing)
	if err != nil {
		return fmt.Errorf("Cron job error, %w", err)
	}

	c.Run()
	return nil
}
