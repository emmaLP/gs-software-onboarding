package consumer

import (
	"context"

	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"go.uber.org/zap"
)

type service struct {
	logger     *zap.Logger
	grpcClient grpc.Client
}

type Service interface {
	ProcessMessages(ctx context.Context, itemChan <-chan commonModel.Item)
}

func New(logger *zap.Logger, grpcClient grpc.Client) *service {
	return &service{
		logger:     logger,
		grpcClient: grpcClient,
	}
}

func (s *service) ProcessMessages(ctx context.Context, itemChan <-chan commonModel.Item) {
	for item := range itemChan {
		err := s.grpcClient.SaveItem(ctx, &item)
		if err != nil {
			s.logger.Error("Failed to save item", zap.Int("id", item.ID), zap.Error(err))
		}
	}
}
