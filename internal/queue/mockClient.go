package queue

import (
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) SendMessage(item commonModel.Item) error {
	args := m.Called(item)
	err := args.Error(0)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mock) ReceiveMessage(msgChan chan commonModel.Item) error {
	args := m.Called(msgChan)

	err := args.Error(0)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mock) CloseConnection() {
	// Do nothing as this is a mock
}
