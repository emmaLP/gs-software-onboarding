package grpc

import (
	"strings"
	"testing"

	"github.com/emmaLP/gs-software-onboarding/internal/caching"
	"github.com/emmaLP/gs-software-onboarding/internal/database"
	pbMock "github.com/emmaLP/gs-software-onboarding/pkg/grpc/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func TestNewServer(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger, err := zap.NewDevelopment()
	require.NoError(t, err)

	apiServer := NewServer(1234, logger, pbMock.NewMockAPIServer(controller))
	assert.NotNil(t, apiServer)
	assert.Equal(t, 1234, apiServer.port)
}

func TestStart(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	logger, err := zap.NewDevelopment()

	require.NoError(t, err)
	tests := map[string]struct {
		port               int
		expectedErrMessage string
	}{
		"Successfully Start Gprc Server": {
			port: 18001,
		},
		"Invalid port": {
			port:               1024 * 1024,
			expectedErrMessage: "failed to listen on port, listen tcp: address 1048576: invalid port",
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			_ = pbMock.NewMockAPIServer(controller)
			apiServer := NewServer(testConfig.port, logger, NewHandler(&caching.Mock{}, &database.Mock{}, logger))

			if strings.TrimSpace(testConfig.expectedErrMessage) != "" {
				_, err := apiServer.Start()
				assert.EqualError(t, err, testConfig.expectedErrMessage)
			} else {
				go func() {
					var s *grpc.Server
					s, err = apiServer.Start()
					defer s.Stop()
				}()
				require.NoError(t, err)
			}
		})
	}
}
