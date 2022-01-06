package grpc

import (
	"context"

	"github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) ListAll(ctx context.Context) ([]*model.Item, error) {
	return handleCall(m.Called(ctx))
}

func (m *Mock) ListStories(ctx context.Context) ([]*model.Item, error) {
	return handleCall(m.Called(ctx))
}

func (m *Mock) ListJobs(ctx context.Context) ([]*model.Item, error) {
	return handleCall(m.Called(ctx))
}

func handleCall(args mock.Arguments) ([]*model.Item, error) {
	itemsArgs, ok := args.Get(0).([]*model.Item)
	if !ok {
		return nil, args.Error(1)
	}
	return itemsArgs, args.Error(1)
}
