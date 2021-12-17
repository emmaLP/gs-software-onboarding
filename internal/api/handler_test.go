package api

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"github.com/emmaLP/gs-software-onboarding/internal/model"
	hnModel "github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestGetAll(t *testing.T) {
	tests := map[string]struct {
		dbConfig             *model.DatabaseConfig
		dbMock               *database.Mock
		expectedMocks        func(t *testing.T, dbMock *database.Mock)
		expectedStatusCode   int
		expectedResultLength int
	}{
		"Successfully GetAll": {
			dbConfig:             nil,
			expectedStatusCode:   200,
			expectedResultLength: 2,
			dbMock:               &database.Mock{},
			expectedMocks: func(t *testing.T, dbMock *database.Mock) {
				dbMock.On("ListAll", context.TODO()).Return([]*hnModel.Item{
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
			expectedMocks: func(t *testing.T, dbMock *database.Mock) {
				dbMock.On("ListAll", context.TODO()).Return(nil, errors.New("Failed to find item"))
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewProduction()
			require.NoError(t, err)
			handler, err := NewHandler(context.TODO(), logger, testConfig.dbConfig, WithDatabaseClient(testConfig.dbMock))
			require.NoError(t, err)

			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.dbMock)
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
