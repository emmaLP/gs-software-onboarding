package database

import (
	"context"

	"github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) SaveItem(ctx context.Context, item *model.Item) error {
	args := m.Called(ctx, item)

	return args.Error(0)
}

func (m *Mock) ListAll(ctx context.Context) ([]*model.Item, error) {
	args := m.Called(ctx)
	return find(args)
}

func (m *Mock) ListStories(ctx context.Context) ([]*model.Item, error) {
	args := m.Called(ctx)
	return find(args)
}

func (m *Mock) ListJobs(ctx context.Context) ([]*model.Item, error) {
	args := m.Called(ctx)
	return find(args)
}

func find(args mock.Arguments) ([]*model.Item, error) {
	collection, ok := args.Get(0).([]*model.Item)
	if !ok {
		return nil, args.Error(1)
	}

	return collection, args.Error(1)
}

func (m *Mock) CloseConnection(ctx context.Context) {
	// Do nothing as this is a mock
}
