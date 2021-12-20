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
	config   *model.DatabaseConfig
	dbClient database.Client
}

type HandlerOptions func(handler *apiHandler)

func NewHandler(ctx context.Context, logger *zap.Logger, config *model.DatabaseConfig, opts ...HandlerOptions) (*apiHandler, error) {
	handler := &apiHandler{
		logger: logger,
		config: config,
	}

	for _, opt := range opts {
		opt(handler)
	}

	if handler.dbClient == nil {
		databaseClient, err := database.New(ctx, logger, config)
		if err != nil {
			return nil, fmt.Errorf("Unexpected error when connecting to the database. %w", err)
		}
		handler.dbClient = databaseClient
	}
	return handler, nil
}

func WithDatabaseClient(client database.Client) HandlerOptions {
	return func(h *apiHandler) {
		h.dbClient = client
	}
}

func (h *apiHandler) GetAll(c echo.Context) error {
	all, err := h.dbClient.ListAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.errorResponse(err, "Error retrieving items"))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": all,
	})
}

func (h *apiHandler) ListStories(c echo.Context) error {
	stories, err := h.dbClient.ListStories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.errorResponse(err, "Error retrieving stories"))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": stories,
	})
}

func (h *apiHandler) ListJobs(c echo.Context) error {
	jobs, err := h.dbClient.ListJobs(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.errorResponse(err, "Error retrieving jobs"))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": jobs,
	})
}

func (h *apiHandler) errorResponse(err error, errMsg string) map[string]interface{} {
	return map[string]interface{}{
		"error_message": errMsg,
		"error":         err,
	}
}
