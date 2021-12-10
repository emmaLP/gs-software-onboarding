package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/emmaLP/gs-software-onboarding/pkg/hackernews"
	"go.uber.org/zap"
)

type service struct {
	logger          *zap.Logger
	numberOfWorkers int
	hnClient        hackernews.Client
	dbClient        database.Client
}

type Client interface {
	processStories(ctx context.Context)
}

func NewService(logger *zap.Logger, config *model.Configuration, dbClient database.Client, hnClient ...hackernews.Client) (*service, error) {
	var hackernewsClient hackernews.Client
	if len(hnClient) == 0 {
		var err error
		hackernewsClient, err = hackernews.New(config.Consumer.BaseUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("Unable to create HackerNew client: %w", err)
		}
	} else {
		//Only ever get the first one, ignore the rest
		hackernewsClient = hnClient[0]
	}

	return &service{
		logger:          logger,
		numberOfWorkers: config.Consumer.NumberOfWorkers,
		hnClient:        hackernewsClient,
		dbClient:        dbClient,
	}, nil
}

func (s *service) processStories(ctx context.Context) error {
	s.logger.Info("Processing stories")

	var wg sync.WaitGroup
	topStoriesChan := make(chan int)

	for i := 0; i < s.numberOfWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.saveItem(topStoriesChan)
		}()
	}

	s.logger.Info("Finished getting the items")

	storyIds, err := s.hnClient.GetTopStories()
	if err != nil {
		return fmt.Errorf("Unable to retrieve the top stories. %w", err)
	}
	for _, id := range storyIds {
		select {
		case <-ctx.Done():
		case topStoriesChan <- id:
		}
	}
	close(topStoriesChan)

	wg.Wait()
	s.logger.Info("Finished processing stories")
	return nil
}

func (s *service) saveItem(topStoriesChan <-chan int) {
	for storyId := range topStoriesChan {
		item, err := s.hnClient.GetItem(storyId)
		if err != nil {
			s.logger.Error("An error occurred when trying to fetch the item.", zap.Error(err))
		} else if !item.Deleted && !item.Dead {
			log.Println(item)
			err := s.dbClient.SaveItem(item)
			if err != nil {
				s.logger.Error("Failed to save item", zap.Error(err))
			}
		}
	}
	s.logger.Info("Finished looping through the channel for story ids")
}
