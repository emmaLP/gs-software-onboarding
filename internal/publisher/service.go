package publisher

import (
	"fmt"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/emmaLP/gs-software-onboarding/internal/queue"
	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews"
	"go.uber.org/zap"
)

type service struct {
	logger      *zap.Logger
	hnClient    hackernews.Client
	queueClient queue.Client
}

type Client interface {
	processStories()
}

type ServiceOptions func(*service)

func NewService(logger *zap.Logger, config *model.Configuration, queueClient queue.Client, opts ...ServiceOptions) (*service, error) {
	hnClient, err := hackernews.New(config.Publisher.BaseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("Unable to create HackerNew client: %w", err)
	}

	service := &service{
		logger:      logger,
		hnClient:    hnClient,
		queueClient: queueClient,
	}

	for _, opt := range opts {
		opt(service)
	}

	return service, nil
}

func WithHackerNewsClient(client hackernews.Client) ServiceOptions {
	return func(s *service) {
		s.hnClient = client
	}
}

func (s *service) processStories() error {
	s.logger.Info("Processing stories")

	storyIds, err := s.hnClient.GetTopStories()
	if err != nil {
		return fmt.Errorf("Unable to retrieve the top stories. %w", err)
	}
	for _, id := range storyIds {
		s.publishItem(id)
	}
	s.logger.Info("Finished processing stories")
	return nil
}

func (s *service) publishItem(storyId int) {
	item, err := s.hnClient.GetItem(storyId)
	if err != nil {
		s.logger.Error("An error occurred when trying to fetch the item.", zap.Error(err))
		return
	}

	if !item.Deleted && !item.Dead {
		err := s.queueClient.SendMessage(*item)
		if err != nil {
			s.logger.Error("Failed to send item to queue", zap.Error(err))
		}
	}
}
