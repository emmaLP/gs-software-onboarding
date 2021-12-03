package database

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/emmaLP/gs-software-onboarding/internal/model"

	hnModel "github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	SaveItem(item hnModel.Item) error
	CloseConnection()
}

type database struct {
	mongoClient *mongo.Client
	context     context.Context
	logger      *zap.Logger
}

const mongoUriTemplate = "mongodb://%s:%s"

func New(ctx context.Context, logger *zap.Logger, config *model.DatabaseConfig) (*database, error) {
	mongoUri := fmt.Sprintf(mongoUriTemplate, config.Host, config.Port)
	credentials := options.Credential{
		Username: config.Username,
		Password: config.Password,
	}
	opts := options.Client().ApplyURI(mongoUri).SetAuth(credentials)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("An error occurred when trying to connect to mongo. %w", err)
	}
	return &database{
		mongoClient: client,
		context:     ctx,
		logger:      logger,
	}, nil
}

func (d *database) SaveItem(item hnModel.Item) error {
	return nil
}

func (d *database) CloseConnection() {
	d.logger.Debug("Closing database connection")
	err := d.mongoClient.Disconnect(d.context)
	if err != nil {
		d.logger.Error("Failed to close connection to database", zap.Error(err))
	}
}
