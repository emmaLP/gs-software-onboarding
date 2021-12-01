package consumer

import (
	"context"
	"fmt"

	"github.com/emmaLP/gs-software-onboarding/internal/model"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func ConfigureCron(ctx context.Context, logger *zap.Logger, config *model.Configuration) error {
	c := cron.New()

	var err error
	service, err := NewService(logger, &config.Consumer, nil)
	if err != nil {
		return fmt.Errorf("An error occurred when trying to instantiate the consumer service: %w", err)
	}
	storyProcessing := func() {
		err = service.processStories(ctx)
	}
	storyProcessing()
	if err != nil {
		return fmt.Errorf("Error occurred processing stories: %w", err)
	}

	_, err = c.AddFunc(config.Consumer.CronSchedule, storyProcessing)
	if err != nil {
		return fmt.Errorf("Cron job error, %w", err)
	}

	c.Run()
	return nil
}
