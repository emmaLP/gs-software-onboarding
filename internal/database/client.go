package database

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	commonModel "github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

type Client interface {
	SaveItem(ctx context.Context, item *commonModel.Item) error
	ListAll(ctx context.Context) ([]*commonModel.Item, error)
	ListStories(ctx context.Context) ([]*commonModel.Item, error)
	ListJobs(ctx context.Context) ([]*commonModel.Item, error)
	CloseConnection(ctx context.Context)
}

type database struct {
	mongoClient  *mongo.Client
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

	database := &database{
		mongoClient:  client,
		logger:       logger,
		databaseName: config.Name,
	}
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("mongo refused to connect: %v %w", ctx.Err(), err)
		default:
			err := client.Ping(ctx, readpref.Primary())
			if err == nil {
				log.Print("mongo is now connected")
				return database, nil
			}
		}
	}
}

func (d *database) SaveItem(ctx context.Context, item *commonModel.Item) error {
	collection := d.getCollection("items")
	opts := options.Update().SetUpsert(true)

	update := bson.M{
		"$set": item,
	}
	_, err := collection.UpdateOne(ctx, bson.M{"id": item.ID}, update, opts)
	if err != nil {
		return fmt.Errorf("Unable to save item. %w", err)
	}
	d.logger.Info("Item saved successfully", zap.Int("ID", item.ID))
	return nil
}

func (d *database) ListAll(ctx context.Context) ([]*commonModel.Item, error) {
	return d.find(ctx, bson.M{})
}

func (d *database) ListStories(ctx context.Context) ([]*commonModel.Item, error) {
	filter := bson.M{"type": "story"}
	return d.find(ctx, filter)
}

func (d *database) ListJobs(ctx context.Context) ([]*commonModel.Item, error) {
	filter := bson.M{"type": "job"}
	return d.find(ctx, filter)
}

func (d *database) find(ctx context.Context, filter interface{}) ([]*commonModel.Item, error) {
	collection := d.getCollection("items")
	all, err := collection.Find(ctx, filter)
	var items []*commonModel.Item
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve items. %w", err)
	}

	if err = all.All(ctx, &items); err != nil {
		return nil, fmt.Errorf("Failed to retrieve items within cursor. %w", err)
	}
	return items, nil
}

func (d *database) CloseConnection(ctx context.Context) {
	d.logger.Debug("Closing database connection")
	err := d.mongoClient.Disconnect(ctx)
	if err != nil {
		d.logger.Error("Failed to close connection to database", zap.Error(err))
	}
}

func (d *database) getCollection(collectionName string) *mongo.Collection {
	database := d.mongoClient.Database(d.databaseName)
	return database.Collection(collectionName)
}
