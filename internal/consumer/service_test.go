package consumer

import (
	"context"
	"testing"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews"
	hnModel "github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
	"go.uber.org/zap"
)

func TestProcessStories(t *testing.T) {
	tests := map[string]struct {
		hnMock         *hackernews.Mock
		ids            []int
		consumerConfig *model.ConsumerConfig
		expectedMocks  func(t *testing.T, hnMock *hackernews.Mock)
	}{
		"One Item": {
			hnMock: &hackernews.Mock{},
			ids:    []int{1},
			consumerConfig: &model.ConsumerConfig{
				BaseUrl:         "test.com",
				CronSchedule:    "",
				NumberOfWorkers: 2,
			},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock) {
				hnMock.On("GetTopStories").Return([]int{1}, nil)
				hnMock.On("GetItem", 1).Return(&hnModel.Item{ID: 1}, nil)
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.hnMock)
			}
			logger, err := zap.NewProduction()
			if err != nil {
				panic(err)
			}
			service, err := NewService(logger, testConfig.consumerConfig, testConfig.hnMock)
			if err != nil {
				logger.Fatal("An unexpected err happened")
				t.FailNow()
			}
			err = service.processStories(context.TODO())
			if err != nil {
				logger.Fatal("An unexpected err happened")
				t.FailNow()
			}

			if testConfig.expectedMocks != nil {
				testConfig.hnMock.AssertExpectations(t)
			}
		})
	}
}
