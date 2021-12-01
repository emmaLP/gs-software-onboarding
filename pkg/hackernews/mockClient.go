package hackernews

import (
	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) GetTopStories() ([]int, error) {
	args := m.Called()

	idsArg, ok := args.Get(0).([]int)
	if !ok {
		return nil, nil
	}

	return idsArg, args.Error(1)
}

func (m *Mock) GetItem(id int) (*model.Item, error) {
	args := m.Called(id)

	itemArg, ok := args.Get(0).(*model.Item)
	if !ok {
		return nil, nil
	}

	return itemArg, args.Error(1)
}
