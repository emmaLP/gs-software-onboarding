package consumer

import (
	"context"
	"fmt"

	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews"

	"go.uber.org/zap"
)

func processStories(ctx context.Context, logger *zap.Logger, baseUrl string) error {
	logger.Info("Processing stories")
	client, err := hackernews.New(baseUrl, nil)
	if err != nil {
		fmt.Errorf("Unable to create HackerNew client: %w", err)
	}
	return nil
}
