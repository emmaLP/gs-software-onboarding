package api

import (
	"fmt"
	"net/http"

	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

type server struct {
	logger *zap.Logger
	router *echo.Echo
}

type ClientServer interface {
	StartServer(address string)
}

// NewServer creates a new http server
func NewServer(logger *zap.Logger, client grpc.Client) (*server, error) {
	router := echo.New()
	router.HideBanner = true
	router.Use(
		middleware.Recover(),
	)

	router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "ok")
	})
	handler, err := NewHandler(logger, client)
	if err != nil {
		return nil, fmt.Errorf("Failed to create the API handler. %w", err)
	}
	router.GET("/all", handler.GetAll)
	router.GET("/stories", handler.ListStories)
	router.GET("/jobs", handler.ListJobs)
	return &server{
		logger: logger,
		router: router,
	}, nil
}

func (s *server) StartServer(address string) {
	err := s.router.Start(address)
	if err != nil {
		s.logger.Fatal("Failed to start API server", zap.Error(err))
	}
}
