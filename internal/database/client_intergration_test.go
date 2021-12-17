package database

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	hnModel "github.com/emmaLP/gs-software-onboarding/pkg/hackernews/model"
	tc "github.com/romnn/testcontainers"
	tcMongo "github.com/romnn/testcontainers/mongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"go.uber.org/zap"
)

func setupMongo(ctx context.Context) (mongoC testcontainers.Container, Config tcMongo.DBConfig, err error) {
	container, dbConfig, err := tcMongo.StartMongoContainer(ctx, tcMongo.ContainerOptions{
		ContainerOptions: tc.ContainerOptions{CollectLogs: false},
		User:             "test",
		Password:         "test",
	})
	if err != nil {
		return nil, tcMongo.DBConfig{}, err
	}

	return container, dbConfig, nil
}

func TestSaveItem(t *testing.T) {
	mongo, dbConfig, err := setupMongo(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, mongo)
	require.NotNil(t, dbConfig)
	defer mongo.Terminate(context.TODO())
	tests := map[string]struct {
		config      *model.DatabaseConfig
		item        *hnModel.Item
		expectedErr string
	}{
		"Successfully saved": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "test",
			},
			item: &hnModel.Item{
				ID:      1,
				Dead:    true,
				Deleted: false,
			},
		},
		"No db to save to": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "",
			},
			item:        &hnModel.Item{},
			expectedErr: "Unable to save item. the Database field must be set on Operation",
		},
	}

	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewProduction()
			require.NoError(t, err)

			client, err := New(context.TODO(), logger, testConfig.config)
			require.NoError(t, err)
			if testConfig.expectedErr != "" {
				dropDatabase(dbConfig, testConfig.config.Name)
			}

			err = client.SaveItem(context.TODO(), testConfig.item)
			if testConfig.expectedErr != "" {
				assert.EqualErrorf(t, err, testConfig.expectedErr, "Request failed should be: %v, got: %v", testConfig.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
			t.Cleanup(func() {
				client.CloseConnection(context.TODO())
			})
		})
	}
}

func TestListAll(t *testing.T) {
	mongo, dbConfig, err := setupMongo(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, mongo)
	require.NotNil(t, dbConfig)
	defer mongo.Terminate(context.TODO())

	item1 := hnModel.Item{
		ID:      1,
		Dead:    true,
		Deleted: false,
		Type:    "story",
	}
	item2 := hnModel.Item{
		ID:      2,
		Dead:    true,
		Deleted: false,
		Type:    "job",
	}

	tests := map[string]struct {
		config           *model.DatabaseConfig
		expectedResponse []*hnModel.Item
		itemsToSave      []*hnModel.Item
		expectedErr      string
	}{
		"No db": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "",
			},
			expectedErr: "Failed to retrieve items. the Database field must be set on Operation",
		},
		"Returns one item": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "test",
			},
			expectedResponse: []*hnModel.Item{
				&item1,
			},
			itemsToSave: []*hnModel.Item{
				&item1,
			},
		},
		"Returns 2 items": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "test",
			},
			expectedResponse: []*hnModel.Item{
				&item1, &item2,
			},
			itemsToSave: []*hnModel.Item{
				&item1, &item2,
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewProduction()
			require.NoError(t, err)

			client, err := New(context.TODO(), logger, testConfig.config)
			require.NoError(t, err)

			if testConfig.expectedErr == "" {
				for _, item := range testConfig.itemsToSave {
					err := client.SaveItem(context.TODO(), item)
					require.NoError(t, err)
				}
			}

			items, err := client.ListAll(context.TODO())
			if testConfig.expectedErr != "" {
				assert.EqualErrorf(t, err, testConfig.expectedErr, "Request failed should be: %v, got: %v", testConfig.expectedErr, err)
				assert.Nil(t, items)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(testConfig.expectedResponse), len(items))
				assert.Equal(t, testConfig.expectedResponse, items)
			}
			t.Cleanup(func() {
				client.CloseConnection(context.TODO())
			})
		})
	}
}

