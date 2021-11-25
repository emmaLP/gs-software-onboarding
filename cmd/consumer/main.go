package main

import (
	"log"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Unable to create logger"))
	}
	defer logger.Sync()
}
