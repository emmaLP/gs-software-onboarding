package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"github.com/emmaLP/gs-software-onboarding/internal/model"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Handler interface {
	GetAll(c echo.Context) error
	ListStories(c echo.Context) error
	ListJobs(c echo.Context) error
}

type apiHandler struct {
	logger   *zap.Logger
	config   *model.Configuration
	dbClient database.Client
}

func NewHandler(ctx context.Context, logger *zap.Logger, config *model.Configuration) (*apiHandler, error) {
	databaseClient, err := database.New(ctx, logger, &config.Database)
	if err != nil {
		return nil, fmt.Errorf("Unexpected error when connecting to the database. %w", err)
	}
	defer databaseClient.CloseConnection(ctx)
	return &apiHandler{
		logger:   logger,
		dbClient: databaseClient,
		config:   config,
	}, nil
}

func (h *apiHandler) GetAll(c echo.Context) error {
	all, err := h.dbClient.ListAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("Error retrieving items from db. %w", err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": all,
	})
}

func (h *apiHandler) ListStories(c echo.Context) error {
	stories, err := h.dbClient.ListStories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("Error retrieving stories from db. %w", err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": stories,
	})
}

func (h *apiHandler) ListJobs(c echo.Context) error {
	jobs, err := h.dbClient.ListJobs(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("Error retrieving jobs from db. %w", err))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": jobs,
	})
}
