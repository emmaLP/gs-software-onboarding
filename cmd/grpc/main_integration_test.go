//go:build integration

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/emmaLP/gs-software-onboarding/internal/caching"
	"github.com/emmaLP/gs-software-onboarding/internal/config"
	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"github.com/emmaLP/gs-software-onboarding/internal/model"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	pb "github.com/emmaLP/gs-software-onboarding/pkg/grpc/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type testHandler struct {
	logger      *zap.Logger
	config      *model.Configuration
	dbClient    database.Client
	cacheClient caching.Client
}

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
			handler := loadTestHandler(t)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			handler.cacheClient.FlushAll(ctx)
			for _, item := range testConfig.itemsToSave {
				handler.saveItemToDatabase(t, item)
			}

			grpcConnection, err := grpc.Dial(handler.address(), grpc.WithInsecure())
			assert.NoError(t, err)
			defer grpcConnection.Close()

			client := pb.NewAPIClient(grpcConnection)

			storiesSteam, err := client.ListStories(ctx, &emptypb.Empty{})
			assert.NoError(t, err)
			assert.NotNil(t, storiesSteam)

			for _, expectedItem := range testConfig.expectedResponse {
				out, err := storiesSteam.Recv()
				if err == nil {
					actualItem := commonModel.PItemToItem(out)
					assert.Equal(t, expectedItem, &actualItem)
				} else {
					t.Fatal("Unexpected error:", err)
				}
			}

			lastRecv, err := storiesSteam.Recv()
			log.Println(lastRecv)
			assert.Equal(t, io.EOF, err)

			t.Cleanup(func() {
				grpcConnection.Close()
				handler.dbClient.CloseConnection(ctx)
				handler.cacheClient.Close()
			})
		})
	}
}

func loadTestHandler(t *testing.T) *testHandler {
	t.Helper()
	conf, err := config.LoadConfig(".")
	require.NoError(t, err)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	dbClient, err := database.New(context.TODO(), logger, &conf.Database)
	require.NoError(t, err)
	cacheClient, err := caching.New(context.TODO(), conf.Cache.Address, dbClient, logger, caching.WithTTL(10*time.Millisecond))
	require.NoError(t, err)
	return &testHandler{
		logger:      logger,
		config:      conf,
		dbClient:    dbClient,
		cacheClient: cacheClient,
	}
}

func (h *testHandler) saveItemToDatabase(t *testing.T, item *commonModel.Item) {
	t.Helper()
	err := h.dbClient.SaveItem(context.TODO(), item)
	require.NoError(t, err)
}

func (h *testHandler) address() string {
	return fmt.Sprintf("%s:%d", "localhost", h.config.Grpc.Port)
}

func (h *testHandler) dropDatabase(ctx context.Context) {
	dbConfig := h.config.Database
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", dbConfig.Host, dbConfig.Port))
	if strings.TrimSpace(dbConfig.Username) != "" && strings.TrimSpace(dbConfig.Password) != "" {
		credentials := options.Credential{
			Username: dbConfig.Username,
			Password: dbConfig.Password,
		}
		opts = opts.SetAuth(credentials)
	}
	client, _ := mongo.Connect(ctx, opts)
	defer client.Disconnect(ctx)

	database := client.Database(dbConfig.Name)
	collection := database.Collection("items")
	collection.Drop(ctx)
	database.Drop(ctx)
}
