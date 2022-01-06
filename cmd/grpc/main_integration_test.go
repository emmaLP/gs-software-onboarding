//go:build integration

package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/emmaLP/gs-software-onboarding/internal/caching"
	"github.com/emmaLP/gs-software-onboarding/internal/config"
	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
	"github.com/emmaLP/gs-software-onboarding/internal/model"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type testHandler struct {
	logger      *zap.Logger
	config      *model.Configuration
	dbClient    database.Client
	cacheClient caching.Client
}

func TestGrpcServer_ListStories(t *testing.T) {
	story := commonModel.Item{
		ID:      1,
		Dead:    false,
		Deleted: false,
		Type:    "story",
	}
	story2 := commonModel.Item{
		ID:      3,
		Dead:    false,
		Deleted: false,
		Type:    "story",
	}
	job := commonModel.Item{
		ID:      2,
		Dead:    true,
		Deleted: false,
		Type:    "job",
	}
	tests := map[string]struct {
		itemsToSave      []*commonModel.Item
		expectedResponse []*commonModel.Item
	}{
		"Successfully list 1 story": {
			itemsToSave:      []*commonModel.Item{&story, &job},
			expectedResponse: []*commonModel.Item{&story},
		},
		"Successfully list 2 stories": {
			itemsToSave:      []*commonModel.Item{&story, &job, &story2},
			expectedResponse: []*commonModel.Item{&story, &story2},
		},
	}

	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			handler := loadTestHandler(t)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			handler.cacheClient.FlushAll(ctx)
			for _, item := range testConfig.itemsToSave {
				handler.saveItemToDatabase(t, item)
			}

			client, err := grpc.NewClient(handler.config.Api.GrpcAddress, handler.logger)
			assert.NoError(t, err)
			defer client.Close()

			stories, err := client.ListStories(ctx)
			assert.NoError(t, err)
			assert.NotNil(t, stories)

			assert.Len(t, stories, len(testConfig.expectedResponse))
			assert.Equal(t, testConfig.expectedResponse, stories)

			t.Cleanup(func() {
				client.Close()
				handler.dbClient.CloseConnection(ctx)
				handler.cacheClient.Close()
			})
		})
	}
}

func loadTestHandler(t *testing.T) *testHandler {
	t.Helper()
	conf, err := config.LoadConfig(".")
	require.NoError(t, err)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	dbClient, err := database.New(context.TODO(), logger, &conf.Database)
	require.NoError(t, err)
	cacheClient, err := caching.New(context.TODO(), conf.Cache.Address, dbClient, logger, caching.WithTTL(10*time.Millisecond))
	require.NoError(t, err)
	return &testHandler{
		logger:      logger,
		config:      conf,
		dbClient:    dbClient,
		cacheClient: cacheClient,
	}
}

func (h *testHandler) saveItemToDatabase(t *testing.T, item *commonModel.Item) {
	t.Helper()
	err := h.dbClient.SaveItem(context.TODO(), item)
	require.NoError(t, err)
}

func (h *testHandler) address() string {
	return fmt.Sprintf("%s:%d", "localhost", h.config.Grpc.Port)
}
