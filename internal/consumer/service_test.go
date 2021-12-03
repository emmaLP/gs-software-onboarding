package consumer

import (
	"context"
	"errors"
	"testing"

	"github.com/emmaLP/gs-software-onboarding/internal/database"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews"
	hnModel "github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
	"go.uber.org/zap"
)

func TestProcessStories(t *testing.T) {
	tests := map[string]struct {
		hnMock        *hackernews.Mock
		dbMock        *database.Mock
		ids           []int
		config        *model.Configuration
		expectedMocks func(t *testing.T, hnMock *hackernews.Mock)
	}{
		"One Item": {
			hnMock: &hackernews.Mock{},
			dbMock: &database.Mock{},
			ids:    []int{1},
			config: &model.Configuration{
				Consumer: model.ConsumerConfig{
					BaseUrl:         "test.com",
					CronSchedule:    "",
					NumberOfWorkers: 2,
				}},

			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock) {
				hnMock.On("GetTopStories").Return([]int{1}, nil)
				hnMock.On("GetItem", 1).Return(&hnModel.Item{ID: 1}, nil)
			},
		},
		"Two Items": {
			hnMock: &hackernews.Mock{},
			dbMock: &database.Mock{},
			ids:    []int{1},
			config: &model.Configuration{
				Consumer: model.ConsumerConfig{
					BaseUrl:         "test.com",
					CronSchedule:    "",
					NumberOfWorkers: 2,
				}},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock) {
				hnMock.On("GetTopStories").Return([]int{1, 2}, nil)
				hnMock.On("GetItem", 1).Return(&hnModel.Item{ID: 1}, nil)
				hnMock.On("GetItem", 2).Return(&hnModel.Item{ID: 2}, nil)
			},
		},
		"Error on item1": {
			hnMock: &hackernews.Mock{},
			dbMock: &database.Mock{},
			ids:    []int{1},
			config: &model.Configuration{
				Consumer: model.ConsumerConfig{
					BaseUrl:         "test.com",
					CronSchedule:    "",
					NumberOfWorkers: 2,
				}},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock) {
				hnMock.On("GetTopStories").Return([]int{1, 2}, nil)
				hnMock.On("GetItem", 1).Return(nil, errors.New("Failed to retrieve item"))
				hnMock.On("GetItem", 2).Return(&hnModel.Item{ID: 2}, nil)
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
			service, err := NewService(logger, testConfig.config, context.TODO(), testConfig.hnMock, testConfig.dbMock)
			if err != nil {
				logger.Fatal("An unexpected err happened")
				t.FailNow()
			}
			err = service.processStories()
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
