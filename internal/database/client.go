package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	hnModel "github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Client interface {
	SaveItem(item *hnModel.Item) error
	CloseConnection()
}

type database struct {
	mongoClient  *mongo.Client
	context      context.Context
	logger       *zap.Logger
	databaseName string
}

const mongoUriTemplate = "mongodb://%s:%s"

func New(ctx context.Context, logger *zap.Logger, config *model.DatabaseConfig) (*database, error) {
	mongoUri := fmt.Sprintf(mongoUriTemplate, config.Host, config.Port)

	opts := options.Client().ApplyURI(mongoUri)
	if strings.TrimSpace(config.Username) != "" && strings.TrimSpace(config.Password) != "" {
		credentials := options.Credential{
			Username: config.Username,
			Password: config.Password,
		}
		opts = opts.SetAuth(credentials)
	}
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("An error occurred when trying to connect to mongo. %w", err)
	}
	return &database{
		mongoClient:  client,
		context:      ctx,
		logger:       logger,
		databaseName: config.Name,
	}, nil
}

func (d *database) SaveItem(item *hnModel.Item) error {
	collection := d.getCollection("items")
	opts := options.Update().SetUpsert(true)

	update := bson.M{
		"$set": item,
	}
	_, err := collection.UpdateOne(d.context, bson.M{"id": item.ID}, update, opts)
	if err != nil {
		return fmt.Errorf("Unable to save item. %w", err)
	}
	d.logger.Info("Item %d saved successfully", zap.Int("ID", item.ID))
	return nil
}

func (d *database) CloseConnection() {
	d.logger.Debug("Closing database connection")
	err := d.mongoClient.Disconnect(d.context)
	if err != nil {
		d.logger.Error("Failed to close connection to database", zap.Error(err))
	}
}

func (d *database) getCollection(collectionName string) *mongo.Collection {
	database := d.mongoClient.Database(d.databaseName)
	return database.Collection(collectionName)
}
