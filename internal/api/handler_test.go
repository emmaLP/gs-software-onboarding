package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/emmaLP/gs-software-onboarding/internal/caching"
	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"github.com/emmaLP/gs-software-onboarding/internal/model"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetAll(t *testing.T) {
	tests := map[string]struct {
		dbConfig             *model.DatabaseConfig
		dbMock               *database.Mock
		cacheMock            *caching.Mock
		expectedMocks        func(t *testing.T, dbMock *caching.Mock)
		expectedStatusCode   int
		expectedResultLength int
	}{
		"Successfully ListAll": {
			dbConfig:             nil,
			expectedStatusCode:   200,
			expectedResultLength: 2,
			dbMock:               &database.Mock{},
			cacheMock:            &caching.Mock{},
			expectedMocks: func(t *testing.T, dbMock *caching.Mock) {
				dbMock.On("ListAll", context.TODO()).Return([]*commonModel.Item{
					{ID: 1, Type: "story"},
					{ID: 2, Type: "job"},
				}, nil)
			},
		},
		"Failed to get data": {
			dbConfig:             nil,
			expectedStatusCode:   500,
			expectedResultLength: 0,
			dbMock:               &database.Mock{},
			cacheMock:            &caching.Mock{},
			expectedMocks: func(t *testing.T, dbMock *caching.Mock) {
				dbMock.On("ListAll", context.TODO()).Return(nil, errors.New("Failed to find item"))
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewProduction()
			require.NoError(t, err)
			handler, err := NewHandler(logger, testConfig.cacheMock)
			require.NoError(t, err)

			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.cacheMock)
			}
			rec, eCtx := setupRequest(t, "/all")
			err = handler.GetAll(eCtx)
			require.NoError(t, err)

			if testConfig.expectedMocks != nil {
				testConfig.dbMock.AssertExpectations(t)
			}
			assert.Equal(t, testConfig.expectedStatusCode, rec.Code)
			if testConfig.expectedStatusCode == http.StatusOK {
				response := decodeRequest(t, rec.Body)

				assert.Equal(t, testConfig.expectedResultLength, len(response.Items))
			}
		})
	}
}

func TestListStories(t *testing.T) {
	tests := map[string]struct {
		dbConfig             *model.DatabaseConfig
		dbMock               *database.Mock
		cacheMock            *caching.Mock
		expectedMocks        func(t *testing.T, dbMock *caching.Mock)
		expectedStatusCode   int
		expectedResultLength int
	}{
		"Successfully ListStories": {
			dbConfig:             nil,
			expectedStatusCode:   200,
			expectedResultLength: 2,
			dbMock:               &database.Mock{},
			cacheMock:            &caching.Mock{},
			expectedMocks: func(t *testing.T, dbMock *caching.Mock) {
				dbMock.On("ListStories", context.TODO()).Return([]*commonModel.Item{
					{ID: 1, Type: "story"},
					{ID: 2, Type: "story"},
				}, nil)
			},
		},
		"Failed to get data": {
			dbConfig:             nil,
			expectedStatusCode:   500,
			expectedResultLength: 0,
			dbMock:               &database.Mock{},
			cacheMock:            &caching.Mock{},
			expectedMocks: func(t *testing.T, dbMock *caching.Mock) {
				dbMock.On("ListStories", context.TODO()).Return(nil, errors.New("Failed to find item"))
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewProduction()
			require.NoError(t, err)
			handler, err := NewHandler(logger, testConfig.cacheMock)
			require.NoError(t, err)

			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.cacheMock)
			}
			rec, eCtx := setupRequest(t, "/stories")
			err = handler.ListStories(eCtx)
			require.NoError(t, err)

			if testConfig.expectedMocks != nil {
				testConfig.dbMock.AssertExpectations(t)
			}
			assert.Equal(t, testConfig.expectedStatusCode, rec.Code)
			if testConfig.expectedStatusCode == http.StatusOK {
				response := decodeRequest(t, rec.Body)

				assert.Equal(t, testConfig.expectedResultLength, len(response.Items))
			}
		})
	}
}

func TestListJobs(t *testing.T) {
	tests := map[string]struct {
		dbConfig             *model.DatabaseConfig
		dbMock               *database.Mock
		cacheMock            *caching.Mock
		expectedMocks        func(t *testing.T, dbMock *caching.Mock)
		expectedStatusCode   int
		expectedResultLength int
	}{
		"Successfully ListJobs": {
			dbConfig:             nil,
			expectedStatusCode:   200,
			expectedResultLength: 2,
			dbMock:               &database.Mock{},
			cacheMock:            &caching.Mock{},
			expectedMocks: func(t *testing.T, dbMock *caching.Mock) {
				dbMock.On("ListJobs", context.TODO()).Return([]*commonModel.Item{
					{ID: 1, Type: "job"},
					{ID: 2, Type: "job"},
				}, nil)
			},
		},
		"Failed to get data": {
			dbConfig:             nil,
			expectedStatusCode:   500,
			expectedResultLength: 0,
			dbMock:               &database.Mock{},
			cacheMock:            &caching.Mock{},
			expectedMocks: func(t *testing.T, dbMock *caching.Mock) {
				dbMock.On("ListJobs", context.TODO()).Return(nil, errors.New("Failed to find item"))
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewProduction()
			require.NoError(t, err)
			handler, err := NewHandler(logger, testConfig.cacheMock)
			require.NoError(t, err)

			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.cacheMock)
			}
			rec, eCtx := setupRequest(t, "/jobs")
			err = handler.ListJobs(eCtx)
			require.NoError(t, err)

			if testConfig.expectedMocks != nil {
				testConfig.dbMock.AssertExpectations(t)
			}
			assert.Equal(t, testConfig.expectedStatusCode, rec.Code)
			if testConfig.expectedStatusCode == http.StatusOK {
				response := decodeRequest(t, rec.Body)

				assert.Equal(t, testConfig.expectedResultLength, len(response.Items))
			}
		})
	}
}
