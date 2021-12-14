package database

import (
	"context"
	"fmt"
	"testing"

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

			err = client.SaveItem(testConfig.item)
			if testConfig.expectedErr != "" {
				assert.EqualErrorf(t, err, testConfig.expectedErr, "Request failed should be: %v, got: %v", testConfig.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}
			t.Cleanup(func() {
				client.CloseConnection()
			})
		})

	}
}
