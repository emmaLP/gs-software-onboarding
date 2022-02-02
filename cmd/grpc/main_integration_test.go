//go:build integration

package main

import (
	"context"
	"testing"
	"time"

	"github.com/emmaLP/gs-software-onboarding/pkg/test"

	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/stretchr/testify/assert"
)

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
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			handler := test.LoadTestHandler(t, ctx)
			handler.FlushCache(ctx)
			for _, item := range testConfig.itemsToSave {
				handler.SaveItemToDatabase(ctx, item)
			}

			client, err := grpc.NewClient(handler.Config.GrpcClient.GrpcAddress, handler.Logger)
			assert.NoError(t, err)
			defer client.Close()

			stories, err := client.ListStories(ctx)
			assert.NoError(t, err)
			assert.NotNil(t, stories)

			assert.Len(t, stories, len(testConfig.expectedResponse))
			assert.Equal(t, testConfig.expectedResponse, stories)

			t.Cleanup(func() {
				client.Close()
				handler.CloseConnections(ctx)
			})
		})
	}
}
