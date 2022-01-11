package consumer

import (
	"context"
	"errors"
	"testing"

	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
	"github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

func TestProcessMessages(t *testing.T) {
	story := model.Item{ID: 1, Type: "story"}
	job := model.Item{ID: 2, Type: "job"}
	tests := map[string]struct {
		mockGrpc           *grpc.Mock
		expectedMock       func(t *testing.T, mockGrpc *grpc.Mock)
		itemsToSave        []model.Item
		expectedErrMessage string
	}{
		"Successfully save 2 items": {
			mockGrpc: &grpc.Mock{},
			itemsToSave: []model.Item{
				story,
				job,
			},
			expectedMock: func(t *testing.T, mockGrpc *grpc.Mock) {
				mockGrpc.On("SaveItem", context.TODO(), &story).Return(nil).Once()
				mockGrpc.On("SaveItem", context.TODO(), &job).Return(nil).Once()
			},
		},
		"Unable to save item": {
			mockGrpc: &grpc.Mock{},
			itemsToSave: []model.Item{
				story,
				job,
			},
			expectedMock: func(t *testing.T, mockGrpc *grpc.Mock) {
				mockGrpc.On("SaveItem", context.TODO(), &story).Return(nil).Once()
				mockGrpc.On("SaveItem", context.TODO(), &job).Return(errors.New("Fail.")).Once()
			},
			expectedErrMessage: "Failed to save item",
		},
	}

	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger := zaptest.NewLogger(t, zaptest.WrapOptions(zap.Hooks(func(e zapcore.Entry) error {
				if testConfig.expectedErrMessage != "" {
					assert.Equal(t, zap.ErrorLevel, e.Level)
					assert.Equal(t, testConfig.expectedErrMessage, e.Message)
				}
				return nil
			})))
			if testConfig.expectedMock != nil {
				testConfig.expectedMock(t, testConfig.mockGrpc)
			}
			service := New(logger, testConfig.mockGrpc)
			testChan := make(chan model.Item)
			defer close(testChan)

			// execute
			go service.ProcessMessages(context.TODO(), testChan)

			for _, item := range testConfig.itemsToSave {
				testChan <- item
			}

			if testConfig.expectedMock != nil {
				testConfig.mockGrpc.AssertExpectations(t)
			}
		})
	}
}
