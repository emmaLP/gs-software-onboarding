package publisher

import (
	"errors"
	"testing"

	"github.com/emmaLP/gs-software-onboarding/internal/queue"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews"
	"go.uber.org/zap"
)

func TestProcessStories(t *testing.T) {
	tests := map[string]struct {
		hnMock        *hackernews.Mock
		queueMock     *queue.Mock
		ids           []int
		config        *model.Configuration
		expectedMocks func(t *testing.T, hnMock *hackernews.Mock, queueMock *queue.Mock)
	}{
		"One Item": {
			hnMock:    &hackernews.Mock{},
			queueMock: &queue.Mock{},
			ids:       []int{1},
			config: &model.Configuration{
				Publisher: model.PublisherConfig{
					BaseUrl:      "test.com",
					CronSchedule: "",
				},
			},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock, queueMock *queue.Mock) {
				hnMock.On("GetTopStories").Return([]int{1}, nil)
				hnMock.On("GetItem", 1).Return(&commonModel.Item{ID: 1}, nil)
				queueMock.On("SendMessage", commonModel.Item{ID: 1}).Return(nil)
			},
		},
		"Two Items": {
			hnMock:    &hackernews.Mock{},
			queueMock: &queue.Mock{},
			ids:       []int{1},
			config: &model.Configuration{
				Publisher: model.PublisherConfig{
					BaseUrl:      "test.com",
					CronSchedule: "",
				},
			},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock, queueMock *queue.Mock) {
				hnMock.On("GetTopStories").Return([]int{1, 2}, nil)
				hnMock.On("GetItem", 1).Return(&commonModel.Item{ID: 1}, nil)
				hnMock.On("GetItem", 2).Return(&commonModel.Item{ID: 2}, nil)
				queueMock.On("SendMessage", commonModel.Item{ID: 1}).Return(nil).Once()
				queueMock.On("SendMessage", commonModel.Item{ID: 2}).Return(nil).Once()
			},
		},
		"Unable to get item from hackernews": {
			hnMock:    &hackernews.Mock{},
			queueMock: &queue.Mock{},
			ids:       []int{1},
			config: &model.Configuration{
				Publisher: model.PublisherConfig{
					BaseUrl:      "test.com",
					CronSchedule: "",
				},
			},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock, queueMock *queue.Mock) {
				hnMock.On("GetTopStories").Return([]int{1, 2}, nil)
				hnMock.On("GetItem", 1).Return(nil, errors.New("Failed to retrieve item"))
				hnMock.On("GetItem", 2).Return(&commonModel.Item{ID: 2}, nil)
				queueMock.On("SendMessage", commonModel.Item{ID: 2}).Return(nil).Once()
			},
		},
		"Unable send item": {
			hnMock:    &hackernews.Mock{},
			queueMock: &queue.Mock{},
			ids:       []int{1},
			config: &model.Configuration{
				Publisher: model.PublisherConfig{
					BaseUrl:      "test.com",
					CronSchedule: "",
				},
			},
			expectedMocks: func(t *testing.T, hnMock *hackernews.Mock, queueMock *queue.Mock) {
				hnMock.On("GetTopStories").Return([]int{2}, nil)
				hnMock.On("GetItem", 2).Return(&commonModel.Item{ID: 2}, nil)
				queueMock.On("SendMessage", commonModel.Item{ID: 2}).Return(errors.New("Failed to send item")).Once()
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.hnMock, testConfig.queueMock)
			}
			logger, err := zap.NewProduction()
			if err != nil {
				panic(err)
			}
			service, err := NewService(logger, testConfig.config, testConfig.queueMock, WithHackerNewsClient(testConfig.hnMock))
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
				testConfig.queueMock.AssertExpectations(t)
			}
		})
	}
}
