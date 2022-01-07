package grpc

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	pb "github.com/emmaLP/gs-software-onboarding/pkg/grpc/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestListAll(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	tests := map[string]struct {
		grpcClient         *pb.MockAPIClient
		listAllClient      *pb.MockAPI_ListAllClient
		expectedMocks      func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListAllClient)
		expectedNumItems   int
		expectedErrMessage string
	}{
		"Successfully ListAll": {
			grpcClient:       pb.NewMockAPIClient(controller),
			listAllClient:    pb.NewMockAPI_ListAllClient(controller),
			expectedNumItems: 2,
			expectedMocks: func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListAllClient) {
				mock.EXPECT().ListAll(gomock.Any(), gomock.Any()).Return(ret, nil)
				ret.EXPECT().Recv().Return(&pb.Item{
					Id:   1,
					Type: "story",
				}, nil)
				ret.EXPECT().Recv().Return(
					&pb.Item{
						Id:   2,
						Type: "job",
					}, nil)
				ret.EXPECT().Recv().Return(nil, io.EOF)
			},
		},
		"Error in grpc client": {
			grpcClient:         pb.NewMockAPIClient(controller),
			expectedNumItems:   0,
			expectedErrMessage: "An error occurred when streaming all. Failed to stream",
			expectedMocks: func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListAllClient) {
				mock.EXPECT().ListAll(gomock.Any(), gomock.Any()).Return(ret, errors.New("Failed to stream"))
			},
		},
		"Error in streaming": {
			grpcClient:         pb.NewMockAPIClient(controller),
			listAllClient:      pb.NewMockAPI_ListAllClient(controller),
			expectedNumItems:   0,
			expectedErrMessage: "receiving item from server. failed",
			expectedMocks: func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListAllClient) {
				mock.EXPECT().ListAll(gomock.Any(), gomock.Any()).Return(ret, nil)
				ret.EXPECT().Recv().Return(nil, errors.New("failed"))
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewDevelopment()
			require.NoError(t, err)

			c := client{
				grpcClient: testConfig.grpcClient,
				logger:     logger,
			}
			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.grpcClient, testConfig.listAllClient)
			}

			allItems, err := c.ListAll(context.TODO())
			if strings.TrimSpace(testConfig.expectedErrMessage) != "" {
				assert.EqualErrorf(t, err, testConfig.expectedErrMessage, "Request failed should be: %v, got: %v", testConfig.expectedErrMessage, err)
				assert.Nil(t, allItems)
			} else {
				require.NoError(t, err)
				assert.Len(t, allItems, testConfig.expectedNumItems)
			}
		})
	}
}

func TestListStories(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	tests := map[string]struct {
		grpcClient         *pb.MockAPIClient
		listClient         *pb.MockAPI_ListStoriesClient
		expectedMocks      func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListStoriesClient)
		expectedNumItems   int
		expectedErrMessage string
	}{
		"Successfully ListStories": {
			grpcClient:       pb.NewMockAPIClient(controller),
			listClient:       pb.NewMockAPI_ListStoriesClient(controller),
			expectedNumItems: 2,
			expectedMocks: func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListStoriesClient) {
				mock.EXPECT().ListStories(gomock.Any(), gomock.Any()).Return(ret, nil)
				ret.EXPECT().Recv().Return(&pb.Item{
					Id:   1,
					Type: "story",
				}, nil)
				ret.EXPECT().Recv().Return(
					&pb.Item{
						Id:   2,
						Type: "story",
					}, nil)
				ret.EXPECT().Recv().Return(nil, io.EOF)
			},
		},
		"Error in grpc client": {
			grpcClient:         pb.NewMockAPIClient(controller),
			expectedNumItems:   0,
			expectedErrMessage: "An error occurred when streaming stories. Failed to stream",
			expectedMocks: func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListStoriesClient) {
				mock.EXPECT().ListStories(gomock.Any(), gomock.Any()).Return(ret, errors.New("Failed to stream"))
			},
		},
		"Error in streaming": {
			grpcClient:         pb.NewMockAPIClient(controller),
			listClient:         pb.NewMockAPI_ListStoriesClient(controller),
			expectedNumItems:   0,
			expectedErrMessage: "receiving item from server. failed",
			expectedMocks: func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListStoriesClient) {
				mock.EXPECT().ListStories(gomock.Any(), gomock.Any()).Return(ret, nil)
				ret.EXPECT().Recv().Return(nil, errors.New("failed"))
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewDevelopment()
			require.NoError(t, err)

			c := client{
				grpcClient: testConfig.grpcClient,
				logger:     logger,
			}
			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.grpcClient, testConfig.listClient)
			}

			allItems, err := c.ListStories(context.TODO())
			if strings.TrimSpace(testConfig.expectedErrMessage) != "" {
				assert.EqualErrorf(t, err, testConfig.expectedErrMessage, "Request failed should be: %v, got: %v", testConfig.expectedErrMessage, err)
				assert.Nil(t, allItems)
			} else {
				require.NoError(t, err)
				assert.Len(t, allItems, testConfig.expectedNumItems)
			}
		})
	}
}

func TestListJobs(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	tests := map[string]struct {
		grpcClient         *pb.MockAPIClient
		listClient         *pb.MockAPI_ListJobsClient
		expectedMocks      func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListJobsClient)
		expectedNumItems   int
		expectedErrMessage string
	}{
		"Successfully ListJobs": {
			grpcClient:       pb.NewMockAPIClient(controller),
			listClient:       pb.NewMockAPI_ListJobsClient(controller),
			expectedNumItems: 2,
			expectedMocks: func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListJobsClient) {
				mock.EXPECT().ListJobs(gomock.Any(), gomock.Any()).Return(ret, nil)
				ret.EXPECT().Recv().Return(&pb.Item{
					Id:   1,
					Type: "job",
				}, nil)
				ret.EXPECT().Recv().Return(
					&pb.Item{
						Id:   2,
						Type: "job",
					}, nil)
				ret.EXPECT().Recv().Return(nil, io.EOF)
			},
		},
		"Error in grpc client": {
			grpcClient:         pb.NewMockAPIClient(controller),
			expectedNumItems:   0,
			expectedErrMessage: "An error occurred when streaming jobs. Failed to stream",
			expectedMocks: func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListJobsClient) {
				mock.EXPECT().ListJobs(gomock.Any(), gomock.Any()).Return(ret, errors.New("Failed to stream"))
			},
		},
		"Error in streaming": {
			grpcClient:         pb.NewMockAPIClient(controller),
			listClient:         pb.NewMockAPI_ListJobsClient(controller),
			expectedNumItems:   0,
			expectedErrMessage: "receiving item from server. failed",
			expectedMocks: func(t *testing.T, mock *pb.MockAPIClient, ret *pb.MockAPI_ListJobsClient) {
				mock.EXPECT().ListJobs(gomock.Any(), gomock.Any()).Return(ret, nil)
				ret.EXPECT().Recv().Return(nil, errors.New("failed"))
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewDevelopment()
			require.NoError(t, err)

			c := client{
				grpcClient: testConfig.grpcClient,
				logger:     logger,
			}
			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.grpcClient, testConfig.listClient)
			}

			allItems, err := c.ListJobs(context.TODO())
			if strings.TrimSpace(testConfig.expectedErrMessage) != "" {
				assert.EqualErrorf(t, err, testConfig.expectedErrMessage, "Request failed should be: %v, got: %v", testConfig.expectedErrMessage, err)
				assert.Nil(t, allItems)
			} else {
				require.NoError(t, err)
				assert.Len(t, allItems, testConfig.expectedNumItems)
			}
		})
	}
}
