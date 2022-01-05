package caching

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/emmaLP/gs-software-onboarding/internal/database"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestListAll(t *testing.T) {
	redisServer, err := miniredis.Run()
	require.NoError(t, err)
	tests := map[string]struct {
		dbMock             *database.Mock
		expectedMocks      func(t *testing.T, dbMock *database.Mock)
		fromCache          bool
		expectedItemsCount int
	}{
		"From cache": {
			dbMock:             &database.Mock{},
			fromCache:          true,
			expectedItemsCount: 2,
			expectedMocks: func(t *testing.T, dbMock *database.Mock) {
				dbMock.On("ListAll", context.TODO()).Return([]*commonModel.Item{
					{ID: 1, Type: "story"},
					{ID: 2, Type: "job"},
				}, nil).Once()
			},
		},
		"From database": {
			dbMock:             &database.Mock{},
			expectedItemsCount: 2,
			expectedMocks: func(t *testing.T, dbMock *database.Mock) {
				dbMock.On("ListAll", context.TODO()).Return([]*commonModel.Item{
					{ID: 1, Type: "story"},
					{ID: 2, Type: "job"},
				}, nil).Times(2)
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewDevelopment()
			require.NoError(t, err)
			cacheClient, err := New(context.TODO(), redisServer.Addr(), testConfig.dbMock, logger, WithTTL(time.Minute))
			require.NoError(t, err)

			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.dbMock)
			}

			// Prepopulate cache
			items, err := cacheClient.ListAll(context.TODO())
			require.NoError(t, err)
			assert.Equal(t, testConfig.expectedItemsCount, len(items))

			if !testConfig.fromCache {
				// Clear the cache if test pulling from the db
				cacheClient.FlushAll(context.TODO())
			}

			items, err = cacheClient.ListAll(context.TODO())
			require.NoError(t, err)
			assert.Equal(t, testConfig.expectedItemsCount, len(items))

			if testConfig.expectedMocks != nil {
				testConfig.dbMock.AssertExpectations(t)
			}
			t.Cleanup(func() {
				cacheClient.FlushAll(context.TODO())
				cacheClient.Close()
			})
		})
	}
}

func TestListStories(t *testing.T) {
	redisServer, err := miniredis.Run()
	require.NoError(t, err)
	tests := map[string]struct {
		dbMock             *database.Mock
		expectedMocks      func(t *testing.T, dbMock *database.Mock)
		fromCache          bool
		expectedItemsCount int
	}{
		"From cache": {
			dbMock:             &database.Mock{},
			fromCache:          true,
			expectedItemsCount: 2,
			expectedMocks: func(t *testing.T, dbMock *database.Mock) {
				dbMock.On("ListStories", context.TODO()).Return([]*commonModel.Item{
					{ID: 1, Type: "story"},
					{ID: 2, Type: "story"},
				}, nil).Once()
			},
		},
		"From database": {
			dbMock:             &database.Mock{},
			expectedItemsCount: 1,
			expectedMocks: func(t *testing.T, dbMock *database.Mock) {
				dbMock.On("ListStories", context.TODO()).Return([]*commonModel.Item{
					{ID: 1, Type: "story"},
				}, nil).Times(2)
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewDevelopment()
			require.NoError(t, err)
			cacheClient, err := New(context.TODO(), redisServer.Addr(), testConfig.dbMock, logger, WithTTL(time.Minute))
			require.NoError(t, err)

			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.dbMock)
			}

			// Prepopulate cache
			items, err := cacheClient.ListStories(context.TODO())
			require.NoError(t, err)
			assert.Equal(t, testConfig.expectedItemsCount, len(items))

			if !testConfig.fromCache {
				// Clear the cache if test pulling from the db
				cacheClient.FlushAll(context.TODO())
			}

			items, err = cacheClient.ListStories(context.TODO())
			require.NoError(t, err)
			assert.Equal(t, testConfig.expectedItemsCount, len(items))

			if testConfig.expectedMocks != nil {
				testConfig.dbMock.AssertExpectations(t)
			}
			t.Cleanup(func() {
				cacheClient.FlushAll(context.TODO())
				cacheClient.Close()
			})
		})
	}
}

func TestListJobs(t *testing.T) {
	redisServer, err := miniredis.Run()
	require.NoError(t, err)
	tests := map[string]struct {
		dbMock             *database.Mock
		expectedMocks      func(t *testing.T, dbMock *database.Mock)
		fromCache          bool
		expectedItemsCount int
	}{
		"From cache": {
			dbMock:             &database.Mock{},
			fromCache:          true,
			expectedItemsCount: 2,
			expectedMocks: func(t *testing.T, dbMock *database.Mock) {
				dbMock.On("ListJobs", context.TODO()).Return([]*commonModel.Item{
					{ID: 1, Type: "jobs"},
					{ID: 2, Type: "jobs"},
				}, nil).Once()
			},
		},
		"From database": {
			dbMock:             &database.Mock{},
			expectedItemsCount: 3,
			expectedMocks: func(t *testing.T, dbMock *database.Mock) {
				dbMock.On("ListJobs", context.TODO()).Return([]*commonModel.Item{
					{ID: 1, Type: "jobs"},
					{ID: 2, Type: "jobs"},
					{ID: 3, Type: "jobs"},
				}, nil).Times(2)
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewProduction()
			require.NoError(t, err)
			cacheClient, err := New(context.TODO(), redisServer.Addr(), testConfig.dbMock, logger, WithTTL(time.Minute))
			require.NoError(t, err)

			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.dbMock)
			}

			// Prepopulate cache
			items, err := cacheClient.ListJobs(context.TODO())
			require.NoError(t, err)
			assert.Equal(t, testConfig.expectedItemsCount, len(items))

			if !testConfig.fromCache {
				// Clear the cache if test pulling from the db
				cacheClient.FlushAll(context.TODO())
			}

			items, err = cacheClient.ListJobs(context.TODO())
			require.NoError(t, err)
			assert.Equal(t, testConfig.expectedItemsCount, len(items))

			if testConfig.expectedMocks != nil {
				testConfig.dbMock.AssertExpectations(t)
			}
			t.Cleanup(func() {
				cacheClient.FlushAll(context.TODO())
				cacheClient.Close()
			})
		})
	}
}
