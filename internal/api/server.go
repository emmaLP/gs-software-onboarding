package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type server struct {
	config *model.Configuration
	logger *zap.Logger
	router *echo.Echo
}

type ClientServer interface {
	StartServer(address string)
}

// NewServer creates a new http server
func NewServer(ctx context.Context, logger *zap.Logger, config *model.Configuration) (*server, error) {
	router := echo.New()
	router.HideBanner = true
	router.Use(
		middleware.Recover(),
	)

	router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})
	handler, err := NewHandler(ctx, logger, &config.Database)
	if err != nil {
		return nil, fmt.Errorf("Failed to create the API handler. %w", err)
	}
	router.GET("/all", handler.GetAll)
	router.GET("/stories", handler.ListStories)
	router.GET("/jobs", handler.ListJobs)
	return &server{
		logger: logger,
		router: router,
		config: config,
	}, nil
}

func (s *server) StartServer() {
	err := s.router.Start(s.config.Api.Address)
	if err != nil {
		s.logger.Fatal("Failed to start API server", zap.Error(err))
	}
}
