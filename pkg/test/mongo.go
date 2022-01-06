//go:build integration

package test

import (
	"context"
	"testing"

	tc "github.com/romnn/testcontainers"
	tcMongo "github.com/romnn/testcontainers/mongo"
	"github.com/testcontainers/testcontainers-go"
)

// SetupMongo can we used within a specific test to run a mongo docker container
func SetupMongo(t *testing.T, ctx context.Context) (mongoC testcontainers.Container, Config tcMongo.DBConfig, err error) {
	t.Helper()
	return SetupMongoContainer(ctx)
}

// SetupMongoContainer to run a mongo docker container used test containers. In order to use this you need to have the docker.sock running
func SetupMongoContainer(ctx context.Context) (mongoC testcontainers.Container, Config tcMongo.DBConfig, err error) {
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
