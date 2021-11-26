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

	storyProcessing := func() {
		processStories(ctx, logger)
	}
	storyProcessing()

	_, err := c.AddFunc(config.Consumer.CronSchedule, storyProcessing)
	if err != nil {
		return fmt.Errorf("Cron job error, %w", err)
	}

	c.Run()
	return nil
}
