package grpc

import (
	"fmt"

	"github.com/emmaLP/gs-software-onboarding/internal/caching"
	"github.com/emmaLP/gs-software-onboarding/pkg/common/model"
	pb "github.com/emmaLP/gs-software-onboarding/pkg/grpc/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	pb.UnimplementedAPIServer
	ItemCache caching.Client
}

func (h Handler) ListAll(empty *emptypb.Empty, s pb.API_ListAllServer) error {
	return h.streamItems(s, func() ([]*model.Item, error) {
		return h.ItemCache.ListAll(s.Context())
	})
}

func (h Handler) ListStories(empty *emptypb.Empty, s pb.API_ListStoriesServer) error {
	return h.streamItems(s, func() ([]*model.Item, error) {
		return h.ItemCache.ListStories(s.Context())
	})
}

func (h Handler) ListJobs(empty *emptypb.Empty, s pb.API_ListJobsServer) error {
	return h.streamItems(s, func() ([]*model.Item, error) {
		return h.ItemCache.ListJobs(s.Context())
	})
}

func (h Handler) streamItems(server interface{ Send(item *pb.Item) error }, itemsFunc func() ([]*model.Item, error)) error {
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
