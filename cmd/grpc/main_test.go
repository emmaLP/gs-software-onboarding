//go:build integration

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"testing"

	"github.com/emmaLP/gs-software-onboarding/pkg/test"
)

func TestMain(m *testing.M) {
	log.Println("Starting integration tests")
	os.Setenv("GRPC_PORT", strconv.Itoa(randomPort()))

	ctx := context.TODO()
	_, dbConfig, err := test.SetupMongoContainer(ctx)
	if err != nil {
		log.Println("FAIL - unable to setup mongo")
		os.Exit(1)
	}

	os.Setenv("DATABASE_NAME", "int-test")
	os.Setenv("DATABASE_USERNAME", dbConfig.User)
	os.Setenv("DATABASE_PASSWORD", dbConfig.Password)
	os.Setenv("DATABASE_HOST", dbConfig.Host)
	os.Setenv("DATABASE_PORT", fmt.Sprint(dbConfig.Port))

	redis, err := test.SetupRedis(ctx)
	if err != nil {
		log.Println("FAIL - unable to setup redis")
		os.Exit(1)
	}
	os.Setenv("CACHE_ADDRESS", redis.URI)

	go func() {
		main()
	}()
	os.Exit(m.Run())
}

func randomPort() int {
	min := 51000
	max := 55000
	return rand.Intn(max-min) + min
}
