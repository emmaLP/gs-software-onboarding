package logging

import (
	"fmt"

	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("Unable to create logger: %w", err)
	}

	return logger, nil
}
