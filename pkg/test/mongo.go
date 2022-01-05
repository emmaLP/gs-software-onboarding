//go:build integration

package test

import (
	"context"
	"testing"

	tc "github.com/romnn/testcontainers"
	tcMongo "github.com/romnn/testcontainers/mongo"
	"github.com/testcontainers/testcontainers-go"
)

func SetupMongo(t *testing.T, ctx context.Context) (mongoC testcontainers.Container, Config tcMongo.DBConfig, err error) {
	t.Helper()
	return SetupMongoContainer(ctx)
}

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
