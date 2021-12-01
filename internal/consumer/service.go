package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	logger          *zap.Logger
	numberOfWorkers int
	hnClient        hackernews.Client
}

type Client interface {
	processStories(ctx context.Context)
}

func NewService(logger *zap.Logger, config *model.ConsumerConfig, hnClient hackernews.Client) (*Service, error) {
	if hnClient == nil {
		var err error
		hnClient, err = hackernews.New(config.BaseUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("Unable to create HackerNew client: %w", err)
		}
	}
	return &Service{
		logger:          logger,
		numberOfWorkers: config.NumberOfWorkers,
		hnClient:        hnClient,
	}, nil
}

func (s *Service) processStories(ctx context.Context) error {
	s.logger.Info("Processing stories")

	errGroup := new(errgroup.Group)
	topStoriesChan := make(chan int)

	for i := 0; i < s.numberOfWorkers; i++ {
		errGroup.Go(func() error {
			return s.saveItem(topStoriesChan)
		})
	}

	s.logger.Info("Finished getting the items")

	storyIds, err := s.hnClient.GetTopStories()
	if err != nil {
		return fmt.Errorf("Unable to retrieve the top stories. %w", err)
	}
	for _, id := range storyIds {
		topStoriesChan <- id
	}
	close(topStoriesChan)

	if err := errGroup.Wait(); err != nil {
		return fmt.Errorf("An error occurred while saving the item. %w", err)
	}
	s.logger.Info("Finished processing stories")
	return nil
}

func (s *Service) saveItem(topStoriesChan <-chan int) error {
	for storyId := range topStoriesChan {
		item, err := s.hnClient.GetItem(storyId)
		if err != nil {
			return fmt.Errorf("An error occurred when trying to fetch the item. %w", err)
		} else if !item.Deleted && !item.Dead {
			log.Println(item)
		}
	}
	s.logger.Info("Finished looping through the channel for story ids")
	return nil
}