func TestListStories(t *testing.T) {
	mongo, dbConfig, err := setupMongo(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, mongo)
	require.NotNil(t, dbConfig)
	defer mongo.Terminate(context.TODO())

	story1 := hnModel.Item{
		ID:      1,
		Dead:    true,
		Deleted: false,
		Type:    "story",
	}
	story2 := hnModel.Item{
		ID:      3,
		Dead:    false,
		Deleted: false,
		Type:    "story",
	}
	job1 := hnModel.Item{
		ID:      2,
		Dead:    true,
		Deleted: false,
		Type:    "job",
	}

	tests := map[string]struct {
		config           *model.DatabaseConfig
		expectedResponse []*hnModel.Item
		itemsToSave      []*hnModel.Item
		expectedErr      string
	}{
		"No db": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "",
			},
			expectedErr: "Failed to retrieve items. the Database field must be set on Operation",
		},
		"Returns one item": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "stories",
			},
			expectedResponse: []*hnModel.Item{
				&story1,
			},
			itemsToSave: []*hnModel.Item{
				&story1, &job1,
			},
		},
		"Returns 2 items": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "stories",
			},
			expectedResponse: []*hnModel.Item{
				&story1, &story2,
			},
			itemsToSave: []*hnModel.Item{
				&story1, &story2, &job1,
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewProduction()
			require.NoError(t, err)

			client, err := New(context.TODO(), logger, testConfig.config)
			require.NoError(t, err)

			if testConfig.expectedErr == "" {
				dropDatabase(dbConfig, testConfig.config.Name)
				for _, item := range testConfig.itemsToSave {
					err := client.SaveItem(context.TODO(), item)
					require.NoError(t, err)
				}
			}

			items, err := client.ListStories(context.TODO())
			if testConfig.expectedErr != "" {
				assert.EqualErrorf(t, err, testConfig.expectedErr, "Request failed should be: %v, got: %v", testConfig.expectedErr, err)
				assert.Nil(t, items)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(testConfig.expectedResponse), len(items))
				assert.Equal(t, testConfig.expectedResponse, items)
			}
			t.Cleanup(func() {
				client.CloseConnection(context.TODO())
			})
		})
	}
}

func TestListJobs(t *testing.T) {
	mongo, dbConfig, err := setupMongo(context.TODO())
	require.NoError(t, err)
	require.NotNil(t, mongo)
	require.NotNil(t, dbConfig)
	defer mongo.Terminate(context.TODO())

	story1 := hnModel.Item{
		ID:      1,
		Dead:    true,
		Deleted: false,
		Type:    "story",
	}
	job2 := hnModel.Item{
		ID:      3,
		Dead:    false,
		Deleted: false,
		Type:    "job",
	}
	job1 := hnModel.Item{
		ID:      2,
		Dead:    true,
		Deleted: false,
		Type:    "job",
	}

	tests := map[string]struct {
		config           *model.DatabaseConfig
		expectedResponse []*hnModel.Item
		itemsToSave      []*hnModel.Item
		expectedErr      string
	}{
		"No db": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "",
			},
			expectedErr: "Failed to retrieve items. the Database field must be set on Operation",
		},
		"Returns one item": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "jobs",
			},
			expectedResponse: []*hnModel.Item{
				&job1,
			},
			itemsToSave: []*hnModel.Item{
				&story1, &job1,
			},
		},
		"Returns 2 items": {
			config: &model.DatabaseConfig{
				Username: dbConfig.User,
				Password: dbConfig.Password,
				Host:     dbConfig.Host,
				Port:     fmt.Sprint(dbConfig.Port),
				Name:     "jobs",
			},
			expectedResponse: []*hnModel.Item{
				&job1, &job2,
			},
			itemsToSave: []*hnModel.Item{
				&story1, &job2, &job1,
			},
		},
	}
	for testName, testConfig := range tests {
		t.Run(testName, func(t *testing.T) {
			logger, err := zap.NewProduction()
			require.NoError(t, err)

			client, err := New(context.TODO(), logger, testConfig.config)
			require.NoError(t, err)

			if testConfig.expectedErr == "" {
				dropDatabase(dbConfig, testConfig.config.Name)
				for _, item := range testConfig.itemsToSave {
					err := client.SaveItem(context.TODO(), item)
					require.NoError(t, err)
				}
			}

			items, err := client.ListJobs(context.TODO())
			if testConfig.expectedErr != "" {
				assert.EqualErrorf(t, err, testConfig.expectedErr, "Request failed should be: %v, got: %v", testConfig.expectedErr, err)
				assert.Nil(t, items)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(testConfig.expectedResponse), len(items))
				assert.Equal(t, testConfig.expectedResponse, items)
			}
			t.Cleanup(func() {
				client.CloseConnection(context.TODO())
			})
		})
	}
}

func dropDatabase(config tcMongo.DBConfig, dbName string) {
	opts := options.Client().ApplyURI(config.ConnectionURI())
	if strings.TrimSpace(config.User) != "" && strings.TrimSpace(config.Password) != "" {
		credentials := options.Credential{
			Username: config.User,
			Password: config.Password,
		}
		opts = opts.SetAuth(credentials)
	}
	client, _ := mongo.Connect(context.TODO(), opts)

	defer client.Disconnect(context.TODO())
	database := client.Database(dbName)
	collection := database.Collection("items")
	collection.Drop(context.TODO())
	database.Drop(context.TODO())
}
