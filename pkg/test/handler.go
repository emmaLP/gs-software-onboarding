//go:build integration

package test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/emmaLP/gs-software-onboarding/internal/caching"
	"github.com/emmaLP/gs-software-onboarding/internal/config"
	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"github.com/emmaLP/gs-software-onboarding/internal/model"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type testHandler struct {
	Logger      *zap.Logger
	Config      *model.Configuration
	dbClient    database.Client
	cacheClient caching.Client
	t           *testing.T
}

func LoadTestHandler(t *testing.T, ctx context.Context) *testHandler {
	t.Helper()
	conf, err := config.LoadConfig(".")
	require.NoError(t, err)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	dbClient, err := database.New(ctx, logger, &conf.Database)
	require.NoError(t, err)
	cacheClient, err := caching.New(ctx, conf.Cache.Address, dbClient, logger, caching.WithTTL(10*time.Millisecond))
	require.NoError(t, err)
	return &testHandler{
		Logger:      logger,
		Config:      conf,
		dbClient:    dbClient,
		cacheClient: cacheClient,
		t:           t,
	}
}

func (h *testHandler) CloseConnections(ctx context.Context) {
	h.t.Helper()
	h.dbClient.CloseConnection(ctx)
	h.cacheClient.Close()
}

func (h *testHandler) SaveItemToDatabase(ctx context.Context, item *commonModel.Item) {
	h.t.Helper()
	err := h.dbClient.SaveItem(ctx, item)
	require.NoError(h.t, err)
}

func (h *testHandler) GRPCAddress() string {
	h.t.Helper()
	return fmt.Sprintf("%s:%d", "localhost", h.Config.Grpc.Port)
}

func (h *testHandler) FlushCache(ctx context.Context) {
	h.t.Helper()
	h.cacheClient.FlushAll(ctx)
}

func (h *testHandler) DropDatabase(ctx context.Context) {
	uri := fmt.Sprintf("mongodb://%s:%s", h.Config.Database.Host, h.Config.Database.Port)
	opts := options.Client().ApplyURI(uri)
	if strings.TrimSpace(h.Config.Database.Username) != "" && strings.TrimSpace(h.Config.Database.Password) != "" {
		credentials := options.Credential{
			Username: h.Config.Database.Username,
			Password: h.Config.Database.Password,
		}
		opts = opts.SetAuth(credentials)
	}
	client, _ := mongo.Connect(ctx, opts)

	defer client.Disconnect(ctx)
	database := client.Database(h.Config.Database.Name)
	collection := database.Collection("items")
	collection.Drop(ctx)
	database.Drop(ctx)
}
