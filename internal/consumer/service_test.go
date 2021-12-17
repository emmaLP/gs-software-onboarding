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
		expectedMocks func(t *testing.T, hnMock *hackernews.Mock, dbMock *database.Mock)
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
				},
			},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock, dbMock *database.Mock) {
				hnMock.On("GetTopStories").Return([]int{1}, nil)
				hnMock.On("GetItem", 1).Return(&hnModel.Item{ID: 1}, nil)
				dbMock.On("SaveItem", context.TODO(), &hnModel.Item{ID: 1}).Return(nil)
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
				},
			},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock, dbMock *database.Mock) {
				hnMock.On("GetTopStories").Return([]int{1, 2}, nil)
				hnMock.On("GetItem", 1).Return(&hnModel.Item{ID: 1}, nil)
				hnMock.On("GetItem", 2).Return(&hnModel.Item{ID: 2}, nil)
				dbMock.On("SaveItem", context.TODO(), &hnModel.Item{ID: 1}).Return(nil).Once()
				dbMock.On("SaveItem", context.TODO(), &hnModel.Item{ID: 2}).Return(nil).Once()
			},
		},
		"Unable to get item from hackernews": {
			hnMock: &hackernews.Mock{},
			dbMock: &database.Mock{},
			ids:    []int{1},
			config: &model.Configuration{
				Consumer: model.ConsumerConfig{
					BaseUrl:         "test.com",
					CronSchedule:    "",
					NumberOfWorkers: 2,
				},
			},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock, dbMock *database.Mock) {
				hnMock.On("GetTopStories").Return([]int{1, 2}, nil)
				hnMock.On("GetItem", 1).Return(nil, errors.New("Failed to retrieve item"))
				hnMock.On("GetItem", 2).Return(&hnModel.Item{ID: 2}, nil)
				dbMock.On("SaveItem", context.TODO(), &hnModel.Item{ID: 2}).Return(nil).Once()
			},
		},
		"Unable save item": {
			hnMock: &hackernews.Mock{},
			dbMock: &database.Mock{},
			ids:    []int{1},
			config: &model.Configuration{
				Consumer: model.ConsumerConfig{
					BaseUrl:         "test.com",
					CronSchedule:    "",
					NumberOfWorkers: 2,
				},
			},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock, dbMock *database.Mock) {
				hnMock.On("GetTopStories").Return([]int{2}, nil)
				hnMock.On("GetItem", 2).Return(&hnModel.Item{ID: 2}, nil)
				dbMock.On("SaveItem", context.TODO(), &hnModel.Item{ID: 2}).Return(errors.New("Failed to save item")).Once()
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.hnMock, testConfig.dbMock)
			}
			logger, err := zap.NewProduction()
			if err != nil {
				panic(err)
			}
			service, err := NewService(logger, testConfig.config, testConfig.dbMock, WithHackerNewsClient(testConfig.hnMock))
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
				testConfig.dbMock.AssertExpectations(t)
			}
		})
	}
}
