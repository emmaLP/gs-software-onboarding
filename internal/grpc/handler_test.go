package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/emmaLP/gs-software-onboarding/internal/caching"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	pbMock "github.com/emmaLP/gs-software-onboarding/pkg/grpc/proto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestListMethods(t *testing.T) {
	items := []*commonModel.Item{
		{ID: 1, Type: "story"},
		{ID: 2, Type: "job"},
	}
	tests := map[string]struct {
		cacheMock         *caching.Mock
		expectedMocks     func(t *testing.T, cacheMock *caching.Mock)
		itemsToSend       []*commonModel.Item
		listAllServer     func(t *testing.T) *pbMock.MockAPI_ListAllServer
		listStoriesServer func(t *testing.T) *pbMock.MockAPI_ListStoriesServer
		listJobsServer    func(t *testing.T) *pbMock.MockAPI_ListJobsServer
	}{
		"ListAll Successfully": {
			cacheMock:   &caching.Mock{},
			itemsToSend: items,
			expectedMocks: func(t *testing.T, cacheMock *caching.Mock) {
				cacheMock.On("ListAll", context.TODO()).Return(items, nil)
			},
			listAllServer: func(t *testing.T) *pbMock.MockAPI_ListAllServer {
				controller := gomock.NewController(t)
				defer controller.Finish()
				return pbMock.NewMockAPI_ListAllServer(controller)
			},
		},
		"ListStories Successfully": {
			cacheMock:   &caching.Mock{},
			itemsToSend: items,
			expectedMocks: func(t *testing.T, cacheMock *caching.Mock) {
				cacheMock.On("ListStories", context.TODO()).Return(items, nil)
			},
			listStoriesServer: func(t *testing.T) *pbMock.MockAPI_ListStoriesServer {
				controller := gomock.NewController(t)
				defer controller.Finish()
				return pbMock.NewMockAPI_ListStoriesServer(controller)
			},
		},
		"ListJobs Successfully": {
			cacheMock:   &caching.Mock{},
			itemsToSend: items,
			expectedMocks: func(t *testing.T, cacheMock *caching.Mock) {
				cacheMock.On("ListJobs", context.TODO()).Return(items, nil)
			},
			listJobsServer: func(t *testing.T) *pbMock.MockAPI_ListJobsServer {
				controller := gomock.NewController(t)
				defer controller.Finish()
				return pbMock.NewMockAPI_ListJobsServer(controller)
			},
		},
	}

	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.cacheMock)
			}
			handler := Handler{
				itemCache: testConfig.cacheMock,
			}
			var err error
			if testConfig.listAllServer != nil {
				apiServer := testConfig.listAllServer(t)
				apiServer.EXPECT().Context().Return(context.TODO())
				for _, item := range testConfig.itemsToSend {
					apiServer.EXPECT().Send(
						gomock.Eq(commonModel.ItemToPItem(*item)),
					).Return(nil)
				}

				err = handler.ListAll(nil, apiServer)
			} else if testConfig.listJobsServer != nil {
				apiServer := testConfig.listJobsServer(t)
				apiServer.EXPECT().Context().Return(context.TODO())
				for _, item := range testConfig.itemsToSend {
					apiServer.EXPECT().Send(
						gomock.Eq(commonModel.ItemToPItem(*item)),
					).Return(nil)
				}

				err = handler.ListJobs(nil, apiServer)
			} else if testConfig.listStoriesServer != nil {
				apiServer := testConfig.listStoriesServer(t)
				apiServer.EXPECT().Context().Return(context.TODO())
				for _, item := range testConfig.itemsToSend {
					apiServer.EXPECT().Send(
						gomock.Eq(commonModel.ItemToPItem(*item)),
					).Return(nil)
				}

				err = handler.ListStories(nil, apiServer)
			}

			assert.NoError(t, err)
			if testConfig.expectedMocks != nil {
				testConfig.cacheMock.AssertExpectations(t)
			}
		})
	}
}

func TestHandler_SaveItem(t *testing.T) {
	tests := map[string]struct {
		dbMock             *database.Mock
		itemToSave         *pbMock.Item
		expectedMocks      func(t *testing.T, dbMock *database.Mock)
		expectedErrMessage string
	}{
		"Successful save": {
			dbMock:     &database.Mock{},
			itemToSave: &pbMock.Item{Id: 1},
			expectedMocks: func(t *testing.T, dbMock *database.Mock) {
				dbMock.On("SaveItem", context.TODO(), mock.Anything).Return(nil)
			},
		},
		"Unsuccessful save": {
			dbMock:             &database.Mock{},
			itemToSave:         &pbMock.Item{Id: 1},
			expectedErrMessage: "Failed to save.",
			expectedMocks: func(t *testing.T, dbMock *database.Mock) {
				dbMock.On("SaveItem", context.TODO(), mock.Anything).Return(errors.New("Failed to save."))
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			if testConfig.expectedMocks != nil {
				testConfig.expectedMocks(t, testConfig.dbMock)
			}
			logger, err := zap.NewDevelopment()
			require.NoError(t, err)

			handler := NewHandler(nil, testConfig.dbMock, logger)
			itemResponse, err := handler.SaveItem(context.TODO(), testConfig.itemToSave)
			if testConfig.expectedErrMessage != "" {
				assert.EqualErrorf(t, err, testConfig.expectedErrMessage, "Request failed should be: %v, got: %v", testConfig.expectedErrMessage, err)
				assert.Equal(t, testConfig.itemToSave.Id, itemResponse.Id)
				assert.False(t, itemResponse.Success)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testConfig.itemToSave.Id, itemResponse.Id)
				assert.True(t, itemResponse.Success)
			}
			if testConfig.expectedMocks != nil {
				testConfig.dbMock.AssertExpectations(t)
			}
		})
	}
}
