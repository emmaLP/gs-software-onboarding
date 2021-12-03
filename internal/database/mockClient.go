package database

import (
	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) SaveItem(item model.Item) error {
	args := m.Called(item)

	return args.Error(1)
}

func (m *Mock) CloseConnection() {
	// Do nothing as this is a mock
}
