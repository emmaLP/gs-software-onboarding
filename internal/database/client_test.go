package database

import (
	"context"
	"testing"

	hnModel "github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNew(t *testing.T) {
	logger, err := zap.NewProduction()
	require.NoError(t, err)
	client, err := New(context.TODO(), logger, &model.DatabaseConfig{})
	require.NoError(t, err)
	assert.NotNil(t, client)
}

func TestCloseConnection(t *testing.T) {
	logger, err := zap.NewProduction()
	require.NoError(t, err)
	client, err := New(context.TODO(), logger, &model.DatabaseConfig{})
	require.NoError(t, err)

	client.CloseConnection()
}

func TestSaveItem(t *testing.T) {
	tests := map[string]struct {
		config         *model.DatabaseConfig
		expectedResult *hnModel.Item
		expectedErr    string
	}{}

	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewProduction()
			require.NoError(t, err)
			client, err := New(context.TODO(), logger, testConfig.config)
			require.NoError(t, err)
			assert.NotNil(t, client)
		})
	}
}
