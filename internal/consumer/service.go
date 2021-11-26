package consumer

import (
	"context"

	"go.uber.org/zap"
)

func processStories(ctx context.Context, logger *zap.Logger) {
	logger.Info("Processing stories")
}
