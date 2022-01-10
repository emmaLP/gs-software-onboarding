package grpc

import (
	"context"
	"fmt"

	"github.com/emmaLP/gs-software-onboarding/internal/database"
	"go.uber.org/zap"

	"github.com/emmaLP/gs-software-onboarding/internal/caching"
	"github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	pb "github.com/emmaLP/gs-software-onboarding/pkg/grpc/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedAPIServer
	itemCache caching.Client
	dbClient  database.Client
	logger    *zap.Logger
}

func NewHandler(itemCache caching.Client, dbClient database.Client, logger *zap.Logger) *Handler {
	return &Handler{
		itemCache: itemCache,
		dbClient:  dbClient,
		logger:    logger,
	}
}

func (h *Handler) ListAll(empty *emptypb.Empty, s pb.API_ListAllServer) error {
	return h.streamItems(s, func() ([]*model.Item, error) {
		return h.itemCache.ListAll(s.Context())
	})
}

func (h *Handler) ListStories(empty *emptypb.Empty, s pb.API_ListStoriesServer) error {
	return h.streamItems(s, func() ([]*model.Item, error) {
		return h.itemCache.ListStories(s.Context())
	})
}

func (h *Handler) ListJobs(empty *emptypb.Empty, s pb.API_ListJobsServer) error {
	return h.streamItems(s, func() ([]*model.Item, error) {
		return h.itemCache.ListJobs(s.Context())
	})
}

func (h *Handler) SaveItem(ctx context.Context, item *pb.Item) (*pb.ItemResponse, error) {
	toItem := model.PItemToItem(item)
	err := h.dbClient.SaveItem(ctx, &toItem)
	if err != nil {
		h.logger.Error("Failed to save item to the database", zap.Error(err))
		return &pb.ItemResponse{Id: item.Id, Success: false}, err
	}
	return &pb.ItemResponse{Id: item.Id, Success: true}, nil
}

func (h *Handler) streamItems(server interface{ Send(item *pb.Item) error }, itemsFunc func() ([]*model.Item, error)) error {
	items, err := itemsFunc()
	if err != nil {
		return fmt.Errorf("fetching all items, %w", err)
	}

	for _, item := range items {
		if err := server.Send(model.ItemToPItem(*item)); err != nil {
			return fmt.Errorf("steaming item to client. %w", err)
		}
	}
	return nil
}
