package database

import (
	"context"
	"testing"

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
