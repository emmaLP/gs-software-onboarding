package api

import (
	"context"
	"net/http"

	"github.com/emmaLP/gs-software-onboarding/internal/grpc"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Handler is an interface enabling multiple implementations of the methods within
type Handler interface {
	GetAll(c echo.Context) error
	ListStories(c echo.Context) error
	ListJobs(c echo.Context) error
	Close(ctx context.Context)
}

type apiHandler struct {
	logger     *zap.Logger
	grpcClient grpc.Client
}

// HandlerOptions give the ability to inject optional struct variables or override others
type HandlerOptions func(handler *apiHandler)

// NewHandler populates the struct of reusable variables needed for implementing the interface functions
func NewHandler(logger *zap.Logger, client grpc.Client) (*apiHandler, error) {
	return &apiHandler{
		logger:     logger,
		grpcClient: client,
	}, nil
}

func (h *apiHandler) GetAll(c echo.Context) error {
	all, err := h.grpcClient.ListAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.errorResponse(err, "Error retrieving items"))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": all,
	})
}

func (h *apiHandler) ListStories(c echo.Context) error {
	stories, err := h.grpcClient.ListStories(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, h.errorResponse(err, "Error retrieving stories"))
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"items": stories,
	})
}

func (h *apiHandler) ListJobs(c echo.Context) error {
	jobs, err := h.grpcClient.ListJobs(c.Request().Context())
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
