package api

import (
	"testing"

	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewServer(t *testing.T) {
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	server, err := NewServer(logger, &grpc.Mock{})
	require.NoError(t, err)

	assert.NotNil(t, server)
}
